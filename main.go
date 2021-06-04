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

type Context struct {
	response http.ResponseWriter
	request  *http.Request
}

type Response struct {
	Error bool `json:"error"`
	User  User `json:"user"`
}

type EmptyResponse struct {
	Error bool        `json:"error"`
	User  interface{} `json:"user"`
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

func responseWriter(ctxt *Context, user []User) {
	ctxt.response.Header().Set("content-type", "application/json")
	if len(user) == 0 {
		emptyResponse := EmptyResponse{Error: true, User: ""}
		resp, err := json.Marshal(emptyResponse)
		if err != nil {
			http.Error(ctxt.response, err.Error(), http.StatusInternalServerError)
		}
		ctxt.response.Write(resp)
		return
	}
	response := Response{Error: false, User: user[0]}
	output, err := json.Marshal(response)
	if err != nil {
		http.Error(ctxt.response, err.Error(), http.StatusBadRequest)
		return
	}

	ctxt.response.Write(output)
}

func queryRunner(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	var requestBody Request
	var user []User
	context := &Context{response: w, request: r}

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
	err = db.Select(&user, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseWriter(context, user)
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
