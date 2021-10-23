package mainpac

import (
	"Laurene/go-log"
	"Laurene/util"
	"fmt"
	"github.com/pkg/errors"
	tb "gopkg.in/tucnak/telebot.v3"
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

func (s *Service) TgStartCMD(x tb.Context) (errReturn error) {
	text := "" +
		"Приветствую, бот имеет следующие возможности:" +
		"\n" +
		"\n<b>Текст:</b>  (доступно в инлайн режиме)" +
		"\n• Написание текста в обратном порядке." +
		"\n• Текст в верхнем регистре." +
		"\n• Текст в случайном регистре." +
		"\n• Шифр Атбаш" +
		"\n" +
		"\n<b>Разное:</b>" +
		"\n• Склейка фото, если прислать или переслать альбом." +
		"\n• Счёт дат в сообщении информации о пользователе из @YetAnotherBot (переслать)."

	x.Send(text, &tb.ReplyMarkup{RemoveKeyboard: true}, tb.ModeHTML)
	return
}

func (s *Service) TgOnText(x tb.Context) (errReturn error) {
	switch {
	case s.Other.YetAnotherBotInfoUserRGX.MatchString(x.Text()):
		s.TgInfoUserYAB(x)
		return
	default:
		text := "" +
			"Что сделать с текстом?" +
			"\n" +
			"\n1. Обратный порядок" +
			"\n2. В верхнем регистре" +
			"\n3. В случайном регистре" +
			"\n4. Шифр Атбаш"

		x.Send(text, &tb.SendOptions{ReplyTo: x.Message()}, s.TG.menu.textBtns)
	}

	switch s.TG.CallbackQuery[x.Chat().ID] {
	case "": //Нет в CallbackQuery - игнор
	case "test":

	}
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
	ar := &tb.ArticleResult{Title: "Текст в обратном порядке", Text: "<pre>" + text + "</pre>", Description: util.TextCut(text, 50)}
	ar.ParseMode = tb.ModeHTML
	res = append(res, ar)

	// Шифр Атбаша
	text = s.Other.AtbashAlphabet.Replace(q.Text)
	key := util.CreateKey(8)
	s.Other.AtbashCache.Add(key, q.Text)
	ar = &tb.ArticleResult{Title: "Кодировать шифром Атбаша", Text: "<pre>" + text + "</pre>", Description: util.TextCut(text, 50)}
	ar.ParseMode = tb.ModeHTML
	rm := *s.TG.menu.atbashBtns2
	rm.InlineKeyboard[0][0].Data = key
	ar.ReplyMarkup = &rm
	res = append(res, ar)

	// Текст в верхнем регистре
	text = strings.ToUpper(q.Text)
	ar = &tb.ArticleResult{Title: "Текст в верхнем регистре", Text: "<pre>" + text + "</pre>", Description: util.TextCut(text, 50)}
	ar.ParseMode = tb.ModeHTML
	res = append(res, ar)

	// Текст в случайном регистре
	text = textRandom(q.Text, s.Rand)
	ar = &tb.ArticleResult{Title: "Текст в случайном регистре", Text: "<pre>" + text + "</pre>", Description: util.TextCut(text, 50)}
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
	if !s.TG.isAdmin(x.Sender(), x.Chat().ID) {
		return
	}

	rm := &tb.ReplyMarkup{}
	btn := *s.TG.Buttons["test"]
	rm.Inline([]tb.Btn{btn})

	x.Send("Test", rm, tb.ModeHTML, tb.NoPreview)
	return
}

func (s *Service) TgTestBtn(x tb.Context) (errReturn error) {
	if !s.TG.isAdmin(x.Sender(), x.Message().Chat.ID) {
		return
	}

	rm := &tb.ReplyMarkup{}
	btn := *s.TG.Buttons["test"]
	rm.Inline([]tb.Btn{btn})

	x.Send("Test", &tb.SendOptions{ReplyTo: x.Message()}, rm, tb.ModeHTML, tb.NoPreview)
	x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "test"})
	return
}

