package config

var _ Discord = (*discord)(nil)

type Discord interface {
	GetToken() string
	GetDisableRegisterCommand() bool
}

type discord struct {
	Token string

	DisableRegisterCommand bool
}

func (d *discord) GetToken() string {
	return d.Token
}

func (d *discord) GetDisableRegisterCommand() bool {
	return d.DisableRegisterCommand
}
