package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Err struct {
	Message string `json:"message"`
}

type Expense struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

func createUserHandler(c echo.Context) error {
	exp := Expense{}
	err := c.Bind(&exp)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	row := db.QueryRow("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4)  RETURNING id", exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags))
	err = row.Scan(&exp.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, exp)
}

func getUserHandler(c echo.Context) error {
	id := c.Param("id")
	stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query user statment:" + err.Error()})
	}

	row := stmt.QueryRow(id)
	exp := Expense{}
	err = row.Scan(&exp.ID, &exp.Title, &exp.Amount, &exp.Note, pq.Array(&exp.Tags))
	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "user not found"})
	case nil:
		return c.JSON(http.StatusOK, exp)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan user:" + err.Error()})
	}
}

func updateUserHandler(c echo.Context) error {
	id := c.Param("id")
	exp := Expense{}
	err := c.Bind(&exp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query user statment:" + err.Error()})
	}

	row := db.QueryRow("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id=$1 RETURNING id, title, amount, note, tags", id, exp.Title, exp.Amount, exp.Note, pq.Array(exp.Tags))

	err = row.Scan(&exp.ID, &exp.Title, &exp.Amount, &exp.Note, pq.Array(&exp.Tags))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't update:" + err.Error()})
	}
	return c.JSON(http.StatusOK, exp)
}

func getUsersHandler(c echo.Context) error {
	stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't prepare query all users statment:" + err.Error()})
	}

	rows, err := stmt.Query()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't query all users:" + err.Error()})
	}

	expense := []Expense{}

	for rows.Next() {
		exp := Expense{}
		err = rows.Scan(&exp.ID, &exp.Title, &exp.Amount, &exp.Note, pq.Array(&exp.Tags))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan user:" + err.Error()})
		}
		expense = append(expense, exp)
	}

	return c.JSON(http.StatusOK, expense)
}

var db *sql.DB

func main() {

	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	createTb := `
				CREATE TABLE IF NOT EXISTS expenses (
					id SERIAL PRIMARY KEY,
					title TEXT,
					amount FLOAT,
					note TEXT,
					tags TEXT[]
				);`

	if _, err := db.Exec(createTb); err != nil {
		log.Fatal("Cannot create table")
		return
	}
	log.Println("Create Success")

	e := echo.New()

	e.POST("/expenses", createUserHandler)
	e.GET("/expenses/:id", getUserHandler)
	e.PUT("/expenses/:id", updateUserHandler)
	e.GET("/expenses", getUsersHandler)

	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))
	log.Fatal(e.Start(os.Getenv("PORT")))
}
