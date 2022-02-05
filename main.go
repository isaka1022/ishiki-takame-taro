package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/joho/godotenv"
)

type Content struct {
	Content string `json:"content"`
	Link    string `json:"link"`
}

type Annotations struct {
	Bold          bool   `json:"bold"`
	Italic        bool   `json:"italic"`
	StrikeThrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Code          bool   `json:"code"`
	Color         string `json:"color"`
}

type BulletedListItemContent struct {
	Type        string                     `json:"type"`
	Text        Content                    `json:"text"`
	Children    []*BulletedListItemContent `json:"children"`
	Annotations Annotations                `json:"annotations"`
	PlainText   string                     `json:"plain_text"`
	Href        string                     `json:"href"`
}

type BulletedListItem struct {
	Text []*BulletedListItemContent
}

type ChildPage struct {
	Text []*BulletedListItemContent
}

type Block struct {
	Object           string `json:"object"`
	Id               string `json:"id"`
	CreatedTime      string `json:"created_time"`
	LastEditedTime   string `json:"last_edited_time"`
	HasChildren      bool   `json:"has_children"`
	Archived         bool   `json:"archived"`
	Type             string `json:"type"`
	BulletedListItem `json:"bulleted_list_item"`
	ChildPage `json:"child_page"`
}

type Blocks []*Block

type Body struct {
	Object     string `json:"object"`
	Results    Blocks `json:"results"`
	NextCursor string `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
}

func main() {
	godotenv.Load(".env")
	PortNum := os.Getenv("PORT")

	http.HandleFunc("/callback", lineHandler)

	fmt.Println("https://localhost:" + PortNum + "で起動中...")

	log.Fatal(http.ListenAndServe(":"+PortNum, nil))
}

// 中身を取得する
func GetContents(blockId string, isOnlyId bool) []string {
	var ids []string
	var messages []string
	var body Body
	var b = FetchChild(blockId)
	json.Unmarshal(b, &body)

	fmt.Println("into GetContents")

	for _, block := range body.Results {
		ids = append(ids, block.Id)
		var childContents []string
		for _, content := range block.BulletedListItem.Text {
			if block.HasChildren == true && !isOnlyId {
				childContents = GetContents(block.Id, false)
			}
			messages = append([]string{content.PlainText}, childContents...)
		}
	}
	if isOnlyId {
		return ids
	}
	fmt.Println("messages")
	fmt.Println(messages)
	return messages
}

func FetchChild(BlockId string) []byte {
	ApiKey := os.Getenv("NOTION_SECRET_KEY")

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.notion.com/v1/blocks/"+BlockId+"/children?page_size=100", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+ApiKey)
	req.Header.Add("Notion-Version", "2021-08-16")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return b
}

// ランダムで一つ選ぶ
func SelectId(ids []string) string {
	fmt.Println("into SelectId")

	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(len(ids))
	return ids[num]
}

func FormatMessage(id string) string {
	fmt.Println("into FormatMessage")

	message := GetContents(id, false)
	fmt.Println("before return")
	return strings.Join(message, "\n・")
}

func lineHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("called")
	SecretToken := os.Getenv("LINE_CANNEL_SECRET_TOKEN")
	AccessToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	DatabaseId := os.Getenv("NOTION_DATABASE_ID")
	bot, err := linebot.New(
		SecretToken,
		AccessToken,
	)
	if err != nil {
		log.Fatal(err)
	}

	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch event.Message.(type) {
			case *linebot.TextMessage:
				// _, err = bot.pushMessage(event.ReplyToken, linebot.NewTextMessage("考え中...")).Do()
				if err != nil {
					log.Print(err)
				}
				pageIds := GetContents(DatabaseId, true)
				fmt.Println("pageIds")
				fmt.Println(pageIds)
				pageId := SelectId(pageIds)
				fmt.Println("pageId")
				fmt.Println(pageId)
				blockIds := GetContents(pageId, true)
				fmt.Println("blockIds")
				fmt.Println(blockIds)

				blockId := SelectId(blockIds)
				fmt.Println("blockId")
				fmt.Println(blockId)
				sendMessage := FormatMessage(blockId)
				_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(sendMessage)).Do()
				if err != nil {
					log.Print(err)
					bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("失敗しました。")).Do()
				}
			}
		}
	}
}
