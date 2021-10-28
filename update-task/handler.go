package function

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title    string `json:"title"`
	Complete bool   `json:"complete"`
}

type Response struct {
	StatusCode int         `json:"status_code"`
	Status     string      `json:"status"`
	Data       interface{} `json:"data"`
}

func Handle(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")

		res := Response{
			StatusCode: http.StatusMethodNotAllowed,
			Status:     "Method not allowed",
			Data:       "Error",
		}

		err := json.NewEncoder(w).Encode(&res)
		if err != nil {
			w.Write([]byte("Something whent wrong"))
		}
		return
	}

	_, ok := r.URL.Query()["id"]
	if !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")

		res := Response{
			StatusCode: http.StatusBadRequest,
			Status:     "Bad Request",
			Data:       "No ID provided",
		}

		err := json.NewEncoder(w).Encode(&res)
		if err != nil {
			w.Write([]byte("Something whent wrong"))
		}
		return
	} else {
		id := r.URL.Query()["id"]

		var updateTask Task
		var res Response

		err := json.NewDecoder(r.Body).Decode(&updateTask)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")

			res.StatusCode = http.StatusBadRequest
			res.Status = "Bad request"
			res.Data = "Error"
			err = json.NewEncoder(w).Encode(&res)
			if err != nil {
				w.Write([]byte("Something went wrong"))
			}
			return
		}

		// Get Connection information
		DB_USER := os.Getenv("DB_USER")
		DB_PASSWORD := os.Getenv("DB_PASSWORD")
		DB_PORT := os.Getenv("DB_PORT")
		DB_HOST := os.Getenv("DB_HOST")
		DB_NAME := os.Getenv("DB_NAME")

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)

		// Get Database connection
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Println("Unable to connect to database")
		}

		var task Task

		err = db.First(&task, id[0]).Error
		if err != nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Header().Set("Content-Type", "application/json")

			res := Response{
				StatusCode: http.StatusBadRequest,
				Status:     "Bad Request",
				Data:       "Task does not exist",
			}

			err := json.NewEncoder(w).Encode(&res)
			if err != nil {
				w.Write([]byte("Something whent wrong"))
			}
			return
		}

		if updateTask.Title != "" {
			task.Title = updateTask.Title
		}

		if updateTask.Complete {
			task.Complete = updateTask.Complete
		}

		err = db.Save(&task).Error
		if err != nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Header().Set("Content-Type", "application/json")

			res := Response{
				StatusCode: http.StatusInternalServerError,
				Status:     "Internal Server Error",
				Data:       err,
			}

			err := json.NewEncoder(w).Encode(&res)
			if err != nil {
				w.Write([]byte("Something whent wrong"))
			}
			return
		}

		res.StatusCode = http.StatusOK
		res.Status = fmt.Sprintf("Task: %s updated", id)
		res.Data = task

		err = json.NewEncoder(w).Encode(&res)
		if err != nil {
			w.Write([]byte("Something went wrong"))
		}

	}
}
