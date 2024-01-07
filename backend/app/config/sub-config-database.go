package config

type Redis interface {
	GetAddr() string
	GetUsername() string
	GetPassword() string
	GetDB() int
}

type redis struct {
	Addr     string
	Username string
	Password string
	DB       int
}

func (r *redis) GetAddr() string {
	return r.Addr
}

func (r *redis) GetUsername() string {
	return r.Username
}

func (r *redis) GetPassword() string {
	return r.Password
}

func (r *redis) GetDB() int {
	return r.DB
}
