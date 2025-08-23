package repository

import (
	"context"
	"queue-service/domain/entity"
)

type MailProviderRepository interface {
	Create(ctx context.Context, provider *entity.MailProvider) error
	GetByEmail(ctx context.Context, email string) (*entity.MailProvider, error)
	GetAll(ctx context.Context) ([]*entity.MailProvider, error)
	Update(ctx context.Context, provider *entity.MailProvider) error
	Delete(ctx context.Context, email string) error
	GetByTypeId(ctx context.Context, typeId string) ([]*entity.MailProvider, error)
}
