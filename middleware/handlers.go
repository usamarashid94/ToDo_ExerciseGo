package middleware

import (
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"fmt"
	"log"
	"net/http" // used to access the request and response object of the api
	"strconv"

	// used to read the environment variable
	// package used to covert string into int type
	// models package where User schema is defined

	// used to get the params from the route

	// package used to read the .env file
	"github.com/gorilla/mux"
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

//------------------------- handlers ------------------------

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

func DeleteTask(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	// convert the id in string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	deletedRows := deleteTaskHandler(int64(id))

	// format the message string
	msg := fmt.Sprintf("Tasks updated successfully. Total rows/record affected %v", deletedRows)

	// format the reponse message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

// UpdateUser update user's detail in the postgres db
func UpdateTask(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// get the id from the request params, key is "id"
	params := mux.Vars(r)

	// convert the id type from string to int
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	var todo models.ToDo

	err = json.NewDecoder(r.Body).Decode(&todo)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	updatedRows := updateTaskHandler(int64(id), todo)

	// format the message string
	msg := fmt.Sprintf("Task updated successfully. Total rows/record affected %v", updatedRows)

	// format the response message
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// send the response
	json.NewEncoder(w).Encode(res)
}

//------------------------- handler functions ----------------

func deleteTaskHandler(id int64) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE FROM "ToDoList" WHERE id=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

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

func updateTaskHandler(id int64, todo models.ToDo) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE "ToDoList" SET task=$2, status=$3 WHERE id=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id, todo.Task, todo.Status)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
