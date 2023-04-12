package models

type Record struct {
	Id         int    `json:"id"`
	Uid        string `json:"uid"`
	Domain     string `json:"domain"`
	Cn         string `json:"cn"`
	Department string `json:"department"`
	Title      string `json:"title"`
	Who        string `json:"who"`
}
