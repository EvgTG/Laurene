package main

import (
	"github.com/rotisserie/eris"
	"time"
)

var CFG InitConfig

type InitConfig struct {
	LogLevel string `env:"LOGLVL" env-default:"INFO"`

	ProxyTG      string         `env:"PROXYTG"`
	TgApiToken   string         `env:"TOKENTG"`
	UserList     []int64        `env:"USERLIST"`
	AdminList    []int64        `env:"ADMINLIST"`
	NotifyList   []int64        `env:"NOTIFYLIST"`
	ErrorList    []int64        `env:"ERRORLIST"`
	TimeLocation MyTimeLocation `env:"LOC" env-default:"UTC"`

	NameDB   string `env:"NAMEDB" env-default:"tgbotTemplate"`
	MongoUrl string `env:"MONGOURL" env-default:"mongodb://127.0.0.1:27017"`

	PingPort string `env:"PINGPORT" env-default:"6975"`
	PingOn   bool   `env:"PINGON" env-default:"false"`
}

type MyTimeLocation string

func (l *MyTimeLocation) SetValue(s string) error {
	*l = MyTimeLocation(s)
	return nil
}

func (l MyTimeLocation) Get() *time.Location {
	loc, err := time.LoadLocation(string(l))
	if err != nil {
		panic(eris.Wrap(err, "cfg.GetTimeLocation()"))
	}
	return loc
}