func (s *Service) TgAdm(x tb.Context) (errReturn error) {
	if !s.TG.isAdmin(x.Sender(), x.Chat().ID) {
		return
	}

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

func (s *Service) TgStatusCMD(x tb.Context) (errReturn error) {
	if !s.TG.isAdmin(x.Sender(), x.Chat().ID) {
		return
	}

	text, rm := s.TgStatusFunc()
	x.Send(text, rm)
	return
}

func (s *Service) TgStatusUpdate(x tb.Context) (errReturn error) {
	if !s.TG.isAdmin(x.Sender(), x.Message().Chat.ID) {
		return
	}

	text, rm := s.TgStatusFunc()
	s.TG.Bot.Edit(x.Message(), text, rm)
	x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Обновлено"})
	return
}

func (s *Service) TgStatusFunc() (string, *tb.ReplyMarkup) {
	text := fmt.Sprintf("Uptime: %s\nAlbums manager length: %v",
		s.TG.uptimeString(s.TG.Uptime), s.TG.AlbumsManager.Len(),
	)

	rm := &tb.ReplyMarkup{}
	rm.Inline([]tb.Btn{*s.TG.Buttons["status_update"]})

	return text, rm
}

func (s *Service) TgLogsCMD(x tb.Context) (errReturn error) {
	if !s.TG.isAdmin(x.Sender(), x.Chat().ID) {
		return
	}

	text := "1. Получить файл логов\n2. Очистить файл логов"
	rm := &tb.ReplyMarkup{}
	rm.Inline(
		[]tb.Btn{*s.TG.Buttons["get_logs"], *s.TG.Buttons["clear_logs"]},
	)
	x.Send(text, rm, tb.ModeHTML)
	return
}

func (s *Service) TgGetLogsBtn(x tb.Context) (errReturn error) {
	if !s.TG.isAdmin(x.Sender(), x.Message().Chat.ID) {
		return
	}

	err := x.Send(&tb.Document{File: tb.FromDisk("files/logrus.log"), FileName: "logrus.log"})
	if err != nil {
		x.Send(errors.Wrap(err, "Ошибка отправки файла.").Error())
	}
	x.Respond()
	return
}

func (s *Service) TgClearLogsBtn(x tb.Context) (errReturn error) {
	if !s.TG.isAdmin(x.Sender(), x.Message().Chat.ID) {
		return
	}

	os.Truncate("files/logrus.log", 0)
	log.Info("Очищено")

	x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Очищено", ShowAlert: true})
	return
}

func (s *Service) TgDeleteBtn(x tb.Context) (errReturn error) {
	if !s.TG.isAdmin(x.Sender(), x.Message().Chat.ID) {
		return
	}
	x.Respond()
	s.TG.Bot.Delete(x.Message())
	x.Delete()
	return
}

func (s *Service) TgCancelReplyMarkup(x tb.Context) (errReturn error) {
	if !s.TG.isAdmin(x.Sender(), x.Chat().ID) {
		return
	}
	delete(s.TG.CallbackQuery, x.Chat().ID)
	x.Send("Отменено.", &tb.ReplyMarkup{RemoveKeyboard: true})
	return
}

func (s *Service) TgSetCmds(x tb.Context) (errReturn error) {
	if !s.TG.isAdmin(x.Sender(), x.Chat().ID) {
		return
	}

	x.Bot().SetCommands([]tb.Command{
		tb.Command{"help", "Список возможностей"},
	})

	x.Send("Сделано")
	return
}

/*

func (s *Service) TgBtn(x tb.Context) (errReturn error) {
	if !s.TG.isAdmin(x.Sender(), x.Message().Chat.ID) {
		return
	}

	//text, rm := s.TgFunc()
	//x.Send(text, rm, tb.ModeHTML)
	//s.TG.Bot.Edit(x.Message(), text, rm)
	x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: ""})
}
func (s *Service) TgCMD(x tb.Context) (errReturn error) {
	if !s.TG.isAdmin(x.Sender(), x.Chat().ID) {
		return
	}

	text, rm := s.TgInfoFilmFunc(m, 0)
	rm := &tb.ReplyMarkup{}
	rm.Inline(
		[]tb.Btn{*s.TG.Buttons["status_update"]},
	)
	x.Send(text, rm, tb.ModeHTML)
}

func (s *Service) TgFunc() (string, *tb.ReplyMarkup) {
	text := ""
	rm := &tb.ReplyMarkup{}
	rm.Inline(
		[]tb.Btn{*s.TG.Buttons["status_update"]},
	)
}

*/
