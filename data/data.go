// Пакет "data" содержит набор служебных функция для работы с Data объектами
// ToDo:
// Рефакторинг кода

package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

    _ "github.com/lib/pq"
)

type DataDb struct {
    Id          int            `json:"id"`
    Name        string         `json:"name"`
    Project     int            `json:"project"`
    Parent      int            `json:"parent"`
    Coordinates []byte         `json:"coordinates"`
}
type DataJson struct {
    Id          int            `json:"id"`
    Name        string         `json:"name"`
    Project     int            `json:"project"`
    Parent      int            `json:"parent"`
    Coordinates map[string]int `json:"coordinates"`
    Content     []*FieldGroup  `json:"content"`
}
type FieldGroup struct {
    Id     int      `json:"id"`
    Name   string   `json:"name"`
    Order  int      `json:"order"`
    Data   int      `json:"data"`
    Fields []*Field `json:"fields"`
}
type Field struct {
    Id    int    `json:"id"`
    Type  string `json:"type"`
    Value string `json:"value"`
    Order int    `json:"order"`
    Group int    `json:"group"`
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

// Функция "List" отображает список всех Data объектов
func List(w http.ResponseWriter, r *http.Request) {
    // Проверка заголовков
    if r.Method != "GET" {
        http.Error(w, http.StatusText(405), 405)
        return
    }

    // Подготовка запроса
    dataRows, err := db.Query("SELECT * FROM data")
    errorHandler(err, 500, w)
    defer dataRows.Close()

    // Сбор данных из БД в структуру
    dataList := make([]*DataJson, 0)
    for dataRows.Next() {
        data := new(DataDb)

        err := dataRows.Scan(&data.Id, &data.Name, &data.Project, &data.Parent, &data.Coordinates)
        errorHandler(err, 500, w)

        // Обработка координат
        var coordinates map[string]int
        json.Unmarshal([]byte(data.Coordinates), &coordinates)

        // Сбор данных из таблицы field_group связанных с Data объектом
        dataFieldGroupRows, err := db.Query("SELECT * FROM field_groups WHERE data = $1", data.Id)
        errorHandler(err, 500, w)
        defer dataFieldGroupRows.Close()

        content := make([]*FieldGroup, 0)
        for dataFieldGroupRows.Next() {
            dataFieldGroup := new(FieldGroup)

            err := dataFieldGroupRows.Scan(&dataFieldGroup.Id, &dataFieldGroup.Name, &dataFieldGroup.Order, &dataFieldGroup.Data)
            errorHandler(err, 500, w)

            // Сбор данных из таблицы fields связанных с Data объектом
            dataFieldsRows, err := db.Query("SELECT * FROM fields WHERE group_id = $1", dataFieldGroup.Id)
            errorHandler(err, 500, w)
            defer dataFieldsRows.Close()

            dataFieldsList := make([]*Field, 0)
            for dataFieldsRows.Next() {
                dataFields := new(Field)

                err := dataFieldsRows.Scan(&dataFields.Id, &dataFields.Type, &dataFields.Order, &dataFields.Value, &dataFields.Group)
                errorHandler(err, 500, w)

                dataFieldsList = append(dataFieldsList, dataFields)
            }

            // Трансформация группы в новый объект
            dataFieldGroupResult := &FieldGroup {
                Id:          dataFieldGroup.Id,
                Name:        dataFieldGroup.Name,
                Order:       dataFieldGroup.Order,
                Data:        dataFieldGroup.Data,
                Fields:      dataFieldsList,
            }

            content = append(content, dataFieldGroupResult)
        }

        // Трансформация Data в новый объект
        dataResult := &DataJson {
            Id:          data.Id,
            Name:        data.Name,
            Project:     data.Project,
            Parent:      data.Parent,
            Coordinates: coordinates,
            Content:     content,
        }

        dataList = append(dataList, dataResult)
    }
    errorHandler(dataRows.Err(), 500, w)

    // Подготавка выходных данных
    output, err := json.Marshal(dataList)
    errorHandler(err, 500, w)

    // Отображение результата
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Write(output)
}

// Функция "ListByProject" отображает список Data объектов по id проекта
func ListByProject(w http.ResponseWriter, r *http.Request) {
    // Проверка заголовков
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

    // Сбор и анализ входных данных
	project := r.FormValue("project")
	if project == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

    // Подготовка запроса
	dataRows, err := db.Query("SELECT * FROM data WHERE project = $1", project)
	errorHandler(err, 500, w)
	defer dataRows.Close()

    // Сбор данных из БД в структуру
    dataList := make([]*DataJson, 0)
    for dataRows.Next() {
        data := new(DataDb)

        err := dataRows.Scan(&data.Id, &data.Name, &data.Project, &data.Parent, &data.Coordinates)
        errorHandler(err, 500, w)

        // Обработка координат
        var coordinates map[string]int
        json.Unmarshal([]byte(data.Coordinates), &coordinates)

        // Сбор данных из таблицы field_group связанных с Data объектом
        dataFieldGroupRows, err := db.Query("SELECT * FROM field_groups WHERE data = $1", data.Id)
        errorHandler(err, 500, w)
        defer dataFieldGroupRows.Close()

        content := make([]*FieldGroup, 0)
        for dataFieldGroupRows.Next() {
            dataFieldGroup := new(FieldGroup)

            err := dataFieldGroupRows.Scan(&dataFieldGroup.Id, &dataFieldGroup.Name, &dataFieldGroup.Order, &dataFieldGroup.Data)
            errorHandler(err, 500, w)

            // Сбор данных из таблицы fields связанных с Data объектом
            dataFieldsRows, err := db.Query("SELECT * FROM fields WHERE group_id = $1", dataFieldGroup.Id)
            errorHandler(err, 500, w)
            defer dataFieldsRows.Close()

            dataFieldsList := make([]*Field, 0)
            for dataFieldsRows.Next() {
                dataFields := new(Field)

                err := dataFieldsRows.Scan(&dataFields.Id, &dataFields.Type, &dataFields.Order, &dataFields.Value, &dataFields.Group)
                errorHandler(err, 500, w)

                dataFieldsList = append(dataFieldsList, dataFields)
            }

            // Трансформация группы в новый объект
            dataFieldGroupResult := &FieldGroup {
                Id:          dataFieldGroup.Id,
                Name:        dataFieldGroup.Name,
                Order:       dataFieldGroup.Order,
                Data:        dataFieldGroup.Data,
                Fields:      dataFieldsList,
            }

            content = append(content, dataFieldGroupResult)
        }

        // Трансформация Data в новый объект
        dataResult := &DataJson {
            Id:          data.Id,
            Name:        data.Name,
            Project:     data.Project,
            Parent:      data.Parent,
            Coordinates: coordinates,
            Content:     content,
        }

        dataList = append(dataList, dataResult)
	}
    errorHandler(dataRows.Err(), 500, w)

    // Подготавка выходных данных
	output, err := json.Marshal(&dataList)
	errorHandler(err, 500, w)

    // Отображение результата
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(output)
}

// Функция "Item" отображает Data объект по его id
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
	dataRow := db.QueryRow("SELECT * FROM data WHERE id = $1", id)
	data := new(DataDb)

    // Сбор данных из БД в структуру
	err := dataRow.Scan(&data.Id, &data.Name, &data.Project, &data.Parent, &data.Coordinates)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else {
        errorHandler(err, 500, w)
    }

    // Обработка координат
    var coordinates map[string]int
    json.Unmarshal([]byte(data.Coordinates), &coordinates)

    // Сбор данных из таблицы field_group связанных с Data объектом
    dataFieldGroupRows, err := db.Query("SELECT * FROM field_groups WHERE data = $1", data.Id)
    errorHandler(err, 500, w)
    defer dataFieldGroupRows.Close()

    content := make([]*FieldGroup, 0)
    for dataFieldGroupRows.Next() {
        dataFieldGroup := new(FieldGroup)

        err := dataFieldGroupRows.Scan(&dataFieldGroup.Id, &dataFieldGroup.Name, &dataFieldGroup.Order, &dataFieldGroup.Data)
        errorHandler(err, 500, w)

        // Сбор данных из таблицы fields связанных с Data объектом
        dataFieldsRows, err := db.Query("SELECT * FROM fields WHERE group_id = $1", dataFieldGroup.Id)
        errorHandler(err, 500, w)
        defer dataFieldsRows.Close()

        dataFieldsList := make([]*Field, 0)
        for dataFieldsRows.Next() {
            dataFields := new(Field)

            err := dataFieldsRows.Scan(&dataFields.Id, &dataFields.Type, &dataFields.Order, &dataFields.Value, &dataFields.Group)
            errorHandler(err, 500, w)

            dataFieldsList = append(dataFieldsList, dataFields)
        }

        // Трансформация группы в новый объект
        dataFieldGroupResult := &FieldGroup {
            Id:          dataFieldGroup.Id,
            Name:        dataFieldGroup.Name,
            Order:       dataFieldGroup.Order,
            Data:        dataFieldGroup.Data,
            Fields:      dataFieldsList,
        }

        content = append(content, dataFieldGroupResult)
    }

    // Трансформация Data в новый объект
	DataResult := &DataJson{
		Id:          data.Id,
		Name:        data.Name,
		Project:     data.Project,
		Parent:      data.Parent,
		Coordinates: coordinates,
		Content:     content,
	}

    // Подготавка выходных данных
	output, err := json.Marshal(DataResult)
	errorHandler(err, 500, w)

    // Отображение результата
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(output)
}

