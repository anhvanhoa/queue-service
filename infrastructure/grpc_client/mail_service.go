package grpc_client

import (
	"context"
	"queue-service/domain/usecase"
	"time"

	gc "github.com/anhvanhoa/service-core/domain/grpc_client"

	"github.com/anhvanhoa/service-core/common"
	proto_mail_history "github.com/anhvanhoa/sf-proto/gen/mail_history/v1"
	proto_mail_provider "github.com/anhvanhoa/sf-proto/gen/mail_provider/v1"
	proto_mail_template "github.com/anhvanhoa/sf-proto/gen/mail_tmpl/v1"
	proto_status_history "github.com/anhvanhoa/sf-proto/gen/status_history/v1"
)

type MailService struct {
	client *gc.Client
	shc    proto_status_history.StatusHistoryServiceClient
	mtc    proto_mail_template.MailTmplServiceClient
	mpc    proto_mail_provider.MailProviderServiceClient
	mhc    proto_mail_history.MailHistoryServiceClient
}

func NewMailService(client *gc.Client) usecase.MailService {
	if client == nil {
		return &MailService{}
	}
	shsc := proto_status_history.NewStatusHistoryServiceClient(client.GetConnection())
	mtc := proto_mail_template.NewMailTmplServiceClient(client.GetConnection())
	mpc := proto_mail_provider.NewMailProviderServiceClient(client.GetConnection())
	mhc := proto_mail_history.NewMailHistoryServiceClient(client.GetConnection())
	return &MailService{
		client: client,
		shc:    shsc,
		mtc:    mtc,
		mpc:    mpc,
		mhc:    mhc,
	}
}

func (m *MailService) CreateStatusHistory(ctx context.Context, statusHistory *usecase.StatusHistory) error {
	if err := m.validateShClient(); err != nil {
		return err
	}
	_, err := m.shc.CreateStatusHistory(ctx, &proto_status_history.CreateStatusHistoryRequest{
		Status:        statusHistory.Status,
		MailHistoryId: statusHistory.MailHistoryId,
		Message:       statusHistory.Message,
		CreatedAt:     statusHistory.CreatedAt.Format(time.RFC3339),
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *MailService) GetMailTemplateById(ctx context.Context, id string) (*usecase.MailTemplate, error) {
	if err := m.validateMtcClient(); err != nil {
		return nil, err
	}
	res, err := m.mtc.GetMailTmpl(ctx, &proto_mail_template.GetMailTmplRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return &usecase.MailTemplate{
		Id:            res.MailTmpl.Id,
		Name:          res.MailTmpl.Name,
		Keys:          res.MailTmpl.Keys,
		Subject:       res.MailTmpl.Subject,
		Body:          res.MailTmpl.Body,
		ProviderEmail: res.MailTmpl.ProviderEmail,
		Status:        common.Status(res.MailTmpl.Status),
	}, nil
}

func (m *MailService) GetMailProviderByEmail(ctx context.Context, email string) (*usecase.MailProvider, error) {
	if err := m.validateMpcClient(); err != nil {
		return nil, err
	}
	res, err := m.mpc.GetMailProvider(ctx, &proto_mail_provider.GetMailProviderRequest{
		Email: email,
	})
	if err != nil {
		return nil, err
	}

	createdAt, err := time.Parse(time.RFC3339, res.MailProvider.CreatedAt)
	if err != nil {
		return nil, err
	}
	updatedAt, err := time.Parse(time.RFC3339, res.MailProvider.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &usecase.MailProvider{
		Email:      res.MailProvider.Email,
		Password:   res.MailProvider.Password,
		UserName:   res.MailProvider.UserName,
		Port:       int(res.MailProvider.Port),
		Host:       res.MailProvider.Host,
		Encryption: res.MailProvider.Encryption,
		Name:       res.MailProvider.Name,
		TypeId:     res.MailProvider.TypeId,
		CreatedBy:  res.MailProvider.CreatedBy,
		CreatedAt:  createdAt,
		UpdatedAt:  &updatedAt,
	}, nil
}

func (m *MailService) UpdateMailHistoryById(ctx context.Context, id string, mailHistory *usecase.MailHistory) error {
	if err := m.validateMhcClient(); err != nil {
		return err
	}
	_, err := m.mhc.UpdateMailHistory(ctx, &proto_mail_history.UpdateMailHistoryRequest{
		Id:            id,
		Subject:       mailHistory.Subject,
		Body:          mailHistory.Body,
		Tos:           mailHistory.Tos,
		Data:          mailHistory.Data,
		TemplateId:    mailHistory.TemplateId,
		EmailProvider: mailHistory.EmailProvider,
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *MailService) validateShClient() error {
	if m.shc == nil {
		return ErrStatusHistoryClientNil
	}
	return nil
}
func (m *MailService) validateMtcClient() error {
	if m.mtc == nil {
		return ErrMailTemplateClientNil
	}
	return nil
}
func (m *MailService) validateMpcClient() error {
	if m.mpc == nil {
		return ErrMailProviderClientNil
	}
	return nil
}
func (m *MailService) validateMhcClient() error {
	if m.mhc == nil {
		return ErrMailHistoryClientNil
	}
	return nil
}
