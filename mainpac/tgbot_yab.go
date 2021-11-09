package mainpac

import (
	ue "Laurene/unmarshal_export"
	"Laurene/util"
	"fmt"
	tb "gopkg.in/tucnak/telebot.v3"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (s *Service) TgInfoUserYAB(x tb.Context) (errReturn error) {
	text := x.Text()

	i1 := strings.Index(text, "С нами")
	i2 := strings.Index(text, "Последний раз писал")
	tmTextStart := string([]rune(strings.Replace(text[i1:strings.Index(text[i1:], "\n")+i1], " ", "", -1))[5:])
	tmTextLastAction := string([]rune(strings.Replace(text[i2:strings.Index(text[i2:], "\n")+i2], " ", "", -1))[17:])

	tmDurStart, err := durationFromText(tmTextStart)
	if err != nil {
		x.Send("Ошибка преобразования времени")
		return
	}
	tmDurLastAction, err := durationFromText(tmTextLastAction)
	if err != nil {
		x.Send("Ошибка преобразования времени")
		return
	}

	var tmMesInt int64
	if x.Message().IsForwarded() {
		tmMesInt = int64(x.Message().OriginalUnixtime)
	} else {
		tmMesInt = time.Now().Unix()
	}
	tmMes := time.Unix(tmMesInt, 0)
	tmStart, tmLastAction := tmMes.Add(-tmDurStart), tmMes.Add(-tmDurLastAction)

	nick := strings.Replace(text, "[BOT] Информация о ", "", 1)
	nick = nick[:strings.Index(nick, ":\n")]
	loc, _ := time.LoadLocation("Europe/Moscow")
	textSend := fmt.Sprintf("%v\n\nС нами с %v\nПоследний раз писал в %v",
		nick, tmStart.In(loc).Format("2006-01-02 15:04:05"), tmLastAction.In(loc).Format("2006-01-02 15:04:05"))

	x.Send(textSend)
	return
}

func durationFromText(s string) (tm time.Duration, err error) {
	iDay := strings.Index(s, "d")
	iHour := strings.Index(s, "h")
	iMin := strings.Index(s, "m")
	iSec := strings.Index(s, "s")
	n, x := 0, 0

	if iDay >= 0 {
		n, err = strconv.Atoi(s[:iDay])
		if err != nil {
			return 0, err
		}
		tm = tm + time.Hour*24*time.Duration(n)
		x = iDay + 1
	}

	if iHour >= 0 {
		n, err = strconv.Atoi(s[x:iHour])
		if err != nil {
			return 0, err
		}
		tm = tm + time.Hour*time.Duration(n)
		x = iHour + 1
	}

	if iMin >= 0 {
		n, err = strconv.Atoi(s[x:iMin])
		if err != nil {
			return 0, err
		}
		tm = tm + time.Minute*time.Duration(n)
		x = iMin + 1
	}

	if iSec >= 0 {
		n, err = strconv.Atoi(s[x:iSec])
		if err != nil {
			return 0, err
		}
		tm = tm + time.Second*time.Duration(n)
	}

	return
}

func (s *Service) TgStatYABNotif(x tb.Context) (errReturn error) {
	if x.Message() == nil || x.Message().Document == nil {
		x.Send("Ошибка, напишите автору.")
		return
	}

	doc := x.Message().Document
	if doc.FileSize > 52420000 {
		x.Send("Слишком большой размер файла (больше 50мб).")
		return
	}

	path := "files/temp/" + util.CreateKey(10) + ".json"
	defer os.Remove(path)
	err := x.Bot().Download(doc.MediaFile(), path)
	if err != nil {
		x.Send("Ошибка, напишите автору.")
		return
	}

	file, err := os.Open(path)
	if err != nil {
		x.Send("Ошибка, напишите автору.")
		return
	}
	defer file.Close()

	keys := []string{"reply", "hug", "slap", "msg"}
	stat := &statNotif{
		minTime:  time.Unix(math.MaxInt64, 0),
		Top:      make(map[string]map[string]int),
		Nicks:    make(map[string][]string),
		TopFinal: make(map[string][]user),
	}
	for _, key := range keys {
		stat.Top[key] = make(map[string]int)
	}
	{
		errBool := true
		for msg := range ue.UnmarshalChan(file, 20) {
			if msg.FromID == "user163626570" {
				errBool = false
			}
		}
		if errBool {
			x.Send("Ошибка, неизвестный чат/файл.")
			return
		}
		file.Seek(0, 0)

		ch := ue.UnmarshalChan(file, 0)
		defer close(ch)
		for msg := range ch {
			if msg == nil || msg.FromID != "user163626570" {
				continue
			}
			if msg.Date.Unix() < 1549832400 {
				continue
			}
			s.statNotifProcessing(stat, msg)
		}

		for _, key := range keys {
			TopFinal(stat, stat.Top[key], stat.Nicks, key)
		}
	}

	text := fmt.Sprintf(""+
		"От %v"+
		"\nДо %v"+
		"\n\nОтветов: %v"+
		"\nОбнимашек: %v"+
		"\nШлепков: %v"+
		"\nЛичных сообщений: %v",
		stat.minTime.In(s.Loc).Format("2006-01-02 15:04 MST -0700"),
		stat.maxTime.In(s.Loc).Format("2006-01-02 15:04 MST -0700"),
		stat.ReplySum, stat.HugSum, stat.SlapSum, stat.MsgSum)
	x.Send(text)

	for _, key := range keys {
		text = "Топ " + key + "\n\n"
		for i, u := range stat.TopFinal[key] {
			text += fmt.Sprintf("%v. %v %v %v\n", i+1, u.number, u.id, u.nicks)
		}
		x.Send(text)
	}
	return
}

type statNotif struct {
	MsgSum, ReplySum, SlapSum, HugSum int
	minTime, maxTime                  time.Time

	Top   map[string]map[string]int
	Nicks map[string][]string

	TopFinal map[string][]user
}

type user struct {
	id     string
	nicks  string
	number int
}

func (s *Service) statNotifProcessing(stat *statNotif, msg *ue.Message) {
	id := getID(msg.Hashtags)
	nick := getNick(msg.Text)
	topID := ""
	ok := nick != "" && id != ""
	if ok {
		if _, okNick := stat.Nicks[id]; !okNick {
			stat.Nicks[id] = make([]string, 0)
		}
		exist := false
		for i := range stat.Nicks[id] {
			if stat.Nicks[id][i] == nick {
				exist = true
				break
			}
		}
		if !exist {
			stat.Nicks[id] = append(stat.Nicks[id], nick)
		}
	}

	switch {
	case s.Other.YABNotifReply.MatchString(msg.Text):
		stat.ReplySum++
		topID = "reply"
	case s.Other.YABNotifHug.MatchString(msg.Text):
		stat.HugSum++
		topID = "hug"
	case s.Other.YABNotifSlap.MatchString(msg.Text):
		stat.SlapSum++
		topID = "slap"
	case s.Other.YABNotifMsg.MatchString(msg.Text):
		stat.MsgSum++
		topID = "msg"
	default:
		return
	}

	if msg.Date.Unix() < stat.minTime.Unix() {
		stat.minTime = msg.Date
	}
	if msg.Date.Unix() > stat.maxTime.Unix() {
		stat.maxTime = msg.Date
	}

	if ok {
		stat.Top[topID][id]++
	}
}

func getNick(s string) (nick string) {
	i11 := strings.Index(s, "#")
	if i11 < 0 {
		return
	}
	i1 := strings.Index(s[i11:], " ")
	if i1 < 0 {
		return
	}
	i1 += i11 + 1

	i2 := strings.Index(s[i1:], " в чате")
	if i2 < 0 {
		return
	}
	i2 += i1

	nick = s[i1:i2]
	return
}

func getID(tags []string) (id string) {
	if len(tags) < 1 {
		return
	}
	return tags[0]
}

func TopFinal(stat *statNotif, top map[string]int, nicks map[string][]string, key string) {
	users := make([]user, 0, len(top))
	for s, i := range top {
		users = append(users, user{id: s, number: i})
	}
	for i := range users {
		users[i].nicks = strings.Join(nicks[users[i].id], "|")
	}
	sort.Slice(users, func(i, j int) bool { return users[i].number > users[j].number })
	if len(users) > 10 {
		stat.TopFinal[key] = users[:10]
	} else {
		stat.TopFinal[key] = users
	}
}

/*
От 11 февраля 2019
[BOT] Настройки уведомлений:
мсг   [BOT] Тебе отправлено личное сообщение от #S �Sin в чате @YetAnotherBot.
ответ [BOT] Ответ от #Alizee ����‍� в чате @YetAnotherBot - #ttvec.
шлеп  [BOT] Шлепок от #PC Старейшина� в чате @YetAnotherBot - #vszgp.
хуг   [BOT] Обнимашка от #clip �Скрепка в чате @YetAnotherBot - #qszst.
*/
