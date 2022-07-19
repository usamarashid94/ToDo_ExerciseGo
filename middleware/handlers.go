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
	// load .env file
	//err := godotenv.Load(".env")
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, USER, PASSWORD, DBNAME,
	)

	// if err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }

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
	users, err := getAllTasks()

	if err != nil {
		log.Fatalf("Unable to get all user. %v", err)
	}

	// send all the users as response
	json.NewEncoder(w).Encode(users)
}

func getAllTasks() ([]models.ToDo, error) {
	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	var users []models.ToDo

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
		var user models.ToDo

		// unmarshal the row object to user
		err = rows.Scan(&user.Task, &user.ID, &user.Status)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// append the user in the users slice
		users = append(users, user)

	}

	// return empty user on error
	return users, err
}
