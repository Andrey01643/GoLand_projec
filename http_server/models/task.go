package models

type Task struct {
	ID            int
	VariantID     int
	Task          string
	CorrectAnswer string
	Answer1       string
	Answer2       string
	Answer3       string
	Answer4       string
}
