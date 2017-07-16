// Пакет "projects" содержит набор служебных функция для работы с проектами
// ToDo:
// Рефакторинг кода

package projects

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Project struct {
	Id    int
	Name  string
	Pages int
}

var db *sql.DB
var err error

// Функция "Init" устанавливает соединение с БД
func init() {
    // Открываем соединение
	db, err = sql.Open("postgres", "user=urivsky password=123581321 dbname=mindassistant sslmode=disable")
	errorHandler(err, 500, nil)

    // Отслеживаем состояние канала передачи данных
    errorHandler(db.Ping(), 500, nil)
}

// Фукнция "List" отображает список проектов
func List(w http.ResponseWriter, r *http.Request) {
    // Проверка заголовков
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

    // Подготовка запроса
	rows, err := db.Query("SELECT * FROM projects")
	errorHandler(err, 500, w)
	defer rows.Close()

    // Сбор данных из БД в структуру
	projects := make([]*Project, 0)
	for rows.Next() {
		project := new(Project)

		err := rows.Scan(&project.Id, &project.Name, &project.Pages)
		errorHandler(err, 500, w)

		projects = append(projects, project)
	}
    errorHandler(rows.Err(), 500, w)

    // Подготавка выходных данных
	output, err := json.Marshal(projects)
	errorHandler(err, 500, w)

    // Отображение результата
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(output)
}

// Функция "Item" отображет проект по его id
func Item(w http.ResponseWriter, r *http.Request) {
    // Проверка заголовков
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

    // Сбор и анализ входных данных
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

    // Подготовка запроса
	row := db.QueryRow("SELECT * FROM projects WHERE id = $1", id)
	project := new(Project)

    // Сбор данных из БД в структуру
	err := row.Scan(&project.Id, &project.Name, &project.Pages)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else {
        errorHandler(err, 500, w)
    }

    // Подготавка выходных данных
	output, err := json.Marshal(project)
	errorHandler(err, 500, w)

    // Отображение результата
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(output)
}

// Функция "Create" создает новый проект
func Create(w http.ResponseWriter, r *http.Request) {
    // Проверка заголовков
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

    // Сбор и анализ входных данных
	decoder := json.NewDecoder(r.Body)
	project := new(Project)
	err := decoder.Decode(&project)
	errorHandler(err, 500, w)
	defer r.Body.Close()

	if project.Name == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

    // Подготовка запроса
	createProject, err := db.Prepare("INSERT INTO projects(name, pages) VALUES($1, $2);")
	errorHandler(err, 500, w)

    // Выполнение запроса
	result, err := createProject.Exec(project.Name, project.Pages)
	errorHandler(err, 500, w)

    // Проверка на успешность
	rowsAffected, err := result.RowsAffected()
	errorHandler(err, 500, w)

    // Отображение результата
	if rowsAffected > 0 {
		fmt.Fprintf(w, "%t\n", true)
	}
}

// Функция "Update" обновляет данные проекта
func Update(w http.ResponseWriter, r *http.Request) {
    // Проверка заголовков
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

    // Сбор и анализ входных данных
	decoder := json.NewDecoder(r.Body)
	project := new(Project)
	err := decoder.Decode(&project)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if project.Id == 0 || project.Name == "" || project.Pages == 0 {
		http.Error(w, http.StatusText(400), 400)
		return
	}

    // Подготовка запроса
	updateProject, err := db.Prepare("UPDATE projects SET name = $1, pages = $2 WHERE id = $3")
	errorHandler(err, 500, w)

    // Выполнение запроса
	result, err := updateProject.Exec(project.Name, project.Pages, project.Id)
	errorHandler(err, 500, w)

    // Проверка на успешность
	rowsAffected, err := result.RowsAffected()
	errorHandler(err, 500, w)

    // Отображение результата
	if rowsAffected > 0 {
		fmt.Fprintf(w, "%t\n", true)
	}
}

// Функция "Delete" удаляет проект по его id и все связанные с ним данные
func Delete(w http.ResponseWriter, r *http.Request) {
    // Проверка заголовков
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

    // Сбор и анализ входных данных
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

    // Подготовка запроса на удаление fields связанных с проектом
    deleteFields, err := db.Prepare("DELETE FROM fields f USING (field_groups g, data d) WHERE g.id = f.group_id and g.data = d.id and d.project = $1")
    errorHandler(err, 500, w)

    // Подготовка запроса на удаление field_groups связанных с проектом
    deleteFieldGroups, err := db.Prepare("DELETE FROM field_groups g USING data d WHERE g.data = d.id and d.project = $1")
    errorHandler(err, 500, w)

    // Подготовка запроса на удаление data связанных с проектом
	deleteData, err := db.Prepare("DELETE FROM data WHERE project = $1")
	errorHandler(err, 500, w)

    // Подготовка запроса на удаление проекта
	deleteProject, err := db.Prepare("DELETE FROM projects WHERE id = $1")
	errorHandler(err, 500, w)

    // Выполнение запросов
    deleteFields.Exec(id)
    deleteFieldGroups.Exec(id)
	deleteData.Exec(id)

	projectResult, err := deleteProject.Exec(id)
	errorHandler(err, 500, w)

    // Проверка на успешность
	rowsAffected, err := projectResult.RowsAffected()
	errorHandler(err, 500, w)
    
    // Отображение результата
	if rowsAffected > 0 {
		fmt.Fprintf(w, "%t\n", true)
	}
}

// Функция для проверки ошибок
func errorHandler(err error, code int, w http.ResponseWriter) {
    if err != nil {
        // log.Print(err) // <= Режим отладки
        if w != nil {
            http.Error(w, http.StatusText(500), 500)
            return
        } else {
            log.Fatal(err)
        }
    }
}