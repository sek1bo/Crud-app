package main

import (
	"context"
	myhttp "crud-app/internal/delivery/http"
	"crud-app/internal/repository"
	"crud-app/internal/usecase"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
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
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Index)
	mux.HandleFunc("/requests", handler.GetRequestsWithPagination)
	mux.HandleFunc("/requests/create", handler.CreateRequestPage)
	mux.HandleFunc("/create-request", handler.CreateRequest)
	mux.HandleFunc("/requests/delete", handler.DeleteRequest)

	// Создаем сервер (для последующего graceful shutdown)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Инициализируем cron
	c := cron.New(cron.WithSeconds())
	_, err = c.AddFunc("0 0 0 * * *", func() {
		slog.Info("Starting old requests cleanup...")
		if err := uc.CleanOldRequests(); err != nil {
			slog.Error("Failed to clean old requests", slog.String("error", err.Error()))
		}
	})
	if err != nil {
		slog.Error("Failed to schedule cron job", slog.String("error", err.Error()))
		return
	}

	// Запускаем cron
	c.Start()

	go func() {
		slog.Info("Starting server on :8080...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server ListenAndServe error", slog.String("error", err.Error()))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	slog.Info("Stopping cron...")
	c.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	slog.Info("Shutting down server gracefully...")
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", slog.String("error", err.Error()))
	} else {
		slog.Info("Server shutdown gracefully")
	}
}
