package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"Laurene/go-log"
	"Laurene/mainpack"
	"Laurene/model"
	"Laurene/mongodb"
	"Laurene/util"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rotisserie/eris"
	"go.uber.org/fx"
	tb "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/layout"
)

func NewApp() (app *fx.App) {
	app = fx.New(
		fx.Provide(
			NewDB,
			NewService,
		),

		fx.Invoke(
			ReadConfig,
			Logger,
			PingServe,
			Start,
		),
	)
	return
}

func ReadConfig() {
	err := cleanenv.ReadConfig("files/cfg.env", &CFG)
	if err != nil {
		panic(eris.Wrap(err, "ReadConfig"))
	}
}

func Logger() {
	log.SetLogger(log.New(CFG.LogLevel, true))
	log.Info("Go!")
}

func NewDB() *model.Model {
	return model.New(mongodb.NewDB(CFG.NameDB, CFG.MongoUrl))
}

func NewService(lc fx.Lifecycle /*db *model.Model*/) *mainpack.Service {
	lt, err := layout.New("bot.yml")
	util.ErrCheckFatal(err, "layout.New()", "NewService", "init")

	bot, err := tb.NewBot(tb.Settings{
		Token: CFG.TgApiToken,
		Poller: &tb.Webhook{
			Listen:      ":" + CFG.LocalPort,
			IP:          CFG.IP,
			SecretToken: CFG.SecretToken,

			Endpoint: &tb.WebhookEndpoint{
				PublicURL: "https://" + CFG.IP + ":" + CFG.Port + CFG.Path,
				Cert:      "cert/pem.pem",
			},

			// +.tls если без обратного прокси
		},
	})
	// long-polling &tb.LongPoller{Timeout: 30 * time.Second},
	util.ErrCheckFatal(err, "tb.NewBot()", "NewService", "init")
	bot.Use(lt.Middleware("ru"))

	service := &mainpack.Service{
		Bot: &mainpack.Bot{
			Bot:           bot,
			Layout:        lt,
			Username:      bot.Me.Username,
			UserList:      CFG.UserList,
			AdminList:     CFG.AdminList,
			NotifyList:    CFG.NotifyList,
			ErrorList:     CFG.ErrorList,
			Uptime:        time.Now(),
			CallbackQuery: make(map[int64]string),
			AlbumsManager: util.NewAlbumsManager(),
			VideoAlbumsManager: &mainpack.VideoAlbumsManager{
				Map:     make(map[int64]*mainpack.VideoAlbum),
				MapLock: make(map[int64]bool),
				Mutex:   sync.Mutex{},
			},
			Rand: rand.New(rand.NewSource(time.Now().UnixNano())),
		},
		DB:   nil,
		Loc:  CFG.TimeLocation.Get(),
		Rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	err = service.Bot.RemoveWebhook()
	util.ErrCheckFatal(err, "Bot.RemoveWebhook()", "NewService", "init")

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			// service.Bot.Stop()
			service.Bot.RemoveWebhook()
			return nil
		},
	})

	return service
}

func Start(s *mainpack.Service) {
	s.Start()
}

func PingServe() {
	if !CFG.PingOn || CFG.PingPort == "" {
		log.Info("PingServer off")
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/pingLaurene", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "pong")
	})
	log.Info("PingServe on")
	go http.ListenAndServe(":"+CFG.PingPort, mux)
}
