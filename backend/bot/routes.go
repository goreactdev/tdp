package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (app *application) routes(msgConfig tgbotapi.MessageConfig, updateMsg *tgbotapi.Message) error {

	command := updateMsg.Command()
	
	switch command {
	case "start":
		err := app.startHandler(msgConfig, updateMsg)
		if err != nil {
			return err
		}

	case "rating":
		err := app.ratingHandler(msgConfig, updateMsg)
		if err != nil {
			return err
		}

	// Process the command and parameter as needed
	case "whois":
		err := app.whoisHandler(msgConfig, updateMsg)
		if err != nil {
			return err
		}

	case "help":
		err := app.helpHandler(msgConfig, updateMsg)
		if err != nil {
			return err
		}

	default:
		err := app.scoreHandler(msgConfig, updateMsg)
		if err != nil {
			return err
		}
	}

	return nil
}
