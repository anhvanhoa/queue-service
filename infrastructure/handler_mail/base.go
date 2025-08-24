package handler_mail

import (
	"context"
	"queue-service/bootstrap"
	"queue-service/constants"
	loggerI "queue-service/domain/service/logger"
	"queue-service/domain/usecase"
	"queue-service/infrastructure/service/mail"
	mailtpl "queue-service/infrastructure/service/mail_tpl"

	"github.com/hibiken/asynq"
)

type MailHandler struct {
	mailSend usecase.EmailSystemImpl
}

func (e *MailHandler) ProcessTask(ctx context.Context, task *asynq.Task) error {
	return e.mailSend.SendMailQueue(ctx, task.Payload(), task.ResultWriter().TaskID())
}

func NewEmailHandler(
	mux *asynq.ServeMux,
	env *bootstrap.Env,
	log loggerI.Log,
	mailService usecase.MailService,
) {
	var mailS = usecase.NewEmailSystem(
		log,
		mailtpl.NewMailTemplate(),
		mail.NewMailProvider(),
		mailService,
		[]string{"anhnguyen.xmg@xuanmaijsc.vn"},
	)
	mailS.ConfigTest().SetIsProduction(env.IsProduction())
	mux.Handle(string(constants.QUEUE_MAIL), &MailHandler{mailS})
}
