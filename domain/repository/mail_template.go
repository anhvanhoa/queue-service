package repository

import (
	"context"
	"queue-service/domain/common"
	"queue-service/domain/entity"
)

type MailTemplateRepository interface {
	Create(ctx context.Context, template *entity.MailTemplate) error
	GetByID(ctx context.Context, id string) (*entity.MailTemplate, error)
	GetBySubject(ctx context.Context, subject string) (*entity.MailTemplate, error)
	GetAll(ctx context.Context) ([]*entity.MailTemplate, error)
	GetByStatus(ctx context.Context, status common.Status) ([]*entity.MailTemplate, error)
	GetByProviderEmail(ctx context.Context, providerEmail string) ([]*entity.MailTemplate, error)
	Update(ctx context.Context, template *entity.MailTemplate) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status common.Status) error
}
