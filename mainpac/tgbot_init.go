package mainpac

import (
	"Laurene/util"
	lru "github.com/hashicorp/golang-lru"
	tb "gopkg.in/tucnak/telebot.v3"
	"os"
	"regexp"
	"strings"
)

func (s *Service) InitTBot() {
	s.InitOther()

	// Команды роуты

	s.TG.Bot.Handle("/start", s.TgStartCMD)
	s.TG.Bot.Handle("/help", s.TgStartCMD)
	s.TG.Bot.Handle("/yab", s.TgStartYAB)
	s.TG.Bot.Handle("/YABNotification", s.TgYABNotification)
	s.TG.Bot.Handle(tb.OnPhoto, s.TgPic)
	s.TG.Bot.Handle(tb.OnText, s.TgOnText)
	s.TG.Bot.Handle(tb.OnQuery, s.TgOnTextInline)
	s.TG.Bot.Handle(tb.OnDocument, s.TgStatYABNotif)

	// Админские команды

	s.TG.Bot.Handle("/test", s.TgTest)
	s.TG.Bot.Handle("/adm", s.TgAdm)
	s.TG.Bot.Handle("/status", s.TgStatusCMD)
	s.TG.Bot.Handle("/logs", s.TgLogsCMD)
	s.TG.Bot.Handle("/setCmds", s.TgSetCmds)

	// Кнопки роуты

	rm := &tb.ReplyMarkup{}
	im := &tb.ReplyMarkup{ResizeKeyboard: true}
	s.TG.Buttons = make(map[string]*tb.Btn)
	iq := rm.Data("Написать", "")
	iq.InlineQueryChat = " "

	s.TG.addBtn(rm.Data("1️⃣ вниз", "album_to_pic_down", "down"), "album_to_pic_down", s.TgAlbumToPic)
	s.TG.addBtn(rm.Data("1️⃣ вправо", "album_to_pic_right", "right"), "album_to_pic_right", s.TgAlbumToPic)
	s.TG.addBtn(rm.Data("1️⃣ сеткой", "album_to_pic_mesh", "mesh"), "album_to_pic_mesh", s.TgAlbumToPic)
	s.TG.addBtn(rm.Data("2️⃣\U0001F7E9", "album_compress1", "cp1"), "album_compress1", s.TgCompress)
	s.TG.addBtn(rm.Data("2️⃣\U0001F7E8", "album_compress2", "cp2"), "album_compress2", s.TgCompress)
	s.TG.addBtn(rm.Data("2️⃣\U0001F7E5", "album_compress3", "cp3"), "album_compress3", s.TgCompress)
	s.TG.addBtn(rm.Data("1️⃣\U0001F7E9", "pic_compress1", "cp1"), "pic_compress1", s.TgCompress)
	s.TG.addBtn(rm.Data("1️⃣\U0001F7E8", "pic_compress2", "cp2"), "pic_compress2", s.TgCompress)
	s.TG.addBtn(rm.Data("1️⃣\U0001F7E5", "pic_compress3", "cp3"), "pic_compress3", s.TgCompress)
	s.TG.addBtn(rm.Data("🖼 Отправить картинкой", "picfile_to_pic"), "picfile_to_pic", s.TgFilePicToPic)
	s.TG.addBtn(rm.Data("1️⃣", "text_reverse", "1"), "text_reverse", s.TgTextReverse)
	s.TG.addBtn(rm.Data("2️⃣", "text_toupper", "2"), "text_toupper", s.TgTextToUpper)
	s.TG.addBtn(rm.Data("3️⃣", "text_random", "3"), "text_random", s.TgTextRandom)
	s.TG.addBtn(rm.Data("4️⃣", "text_atbash", "4"), "text_atbash", s.TgTextAtbash)
	s.TG.addBtn(rm.Data("Расшифровать", "atbash_btn"), "atbash_btn", s.TgTextAtbashBtn)
	s.TG.addBtn(iq, "iq", s.TgTest)

	s.TG.menu.picAlbumsBtns = &tb.ReplyMarkup{}
	s.TG.menu.picAlbumsBtns.Inline(
		[]tb.Btn{*s.TG.Buttons["album_to_pic_down"], *s.TG.Buttons["album_to_pic_right"], *s.TG.Buttons["album_to_pic_mesh"]},
		[]tb.Btn{*s.TG.Buttons["album_compress1"], *s.TG.Buttons["album_compress2"], *s.TG.Buttons["album_compress3"]},
	)
	s.TG.menu.picBtns = &tb.ReplyMarkup{}
	s.TG.menu.picBtns.Inline([]tb.Btn{*s.TG.Buttons["pic_compress1"], *s.TG.Buttons["pic_compress2"], *s.TG.Buttons["pic_compress3"]})

	s.TG.menu.textBtns = &tb.ReplyMarkup{}
	s.TG.menu.textBtns.Inline([]tb.Btn{*s.TG.Buttons["text_reverse"], *s.TG.Buttons["text_toupper"], *s.TG.Buttons["text_random"], *s.TG.Buttons["text_atbash"]})

	s.TG.menu.atbashBtns = &tb.ReplyMarkup{}
	s.TG.menu.atbashBtns.Inline([]tb.Btn{*s.TG.Buttons["atbash_btn"]}, []tb.Btn{*s.TG.Buttons["iq"]})
	s.TG.menu.atbashBtns2 = &tb.InlineKeyboardMarkup{InlineKeyboard: [][]tb.InlineButton{{*s.TG.Buttons["atbash_btn"].Inline()}, {*s.TG.Buttons["iq"].Inline()}}}

	// Админские кнопки

	s.TG.addBtn(rm.Data("Test", "test"), "test", s.TgTestBtn)
	s.TG.addBtn(rm.Data("🗑Удалить", "delete"), "delete", s.TgDeleteBtn)
	s.TG.addBtn(im.Text("❌Отмена"), "cancel", s.TgCancelReplyMarkup)
	s.TG.addBtn(rm.Data("🔄Обновить", "status_update"), "status_update", s.TgStatusUpdate)
	s.TG.addBtn(rm.Data("1", "get_logs"), "get_logs", s.TgGetLogsBtn)
	s.TG.addBtn(rm.Data("2", "clear_logs"), "clear_logs", s.TgClearLogsBtn)
}

/*

s.TG.Bot.Handle("/", s.Tg)
s.TG.addBtn(rm.Data("", ""), "", s.Tg)

func (s *Service) TgSome(x tb.Context) (errReturn error) {
	return
}

*/

func (s *Service) InitOther() {
	var err error
	dir := "files/temp/"
	os.RemoveAll(dir)
	os.Mkdir(dir, 777)

	// YetAnotherBot RGX
	s.Other.YABInfoUserRGX, err = regexp.Compile("^\\[BOT\\] Информация о .{1,2} #.+:\\n")
	util.ErrCheckFatal(err, "InitRgxs", "YABInfoUserRGX")

	// YetAnotherBot RGX Notif
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
