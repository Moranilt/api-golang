package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Request struct {
	Name string `json:"name"`
}

type User struct {
	Id        int64  `db:"id"`
	Name      string `db:"name"`
	Phone     string `db:"phone"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func checkMethod(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return errors.New("Only method POST allowed")
	}

	return nil
}

func queryRunner(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	var requestBody Request
	var user []User

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error while reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE name='%v'", requestBody.Name)
	db.Select(&user, query)

	output, err := json.Marshal(user[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

func mainHandler(fn func(http.ResponseWriter, *http.Request, *sqlx.DB)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := checkMethod(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
			return
		}
		db, err := sqlx.Connect("mysql", "root:123@/api-golang")
		if err != nil {
			log.Fatalln(err)
			return
		}
		fn(w, r, db)
	}
}

func main() {
	http.HandleFunc("/", mainHandler(queryRunner))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
