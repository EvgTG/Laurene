package mainpack

import (
	"Laurene/go-log"
	"Laurene/util"
	"fmt"
	"github.com/rotisserie/eris"
	tb "gopkg.in/telebot.v3"
	"os"
	"strconv"
	"strings"
)

/*
CMD -        команда
Update/Btn - кнопка обновления/обычная
Func -       логика работы
Но они обязательны только все вместе
*/

func (s *Service) TgStart(x tb.Context) (errReturn error) {
	x.Send(s.Bot.Text(x, "start"), &tb.ReplyMarkup{RemoveKeyboard: true}, tb.ModeHTML)
	return
}

func (s *Service) TgStartYAB(x tb.Context) (errReturn error) {
	x.Send(s.Bot.Text(x, "start_yab"), &tb.ReplyMarkup{RemoveKeyboard: true}, tb.ModeHTML)
	return
}

func (s *Service) TgOnText(x tb.Context) (errReturn error) {
	switch {
	case s.Other.YABInfoUserRGX.MatchString(x.Text()):
		s.TgInfoUserYAB(x)
		return
	default:
		x.Send(s.Bot.Text(x, "ontext"), &tb.SendOptions{ReplyTo: x.Message()}, s.Bot.Markup(x, "text"))
	}

	s.Bot.CallbackQueryMutex.Lock()
	cq := s.Bot.CallbackQuery[x.Chat().ID]
	s.Bot.CallbackQueryMutex.Unlock()

	// Не запускать функции отдельной горутиной - теряется контекст
	switch cq {
	case "": //Нет в CallbackQuery - игнор
	case "test":
	}

	return
}

func (s *Service) TgYABNotification(x tb.Context) (errReturn error) {
	x.Send(s.Bot.Text(x, "yab"))
	return
}

func (s *Service) TgOnTextInline(x tb.Context) (errReturn error) {
	q := x.Query()
	if q == nil {
		return
	}
	if q.Text == "" {
		return
	}

	res := make([]tb.Result, 0, 2)
	text := ""
	q.Text = strings.TrimSpace(q.Text)

	// Текст в обратном порядке
	text = textReverse(q.Text)
	ar := &tb.ArticleResult{Title: s.Bot.Text(x, "ti_reverse"), Text: "<pre>" + text + "</pre>", Description: util.TextCut(text, 50)}
	ar.ParseMode = tb.ModeHTML
	res = append(res, ar)

	// Шифр Атбаша
	text = s.Other.AtbashAlphabet.Replace(q.Text)
	key := util.CreateKey(8)
	s.Other.AtbashCache.Add(key, q.Text)
	ar = &tb.ArticleResult{Title: s.Bot.Text(x, "ti_atbash"), Text: "<pre>" + text + "</pre>", Description: util.TextCut(text, 50)}
	ar.ParseMode = tb.ModeHTML
	rm := *s.Bot.Markup(x, "atbash")
	rm.InlineKeyboard[0][0].Data = key
	ar.ReplyMarkup = &rm
	res = append(res, ar)

	// Текст в переводе на эмоджи
	text = textEmoji(q.Text)
	ar = &tb.ArticleResult{Title: s.Bot.Text(x, "ti_emoji"), Text: "<pre>" + text + "</pre>", Description: util.TextCut(text, 50)}
	ar.ParseMode = tb.ModeHTML
	res = append(res, ar)

	// Текст в верхнем регистре
	text = strings.ToUpper(q.Text)
	ar = &tb.ArticleResult{Title: s.Bot.Text(x, "ti_upper"), Text: "<pre>" + text + "</pre>", Description: util.TextCut(text, 50)}
	ar.ParseMode = tb.ModeHTML
	res = append(res, ar)

	// Текст в случайном регистре
	text = textRandom(q.Text, s.Rand)
	ar = &tb.ArticleResult{Title: s.Bot.Text(x, "ti_random"), Text: "<pre>" + text + "</pre>", Description: util.TextCut(text, 50)}
	ar.ParseMode = tb.ModeHTML
	res = append(res, ar)

	for i := range res {
		res[i].SetResultID(strconv.Itoa(i))
	}
	qr := &tb.QueryResponse{
		QueryID:    q.ID,
		CacheTime:  0,
		IsPersonal: true,
		Results:    res,
	}
	x.Answer(qr)

	return
}

// Ниже только админское

func (s *Service) TgTest(x tb.Context) (errReturn error) {
	x.Send(s.Bot.Text(x, "test"), s.Bot.Markup(x, "test"), tb.ModeHTML, tb.NoPreview)
	return
}

