package mainpac

import tb "gopkg.in/tucnak/telebot.v3"

func (s *Service) TgTextReverse(x tb.Context) (errReturn error) {
	defer x.Respond()
	if x.Message().ReplyTo == nil {
		x.Respond(&tb.CallbackResponse{CallbackID: x.Callback().ID, Text: "Ошибка", ShowAlert: true})
		return
	}
	text := x.Message().ReplyTo.Text

	xlist := []rune(text)
	ln := len(xlist)
	xlistRevers := make([]rune, 0, ln)
	for i := 0; i < ln; i++ {
		xlistRevers = append(xlistRevers, xlist[ln-i-1])
	}
	textRevers := string(xlistRevers)

	x.Send(textRevers)
	return
}
