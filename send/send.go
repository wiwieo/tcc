package send

type Send interface {
	Send(content []byte) error
}

