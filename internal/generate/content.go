package generate

type GenerateContent struct{}

func NewGenerateContent() (*GenerateContent, error) {
	return &GenerateContent{}, nil
}

func (g *GenerateContent) Run() error {
	log.Warn("not implemented")
	return nil
}
