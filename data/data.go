// Package "data" contains utility functions for working with data objects.
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

func init() {
	db, err = sql.Open("postgres", "user=urivsky password=123581321 dbname=mindassistant sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

// Function "List" show list of data
func List(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        http.Error(w, http.StatusText(405), 405)
        return
    }

    // Собираем список всех Data объектов
    dataRows, err := db.Query("SELECT * FROM data")
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }
    defer dataRows.Close()

    dataList := make([]*DataJson, 0)
    for dataRows.Next() {
        data := new(DataDb)

        err := dataRows.Scan(&data.Id, &data.Name, &data.Project, &data.Parent, &data.Coordinates)
        if err != nil {
            http.Error(w, http.StatusText(500), 500)
            return
        }

        // Обрабатываем координаты
        var coordinates map[string]int
        json.Unmarshal([]byte(data.Coordinates), &coordinates)

        // Собираем список всех field_group
        dataFieldGroupRows, err := db.Query("SELECT * FROM field_groups WHERE data = $1", data.Id)
        if err != nil {
            http.Error(w, http.StatusText(500), 500)
            return
        }
        defer dataFieldGroupRows.Close()

        content := make([]*FieldGroup, 0)
        for dataFieldGroupRows.Next() {
            dataFieldGroup := new(FieldGroup)

            err := dataFieldGroupRows.Scan(&dataFieldGroup.Id, &dataFieldGroup.Name, &dataFieldGroup.Order, &dataFieldGroup.Data)
            if err != nil {
                http.Error(w, http.StatusText(500), 500)
                return
            }

            // Собираем список всех fields
            dataFieldsRows, err := db.Query("SELECT * FROM fields WHERE group_id = $1", dataFieldGroup.Id)
            if err != nil {
                http.Error(w, http.StatusText(500), 500)
                return
            }
            defer dataFieldsRows.Close()

            dataFieldsList := make([]*Field, 0)
            for dataFieldsRows.Next() {
                dataFields := new(Field)

                err := dataFieldsRows.Scan(&dataFields.Id, &dataFields.Type, &dataFields.Order, &dataFields.Value, &dataFields.Group)
                if err != nil {
                    http.Error(w, http.StatusText(500), 500)
                    return
                }

                dataFieldsList = append(dataFieldsList, dataFields)
            }

            // Собираем обработанные группы в новый объект
            dataFieldGroupResult := &FieldGroup {
                Id:          dataFieldGroup.Id,
                Name:        dataFieldGroup.Name,
                Order:       dataFieldGroup.Order,
                Data:        dataFieldGroup.Data,
                Fields:      dataFieldsList,
            }

            content = append(content, dataFieldGroupResult)
        }

        // Собираем обработанные Data в новый объект
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
    if err = dataRows.Err(); err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    // Составляем итоговый вывод списка Data объектов
    output, err := json.Marshal(dataList)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Вывод JSON на клиент
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Write(output)
}

// Function "ListByProject" show list of data by project id
func ListByProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	project := r.FormValue("project")
	if project == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

    // Собираем список Data объектов по id проекта
	dataRows, err := db.Query("SELECT * FROM data WHERE project = $1", project)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer dataRows.Close()

    dataList := make([]*DataJson, 0)
    for dataRows.Next() {
        data := new(DataDb)

        err := dataRows.Scan(&data.Id, &data.Name, &data.Project, &data.Parent, &data.Coordinates)
        if err != nil {
            http.Error(w, http.StatusText(500), 500)
            return
        }

        // Обрабатываем координаты
        var coordinates map[string]int
        json.Unmarshal([]byte(data.Coordinates), &coordinates)

        // Собираем список всех field_group
        dataFieldGroupRows, err := db.Query("SELECT * FROM field_groups WHERE data = $1", data.Id)
        if err != nil {
            http.Error(w, http.StatusText(500), 500)
            return
        }
        defer dataFieldGroupRows.Close()

        content := make([]*FieldGroup, 0)
        for dataFieldGroupRows.Next() {
            dataFieldGroup := new(FieldGroup)

            err := dataFieldGroupRows.Scan(&dataFieldGroup.Id, &dataFieldGroup.Name, &dataFieldGroup.Order, &dataFieldGroup.Data)
            if err != nil {
                http.Error(w, http.StatusText(500), 500)
                return
            }

            // Собираем список всех fields
            dataFieldsRows, err := db.Query("SELECT * FROM fields WHERE group_id = $1", dataFieldGroup.Id)
            if err != nil {
                http.Error(w, http.StatusText(500), 500)
                return
            }
            defer dataFieldsRows.Close()

            dataFieldsList := make([]*Field, 0)
            for dataFieldsRows.Next() {
                dataFields := new(Field)

                err := dataFieldsRows.Scan(&dataFields.Id, &dataFields.Type, &dataFields.Order, &dataFields.Value, &dataFields.Group)
                if err != nil {
                    http.Error(w, http.StatusText(500), 500)
                    return
                }

                dataFieldsList = append(dataFieldsList, dataFields)
            }

            // Собираем обработанные группы в новый объект
            dataFieldGroupResult := &FieldGroup {
                Id:          dataFieldGroup.Id,
                Name:        dataFieldGroup.Name,
                Order:       dataFieldGroup.Order,
                Data:        dataFieldGroup.Data,
                Fields:      dataFieldsList,
            }

            content = append(content, dataFieldGroupResult)
        }

        // Собираем обработанные Data в новый объект
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
	if err = dataRows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	output, err := json.Marshal(&dataList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(output)
}

// Function "Item" show data by id
func Item(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	dataRow := db.QueryRow("SELECT * FROM data WHERE id = $1", id)
	data := new(DataDb)

	err := dataRow.Scan(&data.Id, &data.Name, &data.Project, &data.Parent, &data.Coordinates)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

    // Обрабатываем координаты
    var coordinates map[string]int
    json.Unmarshal([]byte(data.Coordinates), &coordinates)

    // Собираем список всех field_group
    dataFieldGroupRows, err := db.Query("SELECT * FROM field_groups WHERE data = $1", data.Id)
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }
    defer dataFieldGroupRows.Close()

    content := make([]*FieldGroup, 0)
    for dataFieldGroupRows.Next() {
        dataFieldGroup := new(FieldGroup)

        err := dataFieldGroupRows.Scan(&dataFieldGroup.Id, &dataFieldGroup.Name, &dataFieldGroup.Order, &dataFieldGroup.Data)
        if err != nil {
            http.Error(w, http.StatusText(500), 500)
            return
        }

        // Собираем список всех fields
        dataFieldsRows, err := db.Query("SELECT * FROM fields WHERE group_id = $1", dataFieldGroup.Id)
        if err != nil {
            http.Error(w, http.StatusText(500), 500)
            return
        }
        defer dataFieldsRows.Close()

        dataFieldsList := make([]*Field, 0)
        for dataFieldsRows.Next() {
            dataFields := new(Field)

            err := dataFieldsRows.Scan(&dataFields.Id, &dataFields.Type, &dataFields.Order, &dataFields.Value, &dataFields.Group)
            if err != nil {
                http.Error(w, http.StatusText(500), 500)
                return
            }

            dataFieldsList = append(dataFieldsList, dataFields)
        }

        // Собираем обработанные группы в новый объект
        dataFieldGroupResult := &FieldGroup {
            Id:          dataFieldGroup.Id,
            Name:        dataFieldGroup.Name,
            Order:       dataFieldGroup.Order,
            Data:        dataFieldGroup.Data,
            Fields:      dataFieldsList,
        }

        content = append(content, dataFieldGroupResult)
    }

	DataResult := &DataJson{
		Id:          data.Id,
		Name:        data.Name,
		Project:     data.Project,
		Parent:      data.Parent,
		Coordinates: coordinates,
		Content:     content,
	}

	output, err := json.Marshal(DataResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(output)
}

// Function "Create" creates a new data object by json
func Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	decoder := json.NewDecoder(r.Body)
	data := new(DataJson)
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if data.Name == "" || data.Project <= 0 {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	stmt, err := db.Prepare("INSERT INTO data(name, project) VALUES($1, $2);")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	result, err := stmt.Exec(data.Name, data.Project)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if rowsAffected > 0 {
		fmt.Fprintf(w, "%t\n", true)
	}
}

// Function "Update" updates a new data by json
func Update(w http.ResponseWriter, r *http.Request) {
    var dataRowsAffected int64
    var fieldGroupRowsAffected int64
    var fieldRowsAffected int64

	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	decoder := json.NewDecoder(r.Body)
	data := new(DataJson)
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if data.Id == 0 || data.Project <= 0 {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	if data.Coordinates == nil {
		data.Coordinates = map[string]int{"x": 0, "y": 0}
	}
	coordinates, _ := json.Marshal(data.Coordinates)

	updateData, err := db.Prepare("UPDATE data SET name = $2, project = $3, parent = $4, coordinates = $5 WHERE id = $1")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	dataResult, err := updateData.Exec(data.Id, data.Name, data.Project, data.Parent, coordinates)
	if err != nil {
		log.Fatal(err)
	}

	dataRowsAffected, err = dataResult.RowsAffected()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

    if data.Content != nil {
        // Подготавливаем запросы на добавление групп полей в БД
        checkFieldGroup, err := db.Prepare("SELECT * FROM field_groups WHERE id = $1")
        if err != nil {
            http.Error(w, http.StatusText(500), 500)
            return
        }

        insertFieldGroup, err := db.Prepare("INSERT INTO field_groups (name, ordr, data) VALUES ($1, $2, $3);")
        if err != nil {
            http.Error(w, http.StatusText(500), 500)
            return
        }

        updateFieldGroup, err := db.Prepare("UPDATE field_groups set name = $2, ordr = $3, data = $4 where id = $1")
        if err != nil {
            http.Error(w, http.StatusText(500), 500)
            return
        }

        // Подготавливаем запросы на добавление полей в БД
        checkField, err := db.Prepare("SELECT * FROM fields WHERE id = $1")
        if err != nil {
            http.Error(w, http.StatusText(500), 500)
            return
        }

        insertField, err := db.Prepare("INSERT INTO fields (type, ordr, value, group_id) VALUES ($1, $2, $3, $4);")
        if err != nil {
            http.Error(w, http.StatusText(500), 500)
            return
        }     

        updateField, err := db.Prepare("UPDATE fields set type = $2, ordr = $3, value = $4, group_id = $5 where id = $1")
        if err != nil {
            http.Error(w, http.StatusText(500), 500)
            return
        }   

        // Добавление групп полей в БД
        for g := 0; g < len(data.Content); g++ {
            checkFieldGroupResult, err := checkFieldGroup.Exec(data.Content[g].Id)
            if err != nil {
                log.Fatal(err)
            }

            checkFieldGroupRowsAffected, err := checkFieldGroupResult.RowsAffected()
            if err != nil {
                http.Error(w, http.StatusText(500), 500)
                return
            }

            // Если этой группы нет в БД - записываем, в обратном случае - обновляем запись
            if checkFieldGroupRowsAffected == 0 {
                fieldGroupResult, err := insertFieldGroup.Exec(data.Content[g].Name, data.Content[g].Order, data.Id)
                if err != nil {
                    log.Fatal(err)
                }

                fieldGroupRowsAffected, err = fieldGroupResult.RowsAffected()
                if err != nil {
                    http.Error(w, http.StatusText(500), 500)
                    return
                }
            } else {
                fieldGroupResult, err := updateFieldGroup.Exec(data.Content[g].Id, data.Content[g].Name, data.Content[g].Order, data.Id)
                if err != nil {
                    log.Fatal(err)
                }

                fieldGroupRowsAffected, err = fieldGroupResult.RowsAffected()
                if err != nil {
                    http.Error(w, http.StatusText(500), 500)
                    return
                }
            }

            // Добавление полей в БД
            for f := 0; f < len(data.Content[g].Fields); f++ {
                checkFieldResult, err := checkField.Exec(data.Content[g].Fields[f].Id)
                if err != nil {
                    log.Fatal(err)
                }

                checkFieldRowsAffected, err := checkFieldResult.RowsAffected()
                if err != nil {
                    http.Error(w, http.StatusText(500), 500)
                    return
                }

                // Если этого поля нет в БД - записываем, в обратном случае - обновляем запись
                if checkFieldRowsAffected == 0 {
                    fieldResult, err := insertField.Exec(data.Content[g].Fields[f].Type, data.Content[g].Fields[f].Order, data.Content[g].Fields[f].Value, data.Content[g].Id)
                    if err != nil {
                        log.Fatal(err)
                    }

                    fieldRowsAffected, err = fieldResult.RowsAffected()
                    if err != nil {
                        http.Error(w, http.StatusText(500), 500)
                        return
                    }
                } else {
                    fieldResult, err := updateField.Exec(data.Content[g].Fields[f].Id, data.Content[g].Fields[f].Type, data.Content[g].Fields[f].Order, data.Content[g].Fields[f].Value, data.Content[g].Id)
                    if err != nil {
                        log.Fatal(err)
                    }

                    fieldRowsAffected, err = fieldResult.RowsAffected()
                    if err != nil {
                        http.Error(w, http.StatusText(500), 500)
                        return
                    }
                }
            }                    
        }
    }    

	if dataRowsAffected > 0 || fieldGroupRowsAffected > 0 || fieldRowsAffected > 0 {
		fmt.Fprintf(w, "%t\n", true)
	}
}

// Function "Delete" delete a data by id
func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

    // Удаляем Fields объекта
    deleteFields, err := db.Prepare("DELETE FROM fields f USING field_groups g WHERE g.id = f.group_id and g.data = $1")
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    deleteFieldsResult, err := deleteFields.Exec(id)
    if err != nil {
        log.Fatal(err)
    }

    deleteFieldsRowsAffected, err := deleteFieldsResult.RowsAffected()
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    // Удаляем FieldGroups объекта
    deleteFieldGroups, err := db.Prepare("DELETE FROM field_groups WHERE data = $1")
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    deleteFieldGroupsResult, err := deleteFieldGroups.Exec(id)
    if err != nil {
        log.Fatal(err)
    }

    deleteFieldGroupsRowsAffected, err := deleteFieldGroupsResult.RowsAffected()
    if err != nil {
        http.Error(w, http.StatusText(500), 500)
        return
    }

    // Удаляем Data объект
	deleteData, err := db.Prepare("DELETE FROM data WHERE id = $1")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	deleteDataResult, err := deleteData.Exec(id)
	if err != nil {
		log.Fatal(err)
	}

	deleteDataRowsAffected, err := deleteDataResult.RowsAffected()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if deleteDataRowsAffected > 0 || deleteFieldsRowsAffected > 0 || deleteFieldGroupsRowsAffected > 0 {
		fmt.Fprintf(w, "%t\n", true)
	}
}
