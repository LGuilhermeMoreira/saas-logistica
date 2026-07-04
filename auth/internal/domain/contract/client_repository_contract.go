package contract

import "auth/internal/domain/entity"

type ClientRepositoryInterface interface {
	Create(model *entity.Client) error
	Login(email, password string) (*entity.Client, error)
}
