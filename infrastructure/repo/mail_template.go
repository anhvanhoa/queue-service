package repo

import (
	"context"
	"queue-service/domain/common"
	"queue-service/domain/entity"
	repository "queue-service/domain/repository"

	"github.com/go-pg/pg/v10"
)

type mailTemplateRepository struct {
	db *pg.DB
}

func NewMailTemplateRepository(db *pg.DB) repository.MailTemplateRepository {
	return &mailTemplateRepository{
		db: db,
	}
}

func (r *mailTemplateRepository) Create(ctx context.Context, template *entity.MailTemplate) error {
	db := getTx(ctx, r.db)
	_, err := db.Model(template).Insert()
	return err
}

func (r *mailTemplateRepository) GetByID(ctx context.Context, id string) (*entity.MailTemplate, error) {
	db := getTx(ctx, r.db)
	template := &entity.MailTemplate{}
	err := db.Model(template).Where("id = ?", id).Select()
	if err != nil {
		return nil, err
	}
	return template, nil
}

func (r *mailTemplateRepository) GetBySubject(ctx context.Context, subject string) (*entity.MailTemplate, error) {
	db := getTx(ctx, r.db)
	template := &entity.MailTemplate{}
	err := db.Model(template).Where("subject = ?", subject).Select()
	if err != nil {
		return nil, err
	}
	return template, nil
}

func (r *mailTemplateRepository) GetAll(ctx context.Context) ([]*entity.MailTemplate, error) {
	db := getTx(ctx, r.db)
	var templates []*entity.MailTemplate
	err := db.Model(&templates).Select()
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *mailTemplateRepository) GetByStatus(ctx context.Context, status common.Status) ([]*entity.MailTemplate, error) {
	db := getTx(ctx, r.db)
	var templates []*entity.MailTemplate
	err := db.Model(&templates).Where("status = ?", status).Select()
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *mailTemplateRepository) GetByProviderEmail(ctx context.Context, providerEmail string) ([]*entity.MailTemplate, error) {
	db := getTx(ctx, r.db)
	var templates []*entity.MailTemplate
	err := db.Model(&templates).Where("provider_email = ?", providerEmail).Select()
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *mailTemplateRepository) Update(ctx context.Context, template *entity.MailTemplate) error {
	db := getTx(ctx, r.db)
	_, err := db.Model(template).Where("id = ?", template.ID).Update()
	return err
}

func (r *mailTemplateRepository) Delete(ctx context.Context, id string) error {
	db := getTx(ctx, r.db)
	template := &entity.MailTemplate{ID: id}
	_, err := db.Model(template).Where("id = ?", id).Delete()
	return err
}

func (r *mailTemplateRepository) UpdateStatus(ctx context.Context, id string, status common.Status) error {
	db := getTx(ctx, r.db)
	_, err := db.Model((*entity.MailTemplate)(nil)).
		Set("status = ?", status).
		Where("id = ?", id).
		Update()
	return err
}
