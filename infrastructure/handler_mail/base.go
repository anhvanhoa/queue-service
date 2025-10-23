package handler_mail

import (
	"context"
	"queue-service/bootstrap"
	"queue-service/constants"
	"queue-service/domain/usecase"

	"github.com/anhvanhoa/service-core/domain/log"
	"github.com/anhvanhoa/service-core/domain/mail"
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
	log log.Logger,
	mailService usecase.MailService,
) {
	var mailS = usecase.NewEmailSystem(
		log,
		mail.NewMailTemplate(),
		mail.NewMailProvider(),
		mailService,
		[]string{"anhnguyen.xmg@xuanmaijsc.vn"},
	)
	mailS.ConfigTest().SetIsProduction(env.IsProduction())
	mux.Handle(string(constants.QUEUE_MAIL), &MailHandler{mailS})
}
