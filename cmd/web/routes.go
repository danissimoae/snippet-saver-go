package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	/* Простой способ организации обработчика HTTP запросов

	Когда сервер получает HTTP запрос, он вызывает метод
	ServeHTTP() от servemux

	Далее он ищет соответствующий обработчик на основе URL
	запроса и вызывает метод ServeHTTP() данного обработчика,
	образуя цепочку ServeHTTP() методов */
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.indexHandler)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	/* Инициализация FileServer - обработчика HTTP запрова к
	статическим файлам. Путь в функции http.Dir является
	относительным корневой папке проекта. */
	fileServer := http.FileServer(http.Dir("./ui/static"))

	/* Регистрация функции mux.Handle() для регистрации обработчика
	всех запросов, начинающихся с "/static/" - убираем его со
	всех запросов, перед тем как запрос достигнет http.FileServer*/
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
