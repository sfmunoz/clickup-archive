package generate

type GenerateContent struct {
	clickupDir string
}

func NewGenerateContent(clickupDir string) (*GenerateContent, error) {
	return &GenerateContent{clickupDir: clickupDir}, nil
}

func (g *GenerateContent) Run() error {
	log.Warn("not implemented", "clickupDir", g.clickupDir)
	return nil
}
