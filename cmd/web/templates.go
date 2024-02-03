package main

import "interlude/pkg/models"

/* Данный файл создан для того чтобы
хранить динамические данные шаблонов
- так как 0пакее html/template позволяет
передавать только один источник динамических
шаблонов, мы создадим структуру с шаблонами
*/

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
