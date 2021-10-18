package mainpac

import (
	"Laurene/go-log"
	"fmt"
	"github.com/pkg/errors"
	tb "gopkg.in/tucnak/telebot.v3"
	"os"
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
		"\n• Склейка фото, если прислать или переслать альбом." +
		"\n• Счёт дат в сообщении информации о пользователе из @YetAnotherBot (переслать)"

	x.Send(text, &tb.ReplyMarkup{RemoveKeyboard: true})
	return
}

func (s *Service) TgOnText(x tb.Context) (errReturn error) {
	//if !s.TG.isAdmin(x.Sender(), x.Chat().ID) {
	//	return
	//}

	switch {
	case s.Other.YetAnotherBotInfoUserRGX.MatchString(x.Text()):
		s.TgInfoUserYAB(x)
		return
	default:
		text := "" +
			"Что сделать с текстом?" +
			"\n" +
			"\n1. Написать в обратном порядке"

		x.Send(text, &tb.SendOptions{ReplyTo: x.Message()}, s.TG.menu.textBtns)
	}

	switch s.TG.CallbackQuery[x.Chat().ID] {
	case "": //Нет в CallbackQuery - игнор
	case "test":

	}
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
		"\n/logs - действия над логами",
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
