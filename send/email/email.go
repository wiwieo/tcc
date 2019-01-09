package email

import (
	"github.com/go-gomail/gomail"
	"tcc_transaction/global/config"
	"tcc_transaction/log"
)

const (
	Port = 465
	Host = "smtp.exmail.qq.com"
)

type Email struct {
	From string
	To   []string
	//Cc []string
	Subject string
	//Attach []byte
	content chan []byte
	err     chan error
}

func NewEmailSender(from, subject string, to []string) *Email {
	e := &Email{
		From:    from,
		To:      to,
		Subject: subject,
		content: make(chan []byte, 1),
		err:     make(chan error, 1),
	}
	e.send()
	return e
}

func (e *Email) Send(content []byte) error {
	e.content <- content
	return <-e.err
}

func (e *Email) send() {
	go func() {
		d := gomail.NewDialer(Host, Port, *config.EmailUsername, *config.EmailPassword)

		var s gomail.SendCloser
		var err error
		if s, err = d.Dial(); err != nil {
			log.Errorf("connect to email failed, please check it. error info is: %s", err)
		}
		for {
			select {
			case c, ok := <-e.content:
				if !ok {
					continue
				}
				m := gomail.NewMessage()
				m.SetHeader("From", e.From)
				m.SetHeader("To", e.To...)
				m.SetHeader("Subject", e.Subject)
				m.SetBody("text/html", string(c))
				if err := gomail.Send(s, m); err != nil {
					log.Errorf("send email to [%+v] failed, please check it. error info is: %s", e.To, err)
					e.err <- err
					continue
				}
				e.err <- nil
			}
		}
	}()
}
