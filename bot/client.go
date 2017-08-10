package bot

import (
	"bot/comm"
)

func GetUpdates(offset int) *comm.Updates {
	return comm.GetUpdates(offset)
}

func UpdateMessage(reply *comm.Reply) error {
	return comm.Update("editMessageText", reply)
}

func SendMessage(reply *comm.Reply) (int, error) {
	return comm.Send("sendMessage", reply)
}

func AnswerInlineQuery(a *comm.InlineQueryAnswer) error {
	return comm.Update("answerInlineQuery", a)
}

func deleteMessage(d *comm.DeleteMessage) error {
	return comm.Update("deleteMessage", d)
}
