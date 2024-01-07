package config

type Discord interface {
	GetToken() string
}

type discord struct {
	Token string
}

func (d *discord) GetToken() string {
	return d.Token
}
