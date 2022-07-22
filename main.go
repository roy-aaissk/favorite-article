package main

import (
	"fmt"
	"log"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

const (
	baseUrl = "https://api.notion.com/v1/"
)

type NotionResposeDataBase struct {
	Object  string `json:"object"`
	Results []struct {
		Object         string    `json:"object"`
		ID             string    `json:"id"`
		CreatedTime    time.Time `json:"created_time"`
		LastEditedTime time.Time `json:"last_edited_time"`
		CreatedBy      struct {
			Object string `json:"object"`
			ID     string `json:"id"`
		} `json:"created_by"`
		LastEditedBy struct {
			Object string `json:"object"`
			ID     string `json:"id"`
		} `json:"last_edited_by"`
		Cover  interface{} `json:"cover"`
		Icon   interface{} `json:"icon"`
		Parent struct {
			Type       string `json:"type"`
			DatabaseID string `json:"database_id"`
		} `json:"parent"`
		Archived   bool `json:"archived"`
		Properties struct {
			URL struct {
				ID string `json:"id"`
			} `json:"URL "`
			Tag struct {
				ID string `json:"id"`
			} `json:"Tag"`
			Title struct {
				ID string `json:"id"`
			} `json:"Title"`
		} `json:"properties"`
		URL string `json:"url"`
	} `json:"results"`
	NextCursor interface{} `json:"next_cursor"`
	HasMore    bool        `json:"has_more"`
	Type       string      `json:"type"`
	Page       struct {
	} `json:"page"`
}

type Config struct {
	ApiKey string `json:"SECRET_API_KEY"`
	DatabaseId string `json:"NOTION_DATABASE_ID"`
	LineBotChannelSecret string `json:"LINE_BOT_CHANNEL_SECRET"`
	LineBotChannelToken string `json:"LINE_BOT_CHANNEL_TOKEN"`
}


// 環境変数読み込み
func LoadConfig(filePath string) (config *Config, err error) {
	content, err := ioutil.ReadFile(filePath)
  if err != nil {
		return nil, fmt.Errorf("error %s",err)
  }
	
  result := &Config{}
  if err := json.Unmarshal(content, result); err != nil {
		return nil, fmt.Errorf("read %s, %s", result, err)
  }
	
  return result, nil
}

func (c *Config) reqNotionAPI(url string, method string, payload *strings.Reader) (r *http.Response, err error){
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, fmt.Errorf("not request %s", err)
	}
	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("Authorization", c.ApiKey )
	client := &http.Client{}
	res, err := client.Do(req)
	return res, err
}

func main() {
	config, err := LoadConfig("config.json")
  if err != nil {
		log.Fatal(err)
  }
	// page情報一覧取得
	dataBaseUrl := baseUrl + "databases/" + config.DatabaseId + "/query"
	method := "POST"
	payload := strings.NewReader("{\"property\":Tag}")
	res, err := config.reqNotionAPI(dataBaseUrl, method, payload)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var result NotionResposeDataBase
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("%+v", err)
	}

	for _, v := range result.Results {
		fmt.Printf(" result ID name: ",v.ID)
	}


	// page情報から情報をランダムで取得
	// page_idからランダムで情報取得が良い


	// LineBotAPI接続設定
	bot, err := linebot.New(
		config.LineBotChannelSecret,
		config.LineBotChannelToken,
	)
	if err != nil {
		log.Fatal(err, bot)
	}
	fmt.Printf("logmessage: LINE BOT : bot %s", bot)
}



