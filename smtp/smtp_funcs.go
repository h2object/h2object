package smtp

func (helper *SMTPHelper) Process(to []string, subject, template string, data map[string]interface{}) error {
	message := NewMessage(to, subject)
	if err := message.RenderTemplate("html", template, data); err != nil {
		return err
	}
	return helper.SendMessage(message)
}