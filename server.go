package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/labstack/echo/v4/middleware"
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
	log.Println("Create Sucess")

	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))
}
