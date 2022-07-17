package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"net/http"
	"io/ioutil"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	// notionAPIdatabase接続設定
	databaseId := os.Getenv("NOTION_DATABASE_ID")
	APIkey := os.Getenv("SECRET_API_KEY")
	databaseurl := "https://api.notion.com/v1/databases/" + databaseId + "/query"
	// databaseからpage情報一覧を取得
	payload := strings.NewReader("{\"property\":Tag}")
	req, err := http.NewRequest("POST", databaseurl, payload)
	if err != nil {
		log.Fatal(err)
	}
	// requestHeaderAPIバージョンとkey情報を追加
	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("Authorization", APIkey )
	fmt.Printf("set databaseId APIkey %s", APIkey)
	
	client := &http.Client{}
	res, _ := client.Do(req)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("response Body: %s", body)

	// LineBotAPI接続設定
	bot, err := linebot.New(
		os.Getenv("LINE_BOT_CHANNEL_SECRET"),
		os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err, bot)
	}
	fmt.Printf("LINE BOT : bot %s", bot)
}
