package mailtpl

type Result struct {
	Subject string
	Body    string
}

type MailTemplate interface {
	RenderWithLayout(layout, subject, body string, data map[string]any) (*Result, error)
	Render(subject, body string, data map[string]any) (*Result, error)
	RenderLayoutFile(fileLayout, subject, body string, data map[string]any) (*Result, error)
}
