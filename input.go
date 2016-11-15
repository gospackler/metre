package metre

type SchInput struct {
	Input     string
	RespChan  chan string
	ErrorChan chan error
}

func NewScheduleInput(inp string) *SchInput {
	return &SchInput{
		Input:     inp,
		RespChan:  make(chan string),
		ErrorChan: make(chan error),
	}
}

func (s *SchInput) Close() {
	close(s.RespChan)
	close(s.ErrorChan)
}
