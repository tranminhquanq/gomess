package usecase

type ChatUsecase struct {
}

func NewChatUsecase() *ChatUsecase {
	return &ChatUsecase{}
}

func (u *ChatUsecase) GetChatHistory() (interface{}, error) {
	return nil, nil
}

func (u *ChatUsecase) SendMessage() (interface{}, error) {
	return nil, nil
}
