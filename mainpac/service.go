package mainpac

import (
	"Laurene/model"
	"Laurene/util"
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
