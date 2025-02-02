package main

import (
    "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "log"
    "os"
    "client/client"
)

/*
TODO:
1. Implement bot subscription which would send notifications about upcoming events to user
*/

func main() {

    _, status := os.LookupEnv("TG_TOKEN")
    if status == false {
        log.Printf("TG_TOKEN env is missing.")
        os.Exit(1)
    }

    bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
    if err != nil {
        log.Panic(err)
    }
    bot.Debug = false

    // Create chan for telegram updates
    var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
    ucfg.Timeout = 60
    updates := bot.GetUpdatesChan(ucfg)
    for update := range updates {
        if update.Message == nil { // ignore any non-Message updates
            continue
        }

        if !update.Message.IsCommand() { // ignore any non-command Messages
            continue
        }

        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
        msg.Text = client.FetchData(update.Message.Command())
        msg.ParseMode = "markdown"
        if len(msg.Text) == 0 {
            log.Printf("Failed to fetch data for user request")
            msg.Text = "Failed to fetch data for some reason"
        }
        if _, err := bot.Send(msg); err != nil {
            log.Panic(err)
        }
    }
}
