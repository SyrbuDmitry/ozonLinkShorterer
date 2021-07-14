package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"ozonLinkShorterer/cmd/ozonLinkShorterer/models"
)

var dbPointer *sql.DB

//Инициалиизация базы данных
func initDataBase(name string) (*sql.DB, error) {
	database, connectionError := sql.Open("sqlite3", name)
	if connectionError != nil {
		return nil, connectionError
	}
	_, err := database.Exec(models.Scheme)
	if err != nil {
		return nil, err
	}
	fmt.Println("Init DB success")
	return database, nil
}

//Запуск сервиса
func main() {
	var errDb error
	dbPointer, errDb = initDataBase("database/urlDatabase.db")
	if errDb != nil {
		log.Fatal("Can't init database!\n", errDb)
		return
	}
	defer dbPointer.Close()
	http.HandleFunc("/short", shortenUrl)
	http.HandleFunc("/long", retrieveUrl)
	fmt.Println("Started server at", "http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
