package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (app *application) welcomeMessage(msg tgbotapi.MessageConfig) {
	msg.Text = `
	🏆 TON Developers Platform 🏆

	Hey! You don't have an account yet. 🙂

	Join now, become #1 in the TON Community and get <a href="https://ton-org.notion.site/How-to-get-rewarded-ad8ab607478d4a7ab8658051d4ce5bf7">unique rewards</a> in TON merch store with your rating!
	`

	msg.ParseMode = "HTML"

	// Add button to the message
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("💎 Join Developers Platform", "https://tdp.tonbuilders.com/"),
		),
	)


	msg.DisableWebPagePreview = true

	if _, err := app.bot.Send(msg); err != nil {
		log.Panic(err)
	}

}