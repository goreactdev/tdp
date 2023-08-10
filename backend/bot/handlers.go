package main

import (
	"encoding/json"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hibiken/asynq"
	"github.com/ton-developer-program/internal/database"
)

// start handler
func (app *application) startHandler(msgConfig tgbotapi.MessageConfig, updateMsg *tgbotapi.Message) error {

	// if it's group return nil
	if updateMsg.Chat.Type == "group" {
		return nil
	}

	msgConfig.Text = `
	👋 Hello, I'm your personal bot for TON Developers Platform!

		I will help you stay up-to-date with your new rewards and all the events for developers in TON. I'm also available in every chat from <a href='https://t.me/c/1913693703/66'>TON Dev Kit</a>.

		Commands:

		/rating - check your current rating

		/whois - get information about a chat participant via reply

		/whois [@username] - get information about a participant via Telegram username

		/help - display the list of commands	
	`

	msgConfig.ParseMode = "HTML"

	if _, err := app.bot.Send(msgConfig); err != nil {
		return err
	}

	return nil
}

// command help

// command help
func (app *application) helpHandler(msgConfig tgbotapi.MessageConfig, updateMsg *tgbotapi.Message) error {
	// if it's group return nil
	if updateMsg.Chat.Type == "group" {
		return nil
	}

	msgConfig.Text = `
	Commands:

	/rating - check your current rating

	/whois - get information about a chat participant via reply

	/whois [@username] - get information about a participant via Telegram username

	/help - display the list of commands	
	`

	if _, err := app.bot.Send(msgConfig); err != nil {
		return err
	}

	return nil

}

func (app *application) ratingHandler(msgConfig tgbotapi.MessageConfig, updateMsg *tgbotapi.Message) error {

	user, err := app.sqlModels.Users.GetByTelegramUserId(updateMsg.From.ID)
	if err != nil {
		return err
	}

	if user == nil {
		app.welcomeMessage(msgConfig)
		return nil
	}

	// // get last award of user
	name, friendlyAddr, weight, err := app.sqlModels.Nfts.GetLastTokenCreated(user.FriendlyAddress)
	if err != nil {
		return err
	}

	// get position of user
	position, allUsers, err := app.sqlModels.Users.GetUserPosition(user.ID)
	if err != nil {
		return err
	}

	var lastRewardText string

	if name != "" {
		var url = fmt.Sprintf("https://getgems.io/nft/%s", friendlyAddr)

		lastRewardText = fmt.Sprintf(`<a href="%s">Last reward: %s (+%d)</a>`, url, name, weight)
	} else {
		lastRewardText = fmt.Sprintf("Last reward: none")
	}

	msgConfig.Text = fmt.Sprintf(`
	🏆 TON Developers Leaderboard 🏆

		▪️ Username: %s
		▪️ Position: %d out of %d
		▪️ Rating: %d 💎
		
		%s
		`, user.Username, position, allUsers, int64(user.Rating), lastRewardText)

	// parse mode html
	msgConfig.ParseMode = "HTML"

	// add button
	msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("🔗 Open profile", "https://tdp.tonbuilders.com/user/"+user.Username),
		),
	)

	// remove preview
	msgConfig.DisableWebPagePreview = true
	if _, err := app.bot.Send(msgConfig); err != nil {
		return err
	}

	return nil
}

