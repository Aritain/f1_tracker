package main

import (
    "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "log"
    "os"
    "client/client"
)

func main() {

    tgToken, status := os.LookupEnv("TG_TOKEN")
    if !status {
        log.Printf("TG_TOKEN env is missing.")
        os.Exit(1)
    }

    notificationToggle, status := os.LookupEnv("NOTIFICATION_TOGGLE")
    if !status {
        log.Printf("NOTIFICATION_TOGGLE env is missing.")
        os.Exit(1)
    }

    bot, err := tgbotapi.NewBotAPI(tgToken)
    if err != nil {
        log.Panic(err)
    }
    bot.Debug = true

    // Create chan for telegram updates
    var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
    ucfg.Timeout = 60
    updates := bot.GetUpdatesChan(ucfg)

    if notificationToggle == "true" {
        go client.AssetWatcher(bot)
    }

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
