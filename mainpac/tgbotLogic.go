package mainpac

import (
	"Loren/go-log"
	"fmt"
	"github.com/pkg/errors"
	tb "gopkg.in/tucnak/telebot.v2"
	"os"
)

/*
CMD -        команда
Update/Btn - кнопка обновления/обычная
Func -       логика работы
Но они обязательны только все вместе
*/

func (s *Service) TgStartCMD(m *tb.Message) {
	s.TG.Bot.Send(m.Sender, "Hello World!", &tb.ReplyMarkup{ReplyKeyboardRemove: true})
}

// Ниже только админское

func (s *Service) TgTest(m *tb.Message) {
	if !s.TG.isAdmin(m.Sender, m.Chat.ID) {
		return
	}

	rm := &tb.ReplyMarkup{}
	btn := *s.TG.Buttons["test"]
	rm.Inline([]tb.Btn{btn})

	s.TG.Bot.Send(m.Sender, "Test", rm, tb.ModeHTML, tb.NoPreview)
}

func (s *Service) TgTestBtn(c *tb.Callback) {
	if !s.TG.isAdmin(c.Sender, c.Message.Chat.ID) {
		return
	}

	rm := &tb.ReplyMarkup{}
	btn := *s.TG.Buttons["test"]
	rm.Inline([]tb.Btn{btn})

	s.TG.Bot.Send(c.Sender, "Test", &tb.SendOptions{ReplyTo: c.Message}, rm, tb.ModeHTML, tb.NoPreview)
	s.TG.Bot.Respond(c, &tb.CallbackResponse{CallbackID: c.ID, Text: "test"})
}

func (s *Service) TgAdm(m *tb.Message) {
	if !s.TG.isAdmin(m.Sender, m.Chat.ID) {
		return
	}

	text := fmt.Sprintf("" +
		"<b>Пользователькие команды:</b>\n" +
		"/start - приветствие\n" +
		"\n<b>Админские команды:</b>\n" +
		"/status - статус работы\n" +
		"/logs - действия над логами",
	)

	s.TG.Bot.Send(m.Sender, text, tb.ModeHTML)
}

func (s *Service) TgStatusCMD(m *tb.Message) {
	if !s.TG.isAdmin(m.Sender, m.Chat.ID) {
		return
	}

	text, rm := s.TgStatusFunc()
	s.TG.Bot.Send(m.Sender, text, rm)
}

func (s *Service) TgStatusUpdate(c *tb.Callback) {
	if !s.TG.isAdmin(c.Sender, c.Message.Chat.ID) {
		return
	}

	text, rm := s.TgStatusFunc()
	s.TG.Bot.Edit(c.Message, text, rm)
	s.TG.Bot.Respond(c, &tb.CallbackResponse{CallbackID: c.ID, Text: "Обновлено"})
}

func (s *Service) TgStatusFunc() (string, *tb.ReplyMarkup) {
	text := fmt.Sprintf("Uptime: %s", s.TG.uptimeString(s.TG.Uptime))

	rm := &tb.ReplyMarkup{}
	rm.Inline([]tb.Btn{*s.TG.Buttons["status_update"]})

	return text, rm
}

func (s *Service) TgLogsCMD(m *tb.Message) {
	if !s.TG.isAdmin(m.Sender, m.Chat.ID) {
		return
	}

	text := "1. Получить файл логов\n2. Очистить файл логов"
	rm := &tb.ReplyMarkup{}
	rm.Inline(
		[]tb.Btn{*s.TG.Buttons["get_logs"], *s.TG.Buttons["clear_logs"]},
	)
	s.TG.Bot.Send(m.Sender, text, rm, tb.ModeHTML)
}

func (s *Service) TgGetLogsBtn(c *tb.Callback) {
	if !s.TG.isAdmin(c.Sender, c.Message.Chat.ID) {
		return
	}

	_, err := s.TG.Bot.Send(c.Sender, &tb.Document{File: tb.FromDisk("files/logrus.log"), FileName: "logrus.log"})
	if err != nil {
		s.TG.Bot.Send(c.Sender, errors.Wrap(err, "Ошибка отправки файла.").Error())
	}
	s.TG.Bot.Respond(c)
}

func (s *Service) TgClearLogsBtn(c *tb.Callback) {
	if !s.TG.isAdmin(c.Sender, c.Message.Chat.ID) {
		return
	}

	file, _ := os.OpenFile("files/logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	file.Truncate(0)
	file.Close()
	log.Info("Очищено")

	s.TG.Bot.Respond(c, &tb.CallbackResponse{CallbackID: c.ID, Text: "Очищено", ShowAlert: true})
}

func (s *Service) TgCallbackQuery(m *tb.Message) {
	if !s.TG.isAdmin(m.Sender, m.Chat.ID) {
		return
	}

	switch s.TG.CallbackQuery[m.Chat.ID] {
	case "": //Нет в CallbackQuery - игнор
	case "test":

	}
}

func (s *Service) TgDeleteBtn(c *tb.Callback) {
	if !s.TG.isAdmin(c.Sender, c.Message.Chat.ID) {
		return
	}
	s.TG.Bot.Respond(c)
	s.TG.Bot.Delete(c.Message)
}

func (s *Service) TgCancelReplyMarkup(m *tb.Message) {
	if !s.TG.isAdmin(m.Sender, m.Chat.ID) {
		return
	}
	delete(s.TG.CallbackQuery, m.Chat.ID)
	s.TG.Bot.Send(m.Sender, "Отменено.", &tb.ReplyMarkup{ReplyKeyboardRemove: true})
}

/*

func (s *Service) TgBtn(c *tb.Callback) {
	if !s.TG.isAdmin(c.Sender, c.Message.Chat.ID) {
		return
	}

	//text, rm := s.TgFunc()
	//s.TG.Bot.Send(c.Sender, text, rm, tb.ModeHTML)
	//s.TG.Bot.Edit(c.Message, text, rm)
	s.TG.Bot.Respond(c, &tb.CallbackResponse{CallbackID: c.ID, Text: ""})
}
func (s *Service) TgCMD(m *tb.Message) {
	if !s.TG.isAdmin(m.Sender, m.Chat.ID) {
		return
	}

	text, rm := s.TgInfoFilmFunc(m, 0)
	rm := &tb.ReplyMarkup{}
	rm.Inline(
		[]tb.Btn{*s.TG.Buttons["status_update"]},
	)
	s.TG.Bot.Send(m.Sender, text, rm, tb.ModeHTML)
}

func (s *Service) TgFunc() (string, *tb.ReplyMarkup) {
	text := ""
	rm := &tb.ReplyMarkup{}
	rm.Inline(
		[]tb.Btn{*s.TG.Buttons["status_update"]},
	)
}

*/
