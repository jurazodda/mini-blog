package config

import "os"

var (
	env *Env
)

type Env struct {
	Db struct {
		Dsn string
	}
	Jwt struct {
		SigningKey string
	}
}

func GetEnv() *Env {
	env = &Env{
		Db: struct {
			Dsn string
		}{
			Dsn: os.Getenv("DSN_URL"),
		},
		Jwt: struct {
			SigningKey string
		}{
			SigningKey: os.Getenv("SIGNING_KEY"),
		},
	}
	return env
}
