package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

/*
	Помощник serverError записывает сообщение об ошибке в

errorLog и отпарвляет пользователю ответ 500
*/
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

/* Помощник clientError отправляет код состояния пользователю */
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

/* Помощник notFound - удобаня оболочка для clientError */
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
