package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

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

	// Проверяем наличие cookie с идентификатором пользователя
	cookie, err := r.Cookie("user_id")
	if err != nil {
		// Если cookie отсутствует, редиректим на страницу авторизации
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	fmt.Println(cookie.Value)

	// Выполняем запрос на выборку пользователя по идентификатору
	var user models.User
	err = db.QueryRow("SELECT * FROM auth WHERE id=$1", cookie.Value).Scan(&user.ID, &user.Login, &user.Password, &user.IsAuthorized, &user.LoginTime, &user.LogoutTime)
	if err != nil {
		// Если пользователь не найден в БД, редиректим на страницу авторизации
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Получение ID выбранного варианта из URL
	idStr := strings.TrimPrefix(r.URL.Path, "/test/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Получение информации о варианте из БД
	var v models.Variant
	err = db.QueryRow("SELECT id, name FROM variants WHERE id = $1", id).Scan(&v.ID, &v.Name)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Получение всех заданий из БД для выбранного варианта
	rows, err := db.Query("SELECT id, variant_id, task, correct_answer, answer_1, answer_2, answer_3, answer_4 FROM tasks WHERE variant_id = $1 ORDER BY id", id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer rows.Close()

	// Список заданий
	var tasks []models.Task
	var t models.Task
	// Перебираем результаты запроса и добавляем задания в список
	for rows.Next() {

		err = rows.Scan(&t.ID, &t.VariantID, &t.Task, &t.CorrectAnswer, &t.Answer1, &t.Answer2, &t.Answer3, &t.Answer4)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		tasks = append(tasks, t)
	}
	err = rows.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отображение страницы с тестом и передача информации о варианте и заданиях в шаблон
	tpl, err := template.ParseFiles("static/tests.html")
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

	// Проверяем наличие cookie с идентификатором пользователя
	cookie, err := r.Cookie("user_id")
	if err != nil {
		// Если cookie отсутствует, редиректим на страницу авторизации
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	fmt.Println(cookie.Value)

	// Выполняем запрос на выборку пользователя по идентификатору
	var user models.User
	err = db.QueryRow("SELECT * FROM auth WHERE id=$1", cookie.Value).Scan(&user.ID, &user.Login, &user.Password, &user.IsAuthorized, &user.LoginTime, &user.LogoutTime)
	if err != nil {
		// Если пользователь не найден в БД, редиректим на страницу авторизации
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
	var v models.Variant
	err = db.QueryRow("SELECT id, name FROM variants WHERE id = $1", variantID).Scan(&v.ID, &v.Name)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Проверяем правильность ответа
	var correctAnswer string
	err = db.QueryRow("SELECT correct_answer FROM tasks WHERE id = $1", taskID).Scan(&correctAnswer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	isCorrect := answer == correctAnswer

	// Добавление id автора с номером задания, какой был ответ отослан и в какой момент времени
	fmt.Println(answer)
	_, err = db.Exec("INSERT INTO answers (auth_id, variant_id, task_id, answer, answer_time) VALUES ($1, $2, $3, $4, $5)", cookie.Value, variantID, taskID, answer, time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Добавляем результат в таблицу results
	_, err = db.Exec("INSERT INTO results (variant_id, task_id, user_id, is_correct, answer_time) VALUES ($1, $2, $3, $4, $5)", variantID, taskID, cookie.Value, isCorrect, time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Получение всех заданий из БД для выбранного варианта
	rows, err := db.Query(`
	SELECT id, variant_id, task, correct_answer, answer_1, answer_2, answer_3, answer_4
FROM tasks
WHERE variant_id = $1 AND id NOT IN (
    SELECT task_id
    FROM answers
)
ORDER BY id`, variantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	// Список заданий
	var tasks []models.Task
	var t models.Task

	// Перебираем результаты запроса и добавляем задания в список
	for rows.Next() {
		err = rows.Scan(&t.ID, &t.VariantID, &t.Task, &t.CorrectAnswer, &t.Answer1, &t.Answer2, &t.Answer3, &t.Answer4)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		tasks = append(tasks, t)
	}
	if len(tasks) == 0 {
		// Если отсутствуют задания для выбранного варианта, перенаправляем на другую страницу
		http.Redirect(w, r, fmt.Sprintf("/result?variantID=%d", variantID), http.StatusSeeOther) //////////////////////////////////
		return
	}
	err = rows.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отображение страницы с тестом и передача информации о варианте и заданиях в шаблон
	tpl, err := template.ParseFiles("static/tests.html")
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

	// Проверяем наличие cookie с идентификатором пользователя
	cookie, err := r.Cookie("user_id")
	if err != nil {
		// Если cookie отсутствует, редиректим на страницу авторизации
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// Получаем значения параметров Variant и Task из URL
	vID := r.URL.Query().Get("variantID")
	variantID, err := strconv.Atoi(vID)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	// Выполняем запрос на выборку пользователя по идентификатору
	var user models.User
	err = db.QueryRow("SELECT * FROM auth WHERE id=$1", cookie.Value).Scan(&user.ID, &user.Login, &user.Password, &user.IsAuthorized, &user.LoginTime, &user.LogoutTime)
	if err != nil {
		// Если пользователь не найден в БД, редиректим на страницу авторизации
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Выполняем запрос на выборку количества правильных ответов для данного пользователя и выбранного варианта
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM results JOIN tasks ON results.task_id=tasks.id WHERE results.user_id=$1 AND results.is_correct=true AND tasks.variant_id=$2", cookie.Value, variantID).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Выполняем запрос на выборку общего количества заданий
	var total int
	err = db.QueryRow("SELECT COUNT(*) FROM tasks WHERE tasks.variant_id = $1 AND tasks.variant_id IN (SELECT variant_id FROM answers)", variantID).Scan(&total)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Вычисляем процент правильных ответов
	result := float64(count) / float64(total) * 100

	// Отображение страницы с результатами теста и передача списка результатов в шаблон
	tpl, err := template.ParseFiles("static/result.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(w, result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
