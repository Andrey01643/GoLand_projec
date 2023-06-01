package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"httpServer/db/postgres"
	"httpServer/models"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем соединение с БД
	db, err := postgres.ConnectDb()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	user, cookie, err := checkUserCookie(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if cookie == nil {
		http.Error(w, "User cookie not found", http.StatusNotFound)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/test/")
	variantID, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	v, err := postgres.GetVariantByID(db, variantID)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	var t models.Task
	t, err = postgres.GetTasksByVariantID(db, variantID, t)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Отображение страницы с тестом и передача информации о варианте и заданиях в шаблон
	tpl, err := template.ParseFiles("web/tests.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(w, struct {
		Variant models.Variant
		Task    models.Task
	}{
		Variant: v,
		Task:    t,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func NextTestHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем соединение с БД
	db, err := postgres.ConnectDb()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	user, cookie, err := checkUserCookie(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if cookie == nil {
		http.Error(w, "User cookie not found", http.StatusNotFound)
		return
	}

	// Получаем значения параметров Variant и Task из URL
	vID := mux.Vars(r)["variantID"]
	variantID, err := strconv.Atoi(vID)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	tID := mux.Vars(r)["taskID"]
	taskID, err := strconv.Atoi(tID)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	answer := r.FormValue("answer")
	// Получение информации о варианте из БД
	v, err := postgres.GetVariantByID(db, variantID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Проверяем правильность ответа
	correctAnswer, err := postgres.GetCorrectAnswer(db, taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	isCorrect := answer == correctAnswer

	err = postgres.InsertAnswer(db, cookie.Value, variantID, taskID, answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Добавляем результат в таблицу results
	err = postgres.InsertResult(db, variantID, taskID, cookie.Value, isCorrect)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Получение всех заданий из БД для выбранного варианта
	var t models.Task
	t, err = postgres.GetTasks(db, variantID, t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(t.Task) == 0 {
		// Если отсутствуют задания для выбранного варианта, перенаправляем на другую страницу
		http.Redirect(w, r, fmt.Sprintf("/result?variantID=%d", variantID), http.StatusSeeOther) //////////////////////////////////
		return
	}
	// Отображение страницы с тестом и передача информации о варианте и заданиях в шаблон
	tpl, err := template.ParseFiles("web/tests.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(w, struct {
		Variant models.Variant
		Task    models.Task
		Answer  string
	}{
		Variant: v,
		Task:    t,
		Answer:  answer,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func ResultTestsHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем соединение с БД
	db, err := postgres.ConnectDb()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	user, cookie, err := checkUserCookie(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if cookie == nil {
		http.Error(w, "User cookie not found", http.StatusNotFound)
		return
	}

	// Получаем значения параметров Variant и Task из URL
	vID := r.URL.Query().Get("variantID")
	variantID, err := strconv.Atoi(vID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	result, err := postgres.GetTestResults(db, cookie.Value, variantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отображение страницы с результатами теста и передача списка результатов в шаблон
	tpl, err := template.ParseFiles("web/result.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(w, result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