func (s *Service) TgTestBtn(x tb.Context) (errReturn error) {
	x.Send(s.Bot.Text(x, "test"), &tb.SendOptions{ReplyTo: x.Message()}, s.Bot.Markup(x, "test"), tb.ModeHTML, tb.NoPreview)
	x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "test"})
	return
}

func (s *Service) TgAdm(x tb.Context) (errReturn error) {
	text := fmt.Sprintf("" +
		"<b>Пользователькие команды:</b>" +
		"\n/start - приветствие" +
		"\n\n<b>Админские команды:</b>" +
		"\n/status - статус работы" +
		"\n/logs - действия над логами" +
		"\n/setCmds - установить команды бота",
	)

	x.Send(text, tb.ModeHTML)
	return
}

func (s *Service) TgStatus(x tb.Context) (errReturn error) {
	text, rm := s.TgStatusFunc(x)
	mes, err := s.Bot.Send(x.Sender(), text, rm, tb.ModeHTML)
	if err == nil && mes != nil {
		s.Bot.Pin(mes)
	}
	return
}

func (s *Service) TgStatusUpdate(x tb.Context) (errReturn error) {
	text, rm := s.TgStatusFunc(x)
	_, err := s.Bot.Edit(x.Message(), text, rm, tb.ModeHTML)
	if err != nil {
		if err == tb.ErrSameMessageContent {
			x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Обновлено"})
			return
		}

		s.Bot.Delete(x.Message())
		mes, err := s.Bot.Send(x.Sender(), text, rm, tb.ModeHTML)
		if err == nil && mes != nil {
			s.Bot.Pin(mes)
		}
	}

	x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Обновлено"})
	return
}

func (s *Service) TgStatusFunc(x tb.Context) (string, *tb.ReplyMarkup) {
	text := fmt.Sprintf(""+
		"Запущен: %s"+
		"\nUptime: %s"+
		"\n"+
		"\nAlbums manager: %v",
		s.Bot.Uptime.In(s.Loc).Format("2006.01.02 15:04:05 MST"), s.Bot.uptimeString(s.Bot.Uptime),
		s.Bot.AlbumsManager.Len(),
	)

	rm := s.Bot.Markup(x, "status")

	return text, rm
}

func (s *Service) TgLogs(x tb.Context) (errReturn error) {
	text := "1. Получить файл логов\n2. Очистить файл логов"
	x.Send(text, s.Bot.Markup(x, "logs"), tb.ModeHTML)
	return
}

func (s *Service) TgGetLogsBtn(x tb.Context) (errReturn error) {
	err := x.Send(&tb.Document{File: tb.FromDisk("files/logrus.log"), FileName: "logrus.log"})
	if err != nil {
		s.Bot.Send(x.Sender(), eris.Wrap(err, "Ошибка отправки файла.").Error())
	}
	x.Respond()
	return
}

func (s *Service) TgClearLogsBtn(x tb.Context) (errReturn error) {
	os.Truncate("files/logrus.log", 0)
	log.Info("Очищено")

	x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Очищено", ShowAlert: true})
	return
}

func (s *Service) TgDeleteBtn(x tb.Context) (errReturn error) {
	x.Respond()
	x.Delete()
	return
}

func (s *Service) TgCancelReplyMarkup(x tb.Context) (errReturn error) {
	s.DeleteCallbackQuery(x.Chat().ID)

	x.Send("Отменено.", &tb.ReplyMarkup{RemoveKeyboard: true})
	return
}

func (s *Service) DeleteCallbackQuery(chatID int64) {
	s.Bot.CallbackQueryMutex.Lock()
	delete(s.Bot.CallbackQuery, chatID)
	s.Bot.CallbackQueryMutex.Unlock()
}

func (s *Service) TgSetCommands(x tb.Context) (errReturn error) {
	err := s.Bot.SetCommands(s.Bot.Layout.Commands())
	if err != nil {
		x.Send(eris.Wrap(err, "s.Bot.SetCommands()").Error())
		return
	}

	x.Send("Готово.")
	return
}

/*

func (s *Service) TgSome(x tb.Context) (errReturn error) {
	return
}

func (s *Service) TgBtn(x tb.Context) (errReturn error) {
	//text, rm := s.TgFunc()
	//x.Send(text, rm, tb.ModeHTML)
	//s.Bot.Edit(x.Message(), text, rm)
	x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: ""})
}
func (s *Service) TgCMD(x tb.Context) (errReturn error) {
	text, rm := s.TgInfoFilmFunc(m, 0)
	rm := s.Bot.Markup(x, "test")
	x.Send(text, rm, tb.ModeHTML)
}

func (s *Service) TgFunc() (string, *tb.ReplyMarkup) {
	text := ""
	rm := s.Bot.Markup(x, "test")
}

*/
