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

var db *gorm.DB

func init() {
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_PORT := os.Getenv("DB_PORT")
	DB_HOST := os.Getenv("DB_HOST")
	DB_NAME := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)

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

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Content-Type", "application/json")

		res := Response{
			StatusCode: http.StatusMethodNotAllowed,
			Status:     "Method not allowed",
			Data:       "Error",
		}

		err := json.NewEncoder(w).Encode(&res)
		if err != nil {
			w.Write([]byte("Something went wrong"))
		}
		return
	}

	var res Response
	var newTask Task

	err := json.NewDecoder(r.Body).Decode(&newTask)
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

	err = db.Create(&newTask).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")

		res.StatusCode = http.StatusInternalServerError
		res.Status = "Internal server error"
		res.Data = "Error"

		err = json.NewEncoder(w).Encode(&res)
		if err != nil {
			w.Write([]byte("Something went wrong"))
		}
		return
	}

	res.StatusCode = http.StatusOK
	res.Status = "All good, new task created"
	res.Data = newTask

	err = json.NewEncoder(w).Encode(&res)
	if err != nil {
		w.Write([]byte("Something went wrong"))
	}

}
