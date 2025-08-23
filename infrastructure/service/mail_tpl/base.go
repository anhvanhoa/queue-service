package mailtpl

import (
	mailtpl "queue-service/domain/service/mail_tpl"

	"github.com/cbroglie/mustache"
)

type mailTemplate struct{}

func NewMailTemplate() mailtpl.MailTemplate {
	return &mailTemplate{}
}

func (m *mailTemplate) RenderWithLayout(layout, subject, body string, data map[string]any) (*mailtpl.Result, error) {
	var err error
	var result mailtpl.Result
	if subject, err = mustache.Render(subject, data); err != nil {
		return &result, err
	} else if body, err = mustache.RenderInLayout(body, layout, data); err != nil {
		return &result, err
	}
	result.Subject = subject
	result.Body = body
	return &result, nil
}

func (m *mailTemplate) Render(subject, body string, data map[string]any) (*mailtpl.Result, error) {
	var err error
	var result mailtpl.Result
	if subject, err = mustache.Render(subject, data); err != nil {
		return &result, err
	} else if body, err = mustache.Render(body, data); err != nil {
		return &result, err
	}
	result.Subject = subject
	result.Body = body
	return &result, nil
}

func (m *mailTemplate) RenderLayoutFile(fileLayout, subject, body string, data map[string]any) (*mailtpl.Result, error) {
	var result mailtpl.Result
	if layout, err := mustache.ParseFile(fileLayout); err != nil {
		return &result, err
	} else if subject, err = mustache.Render(subject, data); err != nil {
		return &result, err
	} else if renderedBody, err := layout.Render(data); err != nil {
		return &result, err
	} else {
		result.Subject = subject
		result.Body = renderedBody
	}
	return &result, nil
}
