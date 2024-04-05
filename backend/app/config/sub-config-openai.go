package config

type Openai interface {
	GetToken() string
}

var _ Openai = (*openai)(nil)

type openai struct {
	Token string
}

func (o *openai) GetToken() string {
	return o.Token
}
