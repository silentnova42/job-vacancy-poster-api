package pgstorage

import "fmt"

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	Dbname   string
}

func (c *Config) GetUrlConn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.Username, c.Password, c.Host, c.Port, c.Dbname)
}
