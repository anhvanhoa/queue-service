package usecase

import (
	"context"
	"time"

	"github.com/anhvanhoa/service-core/common"
)

type StatusHistory struct {
	Status        string
	MailHistoryId string
	Message       string
	CreatedAt     time.Time
}

type MailTemplate struct {
	Id            string
	Name          string
	Subject       string
	Body          string
	Keys          []string
	ProviderEmail string
	Status        common.Status
	CreatedBy     string
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}

type MailProvider struct {
	Email      string
	Password   string
	UserName   string
	Port       int
	Host       string
	Encryption string
	Name       string
	TypeId     string
	CreatedBy  string
	Status     common.Status
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}

type MailHistory struct {
	ID            string
	TemplateId    string
	Subject       string
	Body          string
	Tos           []string
	Data          string
	EmailProvider string
	CreatedBy     string
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}

type MailService interface {
	CreateStatusHistory(ctx context.Context, statusHistory *StatusHistory) error
	GetMailTemplateById(ctx context.Context, id string) (*MailTemplate, error)
	GetMailProviderByEmail(ctx context.Context, email string) (*MailProvider, error)
	UpdateMailHistoryById(ctx context.Context, id string, mailHistory *MailHistory) error
}
