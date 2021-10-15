package mainpac

import (
	"Laurene/go-log"
	"Laurene/model"
	"Laurene/util"
	"fmt"
	tb "gopkg.in/tucnak/telebot.v3"
	"math/rand"
	"regexp"
	"time"
)

type Service struct {
	TG    *TG
	Other Other
	DB    *model.Model
	Loc   *time.Location
	Rand  *rand.Rand
}

type TG struct {
	Bot           *tb.Bot
	Username      string
	UserList      []int64
	AdminList     []int64
	NotifyList    []int64
	ErrorList     []int64
	Uptime        time.Time
	Buttons       map[string]*tb.Btn
	CallbackQuery map[int64]string //контекстный ввод

	AlbumsManager *util.AlbumsManager
}

type Other struct {
	YetAnotherBotInfoUserRGX *regexp.Regexp
}

func (s Service) Start() {
	log.Info("tgbot init")
	s.InitTBot()
	log.Info("tgbot launch...")
	fmt.Println("tgbot @" + s.TG.Bot.Me.Username)
	go s.GoCheckErrs()
	s.TG.Bot.Start()
}

func (s Service) GoCheckErrs() {
	time.Sleep(time.Second * 30)
	nErr := log.GetErrN()
	if nErr > 0 {
		s.TG.sendToSlice(s.TG.ErrorList, fmt.Sprintf("Новых ошибок: %v.\n Заляните в логи.", nErr))
	}

	for range time.Tick(time.Minute * 5) {
		nErr = log.GetErrN()
		if nErr > 0 {
			s.TG.sendToSlice(s.TG.ErrorList, fmt.Sprintf("Новых ошибок: %v.\n Заляните в логи.", nErr))
		}
	}
}