// Функция "Create" создает новый Data объект
func Create(w http.ResponseWriter, r *http.Request) {
    // Проверка заголовков
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

    // Сбор и анализ входных данных
	decoder := json.NewDecoder(r.Body)
	data := new(DataJson)
	err := decoder.Decode(&data)
	errorHandler(err, 500, w)
	defer r.Body.Close()

	if data.Name == "" || data.Project <= 0 {
		http.Error(w, http.StatusText(400), 400)
		return
	}

    // Подготовка запроса
	createData, err := db.Prepare("INSERT INTO data(name, project) VALUES($1, $2);")
	errorHandler(err, 500, w)

    // Выполнение запроса
	result, err := createData.Exec(data.Name, data.Project)
	errorHandler(err, 500, w)

    // Проверка на успешность
	rowsAffected, err := result.RowsAffected()
	errorHandler(err, 500, w)

    // Отображение результата
	if rowsAffected > 0 {
		fmt.Fprintf(w, "%t\n", true)
	}
}

// Функция "Update" обновляет Data объект
func Update(w http.ResponseWriter, r *http.Request) {
    // Служебные переменные
    var dataRowsAffected int64
    var fieldGroupRowsAffected int64
    var fieldRowsAffected int64

    // Проверка заголовков
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

    // Сбор и анализ входных данных
	decoder := json.NewDecoder(r.Body)
	data := new(DataJson)
	err := decoder.Decode(&data)
	errorHandler(err, 500, w)
	defer r.Body.Close()

	if data.Id == 0 || data.Project <= 0 {
		http.Error(w, http.StatusText(400), 400)
		return
	}

    // Обработка координат
	if data.Coordinates == nil {
		data.Coordinates = map[string]int{"x": 0, "y": 0}
	}
	coordinates, _ := json.Marshal(data.Coordinates)

    // Подготовка запроса
	updateData, err := db.Prepare("UPDATE data SET name = $2, project = $3, parent = $4, coordinates = $5 WHERE id = $1")
	errorHandler(err, 500, w)

    // Выполнение запроса
	dataResult, err := updateData.Exec(data.Id, data.Name, data.Project, data.Parent, coordinates)
	errorHandler(err, 500, w)

    // Проверка на успешность
	dataRowsAffected, err = dataResult.RowsAffected()
	errorHandler(err, 500, w)

    // Обработка контента Data объекта
    if data.Content != nil {
        // Подготавливаем запросы на добавление групп полей в БД
        checkFieldGroup, err := db.Prepare("SELECT * FROM field_groups WHERE id = $1")
        errorHandler(err, 500, w)

        insertFieldGroup, err := db.Prepare("INSERT INTO field_groups (name, ordr, data) VALUES ($1, $2, $3);")
        errorHandler(err, 500, w)

        updateFieldGroup, err := db.Prepare("UPDATE field_groups set name = $2, ordr = $3, data = $4 where id = $1")
        errorHandler(err, 500, w)

        // Подготавливаем запросы на добавление полей в БД
        checkField, err := db.Prepare("SELECT * FROM fields WHERE id = $1")
        errorHandler(err, 500, w)

        insertField, err := db.Prepare("INSERT INTO fields (type, ordr, value, group_id) VALUES ($1, $2, $3, $4);")
        errorHandler(err, 500, w)

        updateField, err := db.Prepare("UPDATE fields set type = $2, ordr = $3, value = $4, group_id = $5 where id = $1")
        errorHandler(err, 500, w)

        // Добавление групп полей в БД
        for g := 0; g < len(data.Content); g++ {
            checkFieldGroupResult, err := checkFieldGroup.Exec(data.Content[g].Id)
            errorHandler(err, 500, w)

            checkFieldGroupRowsAffected, err := checkFieldGroupResult.RowsAffected()
            errorHandler(err, 500, w)

            // Если этой группы нет в БД - записываем, в обратном случае - обновляем запись
            if checkFieldGroupRowsAffected == 0 {
                fieldGroupResult, err := insertFieldGroup.Exec(data.Content[g].Name, data.Content[g].Order, data.Id)
                errorHandler(err, 500, w)

                fieldGroupRowsAffected, err = fieldGroupResult.RowsAffected()
                errorHandler(err, 500, w)
            } else {
                fieldGroupResult, err := updateFieldGroup.Exec(data.Content[g].Id, data.Content[g].Name, data.Content[g].Order, data.Id)
                errorHandler(err, 500, w)

                fieldGroupRowsAffected, err = fieldGroupResult.RowsAffected()
                errorHandler(err, 500, w)
            }

            // Добавление полей в БД
            for f := 0; f < len(data.Content[g].Fields); f++ {
                checkFieldResult, err := checkField.Exec(data.Content[g].Fields[f].Id)
                errorHandler(err, 500, w)

                checkFieldRowsAffected, err := checkFieldResult.RowsAffected()
                errorHandler(err, 500, w)

                // Если этого поля нет в БД - записываем, в обратном случае - обновляем запись
                if checkFieldRowsAffected == 0 {
                    fieldResult, err := insertField.Exec(data.Content[g].Fields[f].Type, data.Content[g].Fields[f].Order, data.Content[g].Fields[f].Value, data.Content[g].Id)
                    errorHandler(err, 500, w)

                    fieldRowsAffected, err = fieldResult.RowsAffected()
                    errorHandler(err, 500, w)
                } else {
                    fieldResult, err := updateField.Exec(data.Content[g].Fields[f].Id, data.Content[g].Fields[f].Type, data.Content[g].Fields[f].Order, data.Content[g].Fields[f].Value, data.Content[g].Id)
                    errorHandler(err, 500, w)

                    fieldRowsAffected, err = fieldResult.RowsAffected()
                    errorHandler(err, 500, w)
                }
            }                    
        }
    }    

    // Отображение результата
	if dataRowsAffected > 0 || fieldGroupRowsAffected > 0 || fieldRowsAffected > 0 {
		fmt.Fprintf(w, "%t\n", true)
	}
}

