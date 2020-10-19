package message

type Massage struct {
	From    string
	To      string
	Content string
}

type MassageService interface {
	Send(massage Massage) error
}

type SimpleMessageService struct {}

func (sms *SimpleMessageService) Send(message Massage) {

}