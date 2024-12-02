package main

import (
	myhttp "crud-app/internal/delivery/http"
	"crud-app/internal/repository"
	"crud-app/internal/usecase"
	"database/sql"
	"log/slog"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	// Настройки подключения к базе данных
	connStr := "user=postgres password=1234 dbname=crud-app sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("Failed to connect to the database", slog.String("error", err.Error()))
		return
	}
	defer db.Close()

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		slog.Error("Failed to ping the database", slog.String("error", err.Error()))
		return
	}

	// Инициализируем зависимости
	repo := repository.NewRequestRepository(db)
	uc := usecase.NewRequestUsecase(repo)
	handler := myhttp.NewHandler(uc)

	// Настраиваем маршруты
	http.HandleFunc("/", handler.Index)
	http.HandleFunc("/requests", handler.GetRequestsWithPagination)
	http.HandleFunc("/requests/create", handler.CreateRequestPage)
	http.HandleFunc("/create-request", handler.CreateRequest)
	http.HandleFunc("/requests/delete", handler.DeleteRequest)

	// Запускаем сервер
	slog.Info("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		slog.Error("Failed to start server", slog.String("error", err.Error()))
	}
}
