package main

import (
    "github.com/Syfaro/telegram-bot-api"
    "log"
    "os"
    "client/client"
)


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
    bot.Debug = true

    // Create chan for telegram updates
    var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
    ucfg.Timeout = 60
    updates, _ := bot.GetUpdatesChan(ucfg)
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


        if _, err := bot.Send(msg); err != nil {
            log.Panic(err)
        }
    }
}