// Функция "Delete" удаляет Data объект по его id и все связанные с ним данные
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

    // Удаление fields связанных с Data объектом
    deleteFields, err := db.Prepare("DELETE FROM fields f USING field_groups g WHERE g.id = f.group_id and g.data = $1")
    errorHandler(err, 500, w)

    deleteFieldsResult, err := deleteFields.Exec(id)
    errorHandler(err, 500, w)

    deleteFieldsRowsAffected, err := deleteFieldsResult.RowsAffected()
    errorHandler(err, 500, w)

    // Удаление FieldGroups связанных с Data объектом
    deleteFieldGroups, err := db.Prepare("DELETE FROM field_groups WHERE data = $1")
    errorHandler(err, 500, w)

    deleteFieldGroupsResult, err := deleteFieldGroups.Exec(id)
    errorHandler(err, 500, w)

    deleteFieldGroupsRowsAffected, err := deleteFieldGroupsResult.RowsAffected()
    errorHandler(err, 500, w)

    // Удаление Data объекта
	deleteData, err := db.Prepare("DELETE FROM data WHERE id = $1")
	errorHandler(err, 500, w)

	deleteDataResult, err := deleteData.Exec(id)
	errorHandler(err, 500, w)

	deleteDataRowsAffected, err := deleteDataResult.RowsAffected()
	errorHandler(err, 500, w)

    // Отображение результата
	if deleteDataRowsAffected > 0 || deleteFieldsRowsAffected > 0 || deleteFieldGroupsRowsAffected > 0 {
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