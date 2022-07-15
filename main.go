package main

import (
	"fmt"
	"os"
	"log"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	// LineBotAPI接続設定
	bot, err := linebot.New(
		os.Getenv("LINE_BOT_CHANNEL_SECRET"),
		os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	
	fmt.Println(bot)
}
