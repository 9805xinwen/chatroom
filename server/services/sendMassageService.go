package services

type Massage struct {
	From    string
	To      string
	Content string
}

type SendMassageService interface {
	Send(massage Massage) error
}