// command sayhi
func (app *application) whoisHandler(msgConfig tgbotapi.MessageConfig, updateMsg *tgbotapi.Message) error {

	incomingUser, err := app.sqlModels.Users.GetByTelegramUserId(updateMsg.From.ID)
	if err != nil {
		return err
	}

	if incomingUser == nil {
		app.welcomeMessage(msgConfig)
		return nil
	}

	var user *database.User

	if updateMsg.CommandArguments() != "" {
		user, err = app.sqlModels.Users.GetByTelegramUsername(updateMsg.CommandArguments()[1:])
		if err != nil {
			return err
		}

	}

	// ic ommandArguments is "" and replyToMessage.From.ID is not nil
	if user == nil && updateMsg.ReplyToMessage != nil && updateMsg.ReplyToMessage.From.ID != 0 {
		user, err = app.sqlModels.Users.GetByTelegramUserId(updateMsg.ReplyToMessage.From.ID)
		if err != nil {
			return err
		}

	}

	if user == nil {
		return nil
	}

	// // get last award of user
	name, friendlyAddr, weight, err := app.sqlModels.Nfts.GetLastTokenCreated(user.FriendlyAddress)
	if err != nil {
		return err
	}

	// get position of user
	position, allUsers, err := app.sqlModels.Users.GetUserPosition(user.ID)
	if err != nil {
		return err
	}

	var lastRewardText string

	if name != "" {
		var url = fmt.Sprintf("https://getgems.io/nft/%s", friendlyAddr)

		lastRewardText = fmt.Sprintf(`<a href="%s">Last reward: %s (+%d)</a>`, url, name, weight)
	} else {
		lastRewardText = fmt.Sprintf("Last reward: none")
	}

	msgConfig.Text = fmt.Sprintf(`
	🏆 TON Developers Platform 🏆

		▪️ Username: %s
		▪️ Position: %d out of %d
		▪️ Rating: %d 💎
		
		%s
		`, user.Username, position, allUsers, int64(user.Rating), lastRewardText)

	// parse mode html

	msgConfig.ParseMode = "HTML"

	msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("🔗 Open profile", "https://tdp.tonbuilders.com/user/"+user.Username),
		),
	)

	msgConfig.DisableWebPagePreview = true

	if _, err := app.bot.Send(msgConfig); err != nil {
		return err
	}

	return nil
}

// user's score handler
func (app *application) scoreHandler(msgConfig tgbotapi.MessageConfig, updateMsg *tgbotapi.Message) error {

	//  only if the chat is negative number
	if updateMsg.Chat.ID > 0 {
		return nil
	}

	// get user by telegram user id
	user, err := app.sqlModels.Users.GetByTelegramUserId(updateMsg.From.ID)
	if err != nil {
		return err
	}

	if user == nil {
		return nil
	}

	telegramMessage := &database.TelegramMessage{
		UserID:    updateMsg.From.ID,
		MessageID: updateMsg.MessageID,
		ChatID:    updateMsg.Chat.ID,
	}

	payload, err := json.Marshal(telegramMessage)

	if err != nil {
		app.logger.Error(err, nil)
		return err
	}

	runGetTelegramMessageQueue := asynq.NewTask(database.TYPE_ADD_TG_MESSAGE, payload)

	info, err := app.asynqClient.Enqueue(runGetTelegramMessageQueue, asynq.ProcessIn(20*time.Second), asynq.MaxRetry(5), asynq.ProcessIn(5*time.Second), asynq.Retention(10*time.Minute), asynq.Queue(database.PRIORITY_NORMAL))
	if err != nil {
		app.logger.Error(fmt.Errorf("error: %s", err), nil)
		return err
	}

	prototypesNFT, err := app.sqlModels.Nfts.GetPrototypesByRating(user.ID)
	if err != nil {
		app.logger.Warning(fmt.Sprintf("error running get user rating handler: %v", err))
		return err
	}

	if len(prototypesNFT) > 0 {
		msg := fmt.Sprintf("🎉 Congratulations %s! You have earned new SBT token! 🎉\n\n", user.Username)

		for _, prototype := range prototypesNFT {
			nft, err := app.sqlModels.Nfts.GetSbtTokenByContentUri("https://tdp.tonbuilders.com/v1/deployed-nft/n/" + prototype.Base64 + "/meta.json")
			if err != nil {
				return err
			}

			rewardText := fmt.Sprintf(`<a href="%s">%s</a>`, nft.ContentUri, *nft.Name)

			msg += fmt.Sprintf(`
			
			🏆 Hey Roman! Here is your reward:

			▪️ SBT: %s
			▪️ Description: %s
			▪️ Rating: +%d points

			Keep up the great work! Enjoy!

			`,
			rewardText, prototype.Description, nft.Weight)
		}

		// inline keyboard
		msgConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("🔗 See In My Profile", "https://tdp.tonbuilders.com/user/"+user.Username),
			),
		)

		msgConfig.Text = msg

		if _, err := app.bot.Send(msgConfig); err != nil {
			return err
		}
	}

	app.logger.Info(fmt.Sprintf("enqueued task with id %s", info.ID))

	return nil
}
