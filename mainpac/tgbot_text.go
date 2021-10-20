package mainpac

import tb "gopkg.in/tucnak/telebot.v3"

func (s *Service) TgTextReverse(x tb.Context) (errReturn error) {
	defer x.Respond()
	if x.Message().ReplyTo == nil {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка", ShowAlert: true})
		return
	}

	x.Send(textReverse(x.Message().ReplyTo.Text))
	return
}

func textReverse(s string) (res string) {
	r := []rune(s)
	for _, v := range r {
		res = string(v) + res
	}
	return
}
