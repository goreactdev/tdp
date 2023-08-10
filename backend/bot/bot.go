package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (app *application) startBot() error {

	app.logger.Info("Authorized on account %s", app.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := app.bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	// var allowedChats = 
	// 	-592623628 

	// create a channel to handle OS signals
	quit := make(chan os.Signal, 1)
	// notify the quit channel for multiple signals
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// run the main loop in a goroutine so that it can be stopped gracefully
	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}

			// check if chat id is negative number
			if update.Message.Chat.ID < 0 && app.config.App.AlloweGroupChatID != update.Message.Chat.ID {
                    continue
			}
			

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

				err := app.routes(msg, update.Message)
				if err != nil {
					log.Panic(err)
				}
			}

	}()

	// listen to OS signals
	<-quit

	app.logger.Info("Stopping bot...")
	// close the updates channel
	app.bot.StopReceivingUpdates()
	// wait for the main loop to finish
	app.logger.Info("Bot stopped")

	return nil
}
