package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

const (
	baseUrl = "https://api.notion.com/v1/"
)

type NotionResponseDataBase struct {
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

type NotionPageTitleList struct {
	Object  string `json:"object"`
	Results []struct {
		Object string `json:"object"`
		Type   string `json:"type"`
		ID     string `json:"id"`
		Title  struct {
			Type string `json:"type"`
			Text struct {
				Content string      `json:"content"`
				Link    interface{} `json:"link"`
			} `json:"text"`
			Annotations struct {
				Bold          bool   `json:"bold"`
				Italic        bool   `json:"italic"`
				Strikethrough bool   `json:"strikethrough"`
				Underline     bool   `json:"underline"`
				Code          bool   `json:"code"`
				Color         string `json:"color"`
			} `json:"annotations"`
			PlainText string      `json:"plain_text"`
			Href      interface{} `json:"href"`
		} `json:"title"`
	} `json:"results"`
	NextCursor   interface{} `json:"next_cursor"`
	HasMore      bool        `json:"has_more"`
	Type         string      `json:"type"`
	PropertyItem struct {
		ID      string      `json:"id"`
		NextURL interface{} `json:"next_url"`
		Type    string      `json:"type"`
		Title   struct {
		} `json:"title"`
	} `json:"property_item"`
}

type NotionPageUrlList struct {
	Object string `json:"object"`
	Type   string `json:"type"`
	ID     string `json:"id"`
	URL    string `json:"url"`
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

func (c *Config) reqNotionAPI(url string, method string, payload io.Reader) (b []byte, err error){
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, fmt.Errorf("not request %s", err)
	}
	req.Header.Add("Notion-Version", "2022-06-28")
	req.Header.Add("Authorization", c.ApiKey )
	client := &http.Client{}
	res, err := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("not found res.Body: %s", err)
	}
	defer res.Body.Close()
	return body, err
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World from Go.")
}

func main() {
	config, err := LoadConfig("config.json")
  if err != nil {
		log.Fatal(err)
  }
	// page情報一覧取得
	resBody, err := config.reqNotionAPI(baseUrl + "databases/" + config.DatabaseId + "/query", "POST", strings.NewReader("{\"property\":Tag}"))

	var result NotionResponseDataBase
	if err := json.Unmarshal(resBody, &result); err != nil {
		log.Printf("not found result: %s", err)
	}

	// pageListID取得する
	id :=  make(map[int]string)
	var titleId string
	var urlId string
	for a, v := range result.Results {
		// page情報から情報を5つ取得する
		if a >= 5{
			break
		}
		titleId = v.Properties.Title.ID
		urlId = v.Properties.URL.ID
		id[a] = v.ID
		fmt.Printf("pageList %v", id[a])
	}
	pageListCount := len(id)
	if pageListCount == 0 {
		log.Printf("not found pageList")
	}

	// 詳細項目取得
	var title NotionPageTitleList
	var url  NotionPageUrlList
	var resTitle []byte
	var resUrl []byte
	var titleText string
	n := 0
	for n < 5{
		resTitle, err = config.reqNotionAPI(baseUrl + "pages/" + id[n] + "/properties/" + titleId, "GET", nil)
		resUrl, err = config.reqNotionAPI(baseUrl + "pages/" + id[n] + "/properties/" + urlId, "GET", nil)
		n++
	}
	// 取得したtitleとurlをjson形式に変換
  if err := json.Unmarshal(resTitle, &title); err != nil  {
		log.Printf("not found result: %s", err)
	}
  if err := json.Unmarshal(resUrl, &url); err != nil  {
		log.Printf("not found result: %s", err)
	}
	for _, v := range title.Results {
		titleText = v.Title.Text.Content
	}
	// タイトルテキスト
	fmt.Print(titleText)
	// URLテキスト
	fmt.Print(url.URL)
	// LineBotAPI接続設定
	bot, err := linebot.New(
		config.LineBotChannelSecret,
		config.LineBotChannelToken,
	)
	if err != nil {
		log.Fatal(err, bot)
	}
	fmt.Printf(" logmessage: LINE BOT : bot %s", bot)

	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/bot/webhook", HandleRequest)
	server.ListenAndServe()
}



