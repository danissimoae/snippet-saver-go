package mysql

import (
	"database/sql"
	"errors"

	"interlude/pkg/models"
)

/* Определение типа обрабатывающего пул подключений */
type SnippetModel struct {
	DB *sql.DB
}

/* Insert - метод для создания новой заметки в БД*/
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	/* Выполнение метода Exec() - передаем SQL-запрос и данные
	о заметке. Метод возвращает объект sql.Result который
	содержит некоторую информацию о том что произошло
	после запроса */
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, nil
	}

	/* Используется метод LastInsertId для получения
	последнего ID созданной записи из таблицы snippets*/
	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(id), nil
}

/*
Get - метод для возвращения данных заметки по ее
идентификатору
*/
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`

	/* Выполнение метода QueryRow() для выполнения SQL
	запроса, передается ненадежный id в качестве значения
	для плейсхолдера, возвращается указатель на объект
	sql.Row, который содержит данные записи*/
	row := m.DB.QueryRow(stmt, id)

	/* Инициализация указателя на новую структуру Snippet */
	s := &models.Snippet{}

	/* Row.Scan() - используется для копирования значения из
	каждого поля от sql.Row() в соостветствующее поле
	в структуре Snippet, аргументы для него - это указатели на место,
	куда требуется скопировать данные и количество аргументов
	должно быть равно количеству столбцов в БД*/
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	/* Проверка функцией - если запрос был с ошибкой
	Если ошибка обнаружена, то возвращаем нашу ошибку из модели
	models.ErrNoRecord
	Также эта ошибка является предопределенной -
	как экспешн в пайтон - мы можем ее "отловить"*/
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

/*
Latest - метод возвращающий 10 наиболее часто
используемых заметок
*/
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `select id, title, content, created, expires from snippets where expires > utc_timestamp() order by created desc limit 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	/* Откладывание вызов rows.Closed() чтобы быть увереным что запрос правильно закроется
	перед выполнением Latest, оператор сработает *после* проверки
	на наличие ошибки в методе Query, иначе Query вернет ошибку
	и это привидет к панике*/
	defer rows.Close()

	var snippets []*models.Snippet

	/* Используем rows.Next() для отбора результатов*/
	for rows.Next() {
		s := &models.Snippet{}
		/* Использование rows.Scan() чтобы скопировать значения полей
		в структуру*/
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		/* Добавление структуры в срез*/
		snippets = append(snippets, s)
	}

	/* Чтобы узнать не возникла ли ошибка когда
	rows.Next() закрывается вызываем rows.Err()*/
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
