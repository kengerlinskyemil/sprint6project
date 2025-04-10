package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// ...

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	// ...
	r.Get("/tasks", ListTasksHandler)
	r.Post("/tasks", AddTaskHandler)
	r.Get("/task/{id}", GetTaskByIdHandler)
	r.Delete("/task/{id}", DeleteTaskByIdHandler)
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}

func ListTasksHandler(w http.ResponseWriter, req *http.Request) {
	tasksMarshaled, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(tasksMarshaled)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func AddTaskHandler(w http.ResponseWriter, req *http.Request) {
	task := Task{}
	err := json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Println(err)
		return
	}

	_, ok := tasks[task.ID]
	if ok {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Println("bad id")
		return
	}

	tasks[task.ID] = task

	w.WriteHeader(http.StatusCreated)
}

func GetTaskByIdHandler(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Println("no id in url")
		return
	}

	task, ok := tasks[id]
	if !ok {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Println("no task with such id")
		return
	}

	taskMarshaled, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	_, err = w.Write(taskMarshaled)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

}
func DeleteTaskByIdHandler(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if id == "" {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Println("no id in url")
		return
	}

	_, ok := tasks[id]
	if !ok {
		http.Error(w, "bad request", http.StatusBadRequest)
		log.Println("no task with such id")
		return
	}
	delete(tasks, id)
	w.WriteHeader(http.StatusOK)
}
