package main

import (
	"Laurene/go-config"
	"Laurene/go-log"
	"Laurene/mainpac"
	"Laurene/model"
	"Laurene/mongodb"
	"Laurene/util"
	"fmt"
	"go.uber.org/fx"
	tb "gopkg.in/tucnak/telebot.v2"
	"math/rand"
	"net/http"
	"time"
)

func New() (app *fx.App) {
	app = fx.New(
		fx.Provide(
			ConfigReader,
			Config,
			ConfigLogger,
			NewDB,
			NewService,
		),

		fx.Invoke(
			Logger,
			PingServe,
			Start,
		),
	)
	return
}

type initConfig struct {
	ProxyTG    string
	TgApiToken string
	UserList   []int64
	AdminList  []int64
	NotifyList []int64
	ErrorList  []int64
	Loc        *time.Location

	NameDB   string
	MongoUrl string
	PingPort string
}

type initConfigLogger struct {
	LogLvl string
}

func ConfigReader() *config.Config {
	return config.New()
}

func Config(cfgRead *config.Config) *initConfig {
	loc, err := cfgRead.GetTimeLocation("LOC")
	util.ErrCheckFatal(err, "cfgRead.GetTimeLocation", "Config", "init")

	cfg := &initConfig{
		ProxyTG:    cfgRead.GetString("PROXYTG"),
		TgApiToken: cfgRead.GetString("TOKENTG"),
		UserList:   cfgRead.GetIntSlice64("USERLIST"),
		AdminList:  cfgRead.GetIntSlice64("ADMINLIST"),
		NotifyList: cfgRead.GetIntSlice64("NOTIFYLIST"),
		ErrorList:  cfgRead.GetIntSlice64("ERRORLIST"),
		Loc:        loc,
		NameDB:     cfgRead.GetString("NAMEDB"),
		MongoUrl:   cfgRead.GetString("MONGOURL"),
		PingPort:   cfgRead.GetString("PINGPORT"),
	}

	return cfg
}

func ConfigLogger(cfgRead *config.Config) *initConfigLogger {
	cfg := &initConfigLogger{
		LogLvl: cfgRead.GetString("LOGLVL"),
	}
	return cfg
}

func Logger(cfg *initConfigLogger) {
	log.SetLogger(log.New(cfg.LogLvl, true))
	log.Info("Go!")
}

func NewDB(cfg *initConfig) *model.Model {
	return model.New(mongodb.NewDB(cfg.NameDB, cfg.MongoUrl))
}

func NewService(cfg *initConfig /*, db *model.Model*/) *mainpac.Service {
	bot, err := tb.NewBot(tb.Settings{
		Token:  cfg.TgApiToken,
		Poller: &tb.LongPoller{Timeout: 30 * time.Second},
	})
	util.ErrCheckFatal(err, "tb.NewBot", "NewService", "init")

	service := &mainpac.Service{
		TG: &mainpac.TG{
			Bot:           bot,
			Username:      bot.Me.Username,
			UserList:      cfg.UserList,
			AdminList:     cfg.AdminList,
			NotifyList:    cfg.NotifyList,
			ErrorList:     cfg.ErrorList,
			Uptime:        time.Now(),
			CallbackQuery: make(map[int64]string),
			AlbumsManager: util.NewAlbumsManager(),
		},
		DB:   nil,
		Loc:  cfg.Loc,
		Rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	return service
}

func Start(s *mainpac.Service) {
	s.Start()
}

func PingServe(cfg *initConfig) {
	if cfg.PingPort == "" {
		log.Info("PingServe off")
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/pingLaurene", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "pong")
	})
	log.Info("PingServe on")
	go http.ListenAndServe(":"+cfg.PingPort, mux)
}
