package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "your_password"
	dbname   = "your_dbname"
)

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Успешно подключено к базе данных!")

	url := "https://countrymeters.info/ru/World#Population_clock"

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("Ошибка: статус %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	population := doc.Find("div#cp1").Text()
	if population == "" {
		log.Fatal("Не удалось найти элемент с численностью населения")
	}

	fmt.Println("Численность населения Земли:", population)

	currentDate := time.Now().Format("02.01.2006")

	sqlStatement := `
    INSERT INTO population (population, date)
    VALUES ($1, $2)`
	_, err = db.Exec(sqlStatement, population, currentDate)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Данные успешно вставлены в базу данных!")
}

//
