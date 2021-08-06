package mainpac

import (
	"Laurene/go-log"
	"Laurene/model"
	"Laurene/util"
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"math/rand"
	"time"
)

type Service struct {
	TG   *TG
	DB   *model.Model
	Loc  *time.Location
	Rand *rand.Rand
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

func (s Service) Start() {
	log.Info("tgbot init")
	s.InitTBot()
	log.Info("tgbot launch...")
	fmt.Println("tgbot @" + s.TG.Bot.Me.Username)
	s.TG.Bot.Start()
}
