package usecase

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"queue-service/domain/repository"
	loggerI "queue-service/domain/service/logger"
	"queue-service/domain/service/mail"
	mailtpl "queue-service/domain/service/mail_tpl"
)

type EmailSystemImpl interface {
	SendMailQueue(ctx context.Context, Payload []byte, Id string) error
	ConfigTest() EmailTestingImpl
}

type EmailTestingImpl interface {
	SetIsProduction(mode bool) *EmailTesting
	SetIsAppedMail(mode bool) *EmailTesting
	SetTestMails(mails []string) *EmailTesting
}

type Payload struct {
	Provider string
	Tos      *[]string
	To       *string
	Template string
	Data     map[string]any
}

type EmailTesting struct {
	TestMails    []string // Danh sách email dùng để test
	IsAppedMail  bool
	IsProduction bool // Biến để xác định môi trường
}

type EmailSystem struct {
	configTest      EmailTesting
	log             loggerI.Log
	mailProvider    mail.MailProvider
	mailTemplate    mailtpl.MailTemplate
	mailTplRepo     repository.MailTemplateRepository
	mailProvierRepo repository.MailProviderRepository
}

func (e *EmailSystem) SendMailQueue(ctx context.Context, payload []byte, Id string) error {
	var pl Payload
	statusErr := make(map[string]string)

	if err := json.Unmarshal(payload, &pl); err != nil {
		statusErr["message"] = "Failed to parse payload: " + err.Error()
		// e.statusHistoryRepo.Create(&statusErr) // Sử dụng gRPC
		e.log.Warn(statusErr["message"])
		return errors.New(statusErr["message"])
	}
	tpl, err := e.mailTplRepo.GetByID(ctx, pl.Template)
	if err != nil {
		statusErr["message"] = "Failed to get mail template: " + err.Error()
		// e.statusHistoryRepo.Create(&statusErr) // Sử dụng gRPC
		return err
	} else if tpl == nil {
		statusErr["message"] = "Template not found"
		// e.statusHistoryRepo.Create(&statusErr) // Sử dụng gRPC
		return errors.New("không tìm thấy mẫu email")
	}

	mailT, err := e.mailTemplate.Render(tpl.Subject, tpl.Body, pl.Data)
	if err != nil {
		statusErr["message"] = "Failed to render mail template: " + err.Error()
		// e.statusHistoryRepo.Create(&statusErr) // Sử dụng gRPC
		return err
	}

	// Lấy thông tin cấu hình gửi email
	provider, err := e.mailProvierRepo.GetByEmail(ctx, pl.Provider)
	if err != nil {
		statusErr["message"] = "Failed to get mail provider: " + err.Error()
		// e.statusHistoryRepo.Create(&statusErr) // Sử dụng gRPC
		return err
	} else if provider == nil {
		statusErr["message"] = "Mail provider not found"
		// e.statusHistoryRepo.Create(&statusErr) // Sử dụng gRPC
		return errors.New("không tìm thấy cấu hình gửi email")
	}
	// Set cấu hình gửi email
	e.mailProvider.SetProvider(&mail.ConfigMail{
		Host:     provider.Host,
		Port:     provider.Port,
		UserName: provider.UserName,
		Password: provider.Password,
		Email:    provider.Email,
		Name:     provider.Name,
		TSL:      &tls.Config{InsecureSkipVerify: true},
	})
	tos := []string{}
	if pl.To != nil {
		tos = append(tos, *pl.To)
	} else if pl.Tos != nil {
		tos = *pl.Tos
	}
	if !e.configTest.IsProduction && len(e.configTest.TestMails) > 0 {
		tos = e.configTest.TestMails
	}
	// Gửi email
	if err := e.mailProvider.SendMail(tos, mailT.Subject, mailT.Body, pl.Data); err != nil {
		statusErr["message"] = "Failed to send mail: " + err.Error()
		// e.statusHistoryRepo.Create(&statusErr)
		return err
	}

	// Update lại subject và body vào mail history
	// err = e.mailHistoryRepo.UpdateSubAndBodyById(ctx, Id, mailT.Subject, mailT.Body) // Sử dụng gRPC
	if err != nil {
		statusErr["message"] = "Failed to update mail history: " + err.Error()
		// e.statusHistoryRepo.Create(&statusErr) // Sử dụng gRPC
		return err
	}

	// Thêm trạng thái gửi email thành công
	// statusErr.Status = entity.MAIL_STATUS_SENT
	statusErr["message"] = "Send mail success"
	// return e.statusHistoryRepo.Create(&statusErr) // Sử dụng gRPC
	return nil
}

func (e *EmailSystem) ConfigTest() EmailTestingImpl {
	return &e.configTest
}

func NewEmailSystem(
	log loggerI.Log,
	mailtemplate mailtpl.MailTemplate,
	mailProvider mail.MailProvider,
	mailTplRepo repository.MailTemplateRepository,
	mailProvierRepo repository.MailProviderRepository,
	testMails []string,
) EmailSystemImpl {
	return &EmailSystem{
		configTest:      EmailTesting{TestMails: testMails, IsAppedMail: false, IsProduction: true},
		log:             log,
		mailTemplate:    mailtemplate,
		mailProvider:    mailProvider,
		mailTplRepo:     mailTplRepo,
		mailProvierRepo: mailProvierRepo,
	}
}

func (et *EmailTesting) SetIsProduction(mode bool) *EmailTesting {
	et.IsProduction = mode
	return et
}

func (et *EmailTesting) SetIsAppedMail(mode bool) *EmailTesting {
	et.IsAppedMail = mode
	return et
}

func (et *EmailTesting) SetTestMails(mails []string) *EmailTesting {
	et.TestMails = mails
	return et
}
