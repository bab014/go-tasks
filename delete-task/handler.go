package function

import (
	"encoding/json"
	"errors"
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

var db *gorm.DB

func init() {
	// Get Connection information
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_PORT := os.Getenv("DB_PORT")
	DB_HOST := os.Getenv("DB_HOST")
	DB_NAME := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)

	// Get Database connection
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Unable to connect to database")
	}
}

func Handle(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodDelete {
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

	// Checking for Query Param of "id"
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
		var task Task
		var res Response

		// Check if ID exists
		err := db.First(&task, id[0]).Error
		if err != nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Header().Set("Content-Type", "application/json")

			res := Response{
				StatusCode: http.StatusBadRequest,
				Status:     "Bad Request",
				Data:       fmt.Sprintf("Task with id %s does not exist", id[0]),
			}

			err := json.NewEncoder(w).Encode(&res)
			if err != nil {
				w.Write([]byte("Something whent wrong"))
			}
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
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

		// Delete Task
		err = db.Delete(&task).Error
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
		res.Status = "Ok"
		res.Data = fmt.Sprintf("Task: %s deleted succesfully", id[0])

		err = json.NewEncoder(w).Encode(&res)
		if err != nil {
			w.Write([]byte("Something went wrong"))
		}
	}
}
