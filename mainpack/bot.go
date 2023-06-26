package mainpack

import (
	"Laurene/util"
	lru "github.com/hashicorp/golang-lru"
	tb "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"io/fs"
	"os"
	"regexp"
	"strings"
)

func (s *Service) InitBot() {
	s.InitOther()

	// Команды

	s.Bot.Bot.Handle("/start", s.TgStart)
	s.Bot.Bot.Handle("/help", s.TgStart)
	s.Bot.Bot.Handle("/yab", s.TgStartYAB)
	s.Bot.Bot.Handle("/YABNotification", s.TgYABNotification)
	s.Bot.Bot.Handle("/emoji", s.TgEmojiAlphabet)
	s.Bot.Bot.Handle(tb.OnPhoto, s.TgPic)
	s.Bot.Bot.Handle(tb.OnText, s.TgOnText)
	s.Bot.Bot.Handle(tb.OnQuery, s.TgOnTextInline)
	s.Bot.Bot.Handle(tb.OnDocument, s.TgStatYABNotif)
	s.Bot.Bot.Handle(tb.OnVideo, s.TgVideoComb)

	// Админские команды

	adminOnly := s.Bot.Group()
	adminOnly.Use(middleware.Whitelist(s.Bot.AdminList...))

	adminOnly.Handle("/test", s.TgTest)
	adminOnly.Handle("/adm", s.TgAdm)
	adminOnly.Handle("/status", s.TgStatus)
	adminOnly.Handle("/logs", s.TgLogs)
	adminOnly.Handle("/setCmds", s.TgSetCommands)

	// Кнопки

	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "delete"), s.TgDeleteBtn)
	adminOnly.Handle(s.Bot.Layout.ButtonLocale("", "test"), s.TgTestBtn)
	adminOnly.Handle(s.Bot.Layout.ButtonLocale("", "cancel"), s.TgCancelReplyMarkup)
	adminOnly.Handle(s.Bot.Layout.ButtonLocale("", "status_update"), s.TgStatusUpdate)
	adminOnly.Handle(s.Bot.Layout.ButtonLocale("", "get_logs"), s.TgGetLogsBtn)
	adminOnly.Handle(s.Bot.Layout.ButtonLocale("", "clear_logs"), s.TgClearLogsBtn)

	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "album_to_pic_down"), s.TgAlbumToPic)
	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "album_to_pic_right"), s.TgAlbumToPic)
	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "album_to_pic_mesh"), s.TgAlbumToPic)
	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "album_compress1"), s.TgCompress)
	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "album_compress2"), s.TgCompress)
	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "album_compress3"), s.TgCompress)
	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "pic_compress1"), s.TgCompress)
	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "pic_compress2"), s.TgCompress)
	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "pic_compress3"), s.TgCompress)
	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "pic_gif"), s.TgPicGif)
	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "picfile_to_pic"), s.TgFilePicToPic)
	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "text_reverse"), s.TgTextReverse)
	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "text_toupper"), s.TgTextToUpper)
	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "text_random"), s.TgTextRandom)
	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "text_atbash"), s.TgTextAtbash)
	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "text_emoji"), s.TgTextEmoji)
	s.Bot.Handle(s.Bot.Layout.ButtonLocale("", "atbash_btn"), s.TgTextAtbashBtn)
}

/*

s.TG.Bot.Handle("/", s.Tg)
s.Bot.Handle(s.Bot.Layout.ButtonLocale("", ""), s.Tg)

*/

func (s *Service) InitOther() {
	var err error
	dir := "files/temp/"
	os.RemoveAll(dir)
	os.Mkdir(dir, fs.ModePerm)

	// YetAnotherBot RGX
	s.Other.YABInfoUserRGX, err = regexp.Compile("^\\[BOT\\] Информация о .{1,2} #.+:\\n")
	util.ErrCheckFatal(err, "InitRgxs", "YABInfoUserRGX")

	// YetAnotherBot Notify RGX
	s.Other.YABNotifMsg, err = regexp.Compile("^\\[BOT\\] Тебе отправлено личное сообщение от")
	util.ErrCheckFatal(err, "InitRgxs", "YABNotifMsg")
	s.Other.YABNotifReply, err = regexp.Compile("^\\[BOT\\] Ответ от")
	util.ErrCheckFatal(err, "InitRgxs", "YABNotifReply")
	s.Other.YABNotifSlap, err = regexp.Compile("^\\[BOT\\] Шлепок от")
	util.ErrCheckFatal(err, "InitRgxs", "YABNotifSlap")
	s.Other.YABNotifHug, err = regexp.Compile("^\\[BOT\\] Обнимашка от")
	util.ErrCheckFatal(err, "InitRgxs", "YABNotifHug")

	// Atbash Cache
	s.Other.AtbashCache, _ = lru.New(1000)
	// Atbash
	eng := "abcdefghijklmnopqrstuvwxyz"
	engr := "zyxwvutsrqponmlkjihgfedcba"
	eng2 := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	eng2r := "ZYXWVUTSRQPONMLKJIHGFEDCBA"
	ru := "абвгдеёжзийклмнопрстуфхцчшщъыьэюя"
	rur := "яюэьыъщшчцхфутсрпонмлкйизжёедгвба"
	ru2 := "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"
	ru2r := "ЯЮЭЬЫЪЩШЧЦХФУТСРПОНМЛКЙИЗЖЁЕДГВБА"

	oldnew := make([]string, 0, 26*4+33*4)
	alphabets := [][][]rune{{[]rune(eng), []rune(engr)}, {[]rune(eng2), []rune(eng2r)}, {[]rune(ru), []rune(rur)}, {[]rune(ru2), []rune(ru2r)}}
	for ii := range alphabets {
		for i := 0; i < len(alphabets[ii][0]); i++ {
			oldnew = append(oldnew, string(alphabets[ii][0][i]))
			oldnew = append(oldnew, string(alphabets[ii][1][i]))
		}
	}
	s.Other.AtbashAlphabet = strings.NewReplacer(oldnew...)
}
