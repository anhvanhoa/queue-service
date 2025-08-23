package repo

import (
	"context"
	"queue-service/domain/entity"
	repository "queue-service/domain/repository"

	"github.com/go-pg/pg/v10"
)

type mailProviderRepository struct {
	db *pg.DB
}

func NewMailProviderRepository(db *pg.DB) repository.MailProviderRepository {
	return &mailProviderRepository{
		db: db,
	}
}

func (r *mailProviderRepository) Create(ctx context.Context, provider *entity.MailProvider) error {
	db := getTx(ctx, r.db)
	_, err := db.Model(provider).Insert()
	return err
}

func (r *mailProviderRepository) GetByEmail(ctx context.Context, email string) (*entity.MailProvider, error) {
	db := getTx(ctx, r.db)
	provider := &entity.MailProvider{}
	err := db.Model(provider).Where("email = ?", email).Select()
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func (r *mailProviderRepository) GetAll(ctx context.Context) ([]*entity.MailProvider, error) {
	db := getTx(ctx, r.db)
	var providers []*entity.MailProvider
	err := db.Model(&providers).Select()
	if err != nil {
		return nil, err
	}
	return providers, nil
}

func (r *mailProviderRepository) Update(ctx context.Context, provider *entity.MailProvider) error {
	db := getTx(ctx, r.db)
	_, err := db.Model(provider).Where("email = ?", provider.Email).Update()
	return err
}

func (r *mailProviderRepository) Delete(ctx context.Context, email string) error {
	db := getTx(ctx, r.db)
	provider := &entity.MailProvider{Email: email}
	_, err := db.Model(provider).Where("email = ?", email).Delete()
	return err
}

func (r *mailProviderRepository) GetByTypeId(ctx context.Context, typeId string) ([]*entity.MailProvider, error) {
	db := getTx(ctx, r.db)
	var providers []*entity.MailProvider
	err := db.Model(&providers).Where("type_id = ?", typeId).Select()
	if err != nil {
		return nil, err
	}
	return providers, nil
}
