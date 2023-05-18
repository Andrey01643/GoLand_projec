package models

// Структура для передачи данных в HTML шаблоны
type TemplateData struct {
	Variants []Variant
	Variant  Variant
	Task     Task
	Result   int
}
