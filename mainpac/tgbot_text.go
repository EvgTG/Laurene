package mainpac

import tb "gopkg.in/tucnak/telebot.v3"

func (s *Service) TgTextReverse(x tb.Context) (errReturn error) {
	defer x.Respond()
	if x.Message().ReplyTo == nil {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка", ShowAlert: true})
		return
	}

	r := []rune(x.Message().ReplyTo.Text)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}

	x.Send(string(r))
	return
}
