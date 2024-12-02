package repository

import (
	"crud-app/internal/entity"
	"database/sql"
	"fmt"
	"log/slog"
)

type RequestRepository struct {
	db *sql.DB
}

func NewRequestRepository(db *sql.DB) *RequestRepository {
	return &RequestRepository{db: db}
}

// Создание новой заявки с использованием транзакции
func (r *RequestRepository) CreateRequest(request *entity.Request) error {
	tx, err := r.db.Begin()
	if err != nil {
		slog.Error("Failed to begin transaction", slog.String("error", err.Error()))
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.Exec("INSERT INTO requests (title, content, status) VALUES ($1, $2, $3)", request.Title, request.Content, request.Status)
	if err != nil {
		slog.Error("Failed to insert request", slog.String("error", err.Error()))
		tx.Rollback()
		return fmt.Errorf("failed to insert request: %w", err)
	}

	if err := tx.Commit(); err != nil {
		slog.Error("Failed to commit transaction", slog.String("error", err.Error()))
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	slog.Info("Request created successfully", slog.String("title", request.Title))
	return nil
}

// Удаление заявки
func (r *RequestRepository) DeleteRequest(id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		slog.Error("Failed to begin transaction", slog.String("error", err.Error()))
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.Exec("DELETE FROM requests WHERE id = $1", id)
	if err != nil {
		slog.Error("Failed to delete request", slog.String("error", err.Error()))
		tx.Rollback()
		return fmt.Errorf("failed to delete request: %w", err)
	}

	if err := tx.Commit(); err != nil {
		slog.Error("Failed to commit transaction", slog.String("error", err.Error()))
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	slog.Info("Request deleted successfully", slog.Int("id", id))
	return nil
}

// Получение заявок с сортировкой
func (r *RequestRepository) GetRequestsWithPagination(offset, limit int) ([]entity.Request, error) {
	query := `
		SELECT id, title, content, status 
		FROM requests 
		ORDER BY id DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		slog.Error("Failed to retrieve requests", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to retrieve requests: %w", err)
	}
	defer rows.Close()

	var requests []entity.Request
	for rows.Next() {
		var request entity.Request
		if err := rows.Scan(&request.ID, &request.Title, &request.Content, &request.Status); err != nil {
			slog.Error("Failed to scan request", slog.String("error", err.Error()))
			return nil, fmt.Errorf("failed to scan request: %w", err)
		}
		requests = append(requests, request)
	}

	slog.Info("Requests retrieved successfully", slog.Int("count", len(requests)))
	return requests, nil
}
