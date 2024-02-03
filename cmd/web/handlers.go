package main

import (
	"errors"
	"fmt"
	"html/template"
	"interlude/pkg/models"
	"net/http"
	"strconv"
)

/* Обработчик для зависимостей application */
func (app *application) indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	/* Экземляр структуры templateData содержащий срез
	с заметками */
	data := &templateData{Snippets: s}

	files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	/* Обработка шаблона, читаем файлы шаблона и если
	возникла ошибка - отправляем тело с ошибкой */
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	/* Вызов метода execute для записи содержимого шаблона
	в тело HTTP ответа, последний параметр отвечает за возможность
	динамической оптравки данных в шаблон */
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}

}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	/* Происходит проверка на валидность query запроса */
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	/* Вызов метода Get из модели Snipping для извелчения
	данных по их ID. Если не найдено, возвращается 404 */
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	/* Экземляр структуры templateData содержащий срез
	с заметками */
	data := &templateData{Snippet: s}

	/* Срез содержащий пути к темплейтам */
	files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	/* Парсинг файлов шаблонов */
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	/* Выполнение шаблонов */
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// Даем пользователю знать, какие методы разрешены
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	title := "Метод черного ящика"
	content := "Метод выяснения ошибки, пришедший из авиастроения"
	expires := "7"

	/* Передача данных в метод SnippetModel.Insert() и
	получение ID созданной записи */
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	/* Перенаправление пользователя на
	соотвтствующую страницу заметки*/
	http.Redirect(w, r, fmt.Sprintf("/snippet&id=%d", id), http.StatusSeeOther)
}
