package usecase

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"queue-service/constants"
	serviceError "queue-service/domain/service/error"
	loggerI "queue-service/domain/service/logger"
	"queue-service/domain/service/mail"
	mailtpl "queue-service/domain/service/mail_tpl"
)

var (
	ErrMailTemplateNotFound    = serviceError.NewErr("Không tìm thấy mẫu email")
	ErrMailProviderNotFound    = serviceError.NewErr("Không tìm thấy cấu hình gửi email")
	ErrMailSendFailed          = serviceError.NewErr("Không thể gửi email")
	ErrMailUpdateHistoryFailed = serviceError.NewErr("Không thể cập nhật lịch sử email")
	ErrMailParseFailed         = serviceError.NewErr("Không thể phân tích dữ liệu payload")
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
	configTest   EmailTesting
	log          loggerI.Log
	mailProvider mail.MailProvider
	mailTemplate mailtpl.MailTemplate
	mailService  MailService
}

func (e *EmailSystem) SendMailQueue(ctx context.Context, payload []byte, Id string) error {
	var pl Payload
	sh := &StatusHistory{
		Status: constants.STATUS_SENT_MAIL_PENDING.String(),
	}

	if err := json.Unmarshal(payload, &pl); err != nil {
		sh.Message = ErrMailParseFailed.Error() + ": " + err.Error()
		e.mailService.CreateStatusHistory(ctx, sh)
		e.log.Warn(sh.Message)
		return ErrMailParseFailed
	}
	tpl, err := e.mailService.GetMailTemplateById(ctx, pl.Template)
	if err != nil {
		sh.Message = ErrMailTemplateNotFound.Error() + ": " + err.Error()
		e.log.Warn(sh.Message)
		e.mailService.CreateStatusHistory(ctx, sh)
		return ErrMailTemplateNotFound
	} else if tpl == nil {
		sh.Message = ErrMailTemplateNotFound.Error()
		e.log.Warn(sh.Message)
		e.mailService.CreateStatusHistory(ctx, sh)
		return ErrMailTemplateNotFound
	}

	mailT, err := e.mailTemplate.Render(tpl.Subject, tpl.Body, pl.Data)
	if err != nil {
		sh.Message = ErrMailTemplateNotFound.Error() + ": " + err.Error()
		e.log.Warn(sh.Message)
		e.mailService.CreateStatusHistory(ctx, sh)
		return ErrMailTemplateNotFound
	}

	provider, err := e.mailService.GetMailProviderByEmail(ctx, pl.Provider)
	if err != nil {
		sh.Message = ErrMailProviderNotFound.Error() + ": " + err.Error()
		e.log.Warn(sh.Message)
		e.mailService.CreateStatusHistory(ctx, sh)
		return ErrMailProviderNotFound
	} else if provider == nil {
		sh.Message = ErrMailProviderNotFound.Error()
		e.log.Warn(sh.Message)
		e.mailService.CreateStatusHistory(ctx, sh)
		return ErrMailProviderNotFound
	}

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

	if err := e.mailProvider.SendMail(tos, mailT.Subject, mailT.Body, pl.Data); err != nil {
		sh.Message = ErrMailSendFailed.Error() + ": " + err.Error()
		e.log.Warn(sh.Message)
		e.mailService.CreateStatusHistory(ctx, sh)
		return ErrMailSendFailed
	}

	err = e.mailService.UpdateMailHistoryById(ctx, Id, &MailHistory{
		Subject: mailT.Subject,
		Body:    mailT.Body,
	})
	if err != nil {
		sh.Message = ErrMailUpdateHistoryFailed.Error() + ": " + err.Error()
		e.log.Warn(sh.Message)
		e.mailService.CreateStatusHistory(ctx, sh)
		return ErrMailUpdateHistoryFailed
	}

	sh.Status = constants.STATUS_SENT_MAIL_SENT.String()
	sh.Message = "Gửi email thành công"
	e.log.Info(sh.Message)
	e.mailService.CreateStatusHistory(ctx, sh)
	return nil
}

func (e *EmailSystem) ConfigTest() EmailTestingImpl {
	return &e.configTest
}

func NewEmailSystem(
	log loggerI.Log,
	mailtemplate mailtpl.MailTemplate,
	mailProvider mail.MailProvider,
	mailService MailService,
	testMails []string,
) EmailSystemImpl {
	return &EmailSystem{
		configTest:   EmailTesting{TestMails: testMails, IsAppedMail: false, IsProduction: true},
		log:          log,
		mailTemplate: mailtemplate,
		mailProvider: mailProvider,
		mailService:  mailService,
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
