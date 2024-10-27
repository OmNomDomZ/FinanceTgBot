package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v\n", err)
	}
	defer db.Close()

	if err := goose.Up(db, "internal/db/migrations"); err != nil {
		log.Fatalf("Ошибка применения миграций: %v\n", err)
	}

	fmt.Println("Миграции успешно применены!")
}
