package config

type Database interface {
	GetDSN() string
	GetAutoMigrate() bool
	GetDisableVersionedMigrate() bool
}

type database struct {
	DSN                     string
	AutoMigrate             bool
	DisableVersionedMigrate bool
}

func (d *database) GetDSN() string {
	return d.DSN
}

func (d *database) GetAutoMigrate() bool {
	return d.AutoMigrate
}

func (d *database) GetDisableVersionedMigrate() bool {
	return d.DisableVersionedMigrate
}
