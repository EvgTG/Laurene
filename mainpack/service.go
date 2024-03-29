package mainpack

import (
	"Laurene/go-log"
	"Laurene/model"
	"Laurene/util"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	tb "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/layout"
	"math/rand"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Service struct {
	Bot   *Bot
	Other Other

	DB *model.Model

	Loc  *time.Location
	Rand *rand.Rand
}

type Bot struct {
	*tb.Bot
	*layout.Layout

	UserList   []int64
	AdminList  []int64
	NotifyList []int64
	ErrorList  []int64

	Username string
	Uptime   time.Time
	Rand     *rand.Rand

	CallbackQuery      map[int64]string // Контекстный ввод
	CallbackQueryMutex sync.Mutex

	AlbumsManager      *util.AlbumsManager
	VideoAlbumsManager *VideoAlbumsManager
}

type Other struct {
	YABInfoUserRGX *regexp.Regexp

	YABNotifMsg   *regexp.Regexp
	YABNotifReply *regexp.Regexp
	YABNotifSlap  *regexp.Regexp
	YABNotifHug   *regexp.Regexp

	AtbashAlphabet *strings.Replacer
	AtbashCache    *lru.Cache
}

func (s Service) Start() {
	log.Info("tgbot init")
	s.InitBot()

	log.Info("tgbot launch...")
	fmt.Println("tgbot @" + s.Bot.Bot.Me.Username)

	go s.GoCheckErrs()
	go s.Bot.Start()
}

func (s Service) GoCheckErrs() {
	time.Sleep(time.Second * 30)
	nErr := log.GetErrN()
	if nErr > 0 {
		s.Bot.sendToSlice(s.Bot.ErrorList, fmt.Sprintf("Новых ошибок: %v.\n Заляните в логи.", nErr))
	}

	for range time.Tick(time.Minute * 5) {
		nErr = log.GetErrN()
		if nErr > 0 {
			s.Bot.sendToSlice(s.Bot.ErrorList, fmt.Sprintf("Новых ошибок: %v.\n Заляните в логи.", nErr))
		}
	}
}
