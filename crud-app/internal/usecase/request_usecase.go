package usecase

import (
    "crud-app/internal/entity"
    "crud-app/internal/repository"
)

type RequestUsecase struct {
    repo *repository.RequestRepository
}

func NewRequestUsecase(repo *repository.RequestRepository) *RequestUsecase {
    return &RequestUsecase{repo: repo}
}

func (u *RequestUsecase) CreateRequest(request *entity.Request) error {
    return u.repo.CreateRequest(request)
}

func (u *RequestUsecase) GetRequestsWithPagination(offset, limit int) ([]entity.Request, error) {
    return u.repo.GetRequestsWithPagination(offset, limit)
}

func (u *RequestUsecase) DeleteRequest(id int) error {
    return u.repo.DeleteRequest(id)
}

func (u *RequestUsecase) CleanOldRequests() error {
    return u.repo.CleanOldRequests()
}
