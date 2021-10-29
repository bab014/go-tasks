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

	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodGet {
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

	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_PORT := os.Getenv("DB_PORT")
	DB_HOST := os.Getenv("DB_HOST")
	DB_NAME := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Unable to connect to database")
	}

	var tasks []Task

	id, ok := r.URL.Query()["id"]
	if !ok {
		err = db.Find(&tasks).Error
		if err != nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Header().Set("Content-Type", "application/json")

			res := Response{
				StatusCode: http.StatusInternalServerError,
				Status:     "Internal Server Error",
				Data:       "Error",
			}

			err := json.NewEncoder(w).Encode(&res)
			if err != nil {
				w.Write([]byte("Something whent wrong"))
			}
			return
		}

		var res Response

		res.StatusCode = http.StatusOK
		res.Status = "Ok"
		res.Data = tasks

		err = json.NewEncoder(w).Encode(&res)
		if err != nil {
			w.Write([]byte("something went wrong"))
		}
	} else {
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
		var res Response

		res.StatusCode = http.StatusOK
		res.Status = "Ok"
		res.Data = task

		err = json.NewEncoder(w).Encode(&res)
		if err != nil {
			w.Write([]byte("something went wrong"))
		}

	}

}
