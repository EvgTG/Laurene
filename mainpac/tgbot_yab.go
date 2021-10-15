package mainpac

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v3"
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
