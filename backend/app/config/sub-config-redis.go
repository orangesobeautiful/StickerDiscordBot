package config

type Database interface {
	GetDSN() string
	GetAutoMigrate() bool
}

type database struct {
	DSN         string
	AutoMigrate bool
}

func (d *database) GetDSN() string {
	return d.DSN
}

func (d *database) GetAutoMigrate() bool {
	return d.AutoMigrate
}
