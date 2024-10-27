package main

//
//import (
//	"fmt"
//	"github.com/mymmrac/telego"
//	th "github.com/mymmrac/telego/telegohandler"
//	tu "github.com/mymmrac/telego/telegoutil"
//)
//
//func SendKeyboardButton(bot *telego.Bot, update telego.Update) {
//	chatID := tu.ID(update.Message.Chat.ID)
//
//	keyboard := tu.Keyboard(
//		tu.KeyboardRow(
//			tu.KeyboardButton("Start"),
//			tu.KeyboardButton("Help"),
//		))
//
//	message := tu.Message(
//		chatID,
//		"Keyboard",
//	).WithReplyMarkup(keyboard)
//
//	bot.SendMessage(message)
//}
//
//func SendAnyMessage(bot *telego.Bot, update telego.Update) {
//	chatID := tu.ID(update.Message.Chat.ID)
//
//	bot.CopyMessage(&telego.CopyMessageParams{
//		ChatID:     chatID,
//		FromChatID: chatID,
//		MessageID:  update.Message.MessageID,
//	})
//}
//
//func main() {
//
//	botToken := "7859455978:AAHhbZLsj8IH5HQ3v5YiIuRAJG6iVSy17oA"
//
//	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	updates, _ := bot.UpdatesViaLongPolling(nil)
//	bh, _ := th.NewBotHandler(bot, updates)
//
//	defer bh.Stop()
//	defer bot.StopLongPolling()
//
//	bh.Handle(SendKeyboardButton, th.CommandEqual("start"))
//	bh.Handle(SendAnyMessage, th.AnyMessage())
//
//	bh.Start()
//
//}
