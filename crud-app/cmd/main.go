package main

import (
	myhttp "crud-app/internal/delivery/http"
	"crud-app/internal/repository"
	"crud-app/internal/usecase"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// Подключение к базе данных
	db, err := sqlx.Connect("postgres", "user=postgres password=1234 dbname=crud-app sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// Инициализация слоёв приложения
	repo := repository.NewRequestRepository(db)
	usecase := usecase.NewRequestUsecase(repo)
	handler := myhttp.NewHandler(usecase)

	// Настройка маршрутов
	http.HandleFunc("/", handler.Index)
	http.HandleFunc("/requests", handler.GetRequestsWithPagination)
	http.HandleFunc("/requests/create", handler.CreateRequestPage)
	http.HandleFunc("/requests/create/submit", handler.CreateRequest)
	http.HandleFunc("/requests/delete", handler.DeleteRequest)

	// Запуск сервера
	log.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
