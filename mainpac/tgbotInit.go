package mainpac

import tb "gopkg.in/tucnak/telebot.v2"

func (s *Service) InitTBot() {
	// Команды роуты

	s.TG.Bot.Handle("/start", s.TgStartCMD)

	// Админские команды

	s.TG.Bot.Handle("/test", s.TgTest)
	s.TG.Bot.Handle("/adm", s.TgAdm)
	s.TG.Bot.Handle("/status", s.TgStatusCMD)
	s.TG.Bot.Handle("/logs", s.TgLogsCMD)
	s.TG.Bot.Handle(tb.OnText, s.TgCallbackQuery)

	// Кнопки роуты

	rm := &tb.ReplyMarkup{}
	im := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	s.TG.Buttons = make(map[string]*tb.Btn)

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

*/
