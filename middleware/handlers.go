package middleware

import (
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"fmt"
	"log"
	"net/http" // used to access the request and response object of the api

	// used to read the environment variable
	// package used to covert string into int type
	// models package where User schema is defined

	// used to get the params from the route

	// package used to read the .env file
	_ "github.com/lib/pq" // postgres golang driver
	"github.com/usamarashid94/ToDo_ExerciseGo/models"
)

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

const (
	HOST     = "localhost"
	PORT     = 5432
	USER     = "postgres"
	PASSWORD = "1234"
	DBNAME   = "ToDo_db"
)

// create connection with postgres db
func createConnection() *sql.DB {

	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, USER, PASSWORD, DBNAME,
	)
	// Open the connection
	db, err := sql.Open("postgres", connString)

	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	// return the connection
	return db
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// get all the users in the db
	tasks, err := getAllTasksHandler()

	if err != nil {
		log.Fatalf("Unable to get all tasks. %v", err)
	}

	json.NewEncoder(w).Encode(tasks)
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// create an empty user of type models.User
	var todo models.ToDo

	// decode the json request to user
	err := json.NewDecoder(r.Body).Decode(&todo)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	// call insert user function and pass the user
	insertID := AddTaskHandler(todo)

	// format a response object
	res := response{
		ID:      insertID,
		Message: "Task added successfully",
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

//------------------------- handler functions ----------------
func AddTaskHandler(todo models.ToDo) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the insert sql query
	sqlStatement := `INSERT INTO "ToDoList" (task, id, status) VALUES ($1, $2, $3) RETURNING id`

	// the inserted id will store in this id
	var id int64

	// execute the sql statement
	// Scan function will save the insert id in the id
	err := db.QueryRow(sqlStatement, todo.Task, todo.ID, todo.Status).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	// return the inserted id
	return id
}

func getAllTasksHandler() ([]models.ToDo, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	var todos []models.ToDo

	// create the select sql query
	sqlStatement := `SELECT * FROM "ToDoList"`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var todo models.ToDo

		// unmarshal the row object to user
		err = rows.Scan(&todo.Task, &todo.ID, &todo.Status)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		todos = append(todos, todo)

	}
	// return empty user on error
	return todos, err
}
