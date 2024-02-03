package main

import (
	"database/sql"
	"flag"
	"interlude/pkg/models/mysql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

/* Сруктура служащая для хранения зависимостей всего приложения*/
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *mysql.SnippetModel
}

func main() {
	/* Флаг командной строки, хранящийся в переменной.
	Хранит данные об сетевом адресе HTTP */
	addr := flag.String("addr", ":4000", "Сетевой адрес HTTP")
	/* Считываем флаг из командной строки - функция
	flag.Parse() считывет флаг из комндной строки и присваивает
	его переменной. Если сделать это до использования переменной,
	то будет выставлено значение по умлочанию*/
	flag.Parse()

	/* Определение нового флага из командной строки для
	настройки SQL подключения */
	dsn := flag.String("dsn", "web:Password/snippetbox?parseTime=true", "Название MySQL источника данных")
	flag.Parse()

	/* Логгирование*/
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	/* Lshortfile нужен для включения названия файла и
	номера строки где обнаружена ошибка */
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	/* В целях экономии места в функции main код для создания
	пула соединений выносится в отедльную функцию openDB(). В
	нее передаем полученный источник данных из командной строки */
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	/* Откладывание вызова db.Close() для того чтобы
	пул соединений был закрыт до выхода из функции main() */
	defer db.Close()

	/* Инициализация экземпляра mysql.SnippetModel и
	добавление его в зависимости*/
	/*Новая структура для зависимостей приложения */
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &mysql.SnippetModel{DB: db},
	}

	/* Структура для логов */
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	/* Применение логеров */
	infoLog.Printf("Запуск сервера на %s . . .", *addr)
	/* Вызов метода от нашей новой структуры */
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

/*
	Функция openDB обертывает sql.Open() и возвращает

пул соединений для заданной строки подключения (DSN)
*/
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

/*
Реализация метода защиты от приема, когда
запрос HTTP ведет к папке - вместо отображения
содержимого мы возвращаем 404
*/
/*type neutredFileSystem struct {
	fs http.FileSystem
}*/

/*
Когда http.FileServer получает запрос, мы создаем метод Open()
В нем мы открываем вызываемый путь, и проверяем является ли он
папкой. Далее проверяем наличие файла index.html внутри папки
с помощью Stat("index.html")

Если файл не существует, метод возвращает ошибку os.ErrNotExist,
которая будет преобразована http.FileServer в 404
Далее вызываем Close() для закрытия только что открытого index.html
для избежания утечки файловго дескриптора

Во всех остальных случаях мы возвращаем файл и http.FileServer
работает стандартно
*/
/*func (nfs neutredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}*/
