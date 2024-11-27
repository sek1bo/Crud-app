package repository

import (
	"crud-app/internal/entity"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type RequestRepository struct {
	db *sqlx.DB
}

func NewRequestRepository(db *sqlx.DB) *RequestRepository {
	return &RequestRepository{db: db}
}

func (r *RequestRepository) CreateRequest(request *entity.Request) error {
	// Логируем данные перед вставкой
	fmt.Println("Inserting request into database:", request)

	_, err := r.db.Exec("INSERT INTO requests (title, content, status) VALUES ($1, $2, $3)", request.Title, request.Content, request.Status)
	if err != nil {
		// Логируем ошибку при вставке
		fmt.Println("Error inserting request:", err)
		return err
	}
	return nil
}

func (r *RequestRepository) GetRequestsWithPagination(offset, limit int) ([]entity.Request, error) {
	var requests []entity.Request
	err := r.db.Select(&requests, "SELECT * FROM requests LIMIT $1 OFFSET $2", limit, offset)
	return requests, err
}

func (r *RequestRepository) DeleteRequest(id int) error {
	_, err := r.db.Exec("DELETE FROM requests WHERE id = $1", id)
	return err
}
