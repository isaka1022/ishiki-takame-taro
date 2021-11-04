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

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
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

type Block struct {
	Object           string `json:"object"`
	Id               string `json:"id"`
	CreatedTime      string `json:"created_time"`
	LastEditedTime   string `json:"last_edited_time"`
	HasChildren      bool   `json:"has_children"`
	Archived         bool   `json:"archived"`
	Type             string `json:"type"`
	BulletedListItem `json:"bulleted_list_item"`
}

type Blocks []*Block

type Body struct {
	Object     string `json:"object"`
	Results    Blocks `json:"results"`
	NextCursor string `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
}

type Title struct {
	Name     string
	Contents []string
}

func main() {

	http.HandleFunc("/callback", lineHandler)

	fmt.Println("https://localhost:8080 で起動中...")

	log.Fatal(http.ListenAndServe(":8080", nil))

	// var out bytes.Buffer
	// json.Indent(&out, b, "", " ")
	// out.WriteTo(os.Stdout)
}

func ShowTitles(titles []Title) string {

	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(len(titles))
	ShowTitle := titles[num]
	texts := ShowTitle.Name
	for _, content := range ShowTitle.Contents {
		texts += content
	}

	return texts
}

func FetchTitles() []Title {
	var titles []Title

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	DatabaseId := os.Getenv("NOTION_DATABASE_ID")

	var body Body
	var b = FetchNotion(DatabaseId)
	json.Unmarshal(b, &body)

	for _, block := range body.Results {
		var ContentArray []string
		if block.HasChildren == true {
			FetchContentsByBlockId(block.Id)
		}
		for _, content := range block.BulletedListItem.Text {
			titles = append(titles, Title{content.PlainText, ContentArray})
		}
	}

	return titles
}

func FetchContentsByBlockId(BlockId string) []string {
	var contents []string
	var child_body Body
	var b = FetchNotion(BlockId)

	json.Unmarshal(b, &child_body)
	for _, block := range child_body.Results {
		for _, content := range block.BulletedListItem.Text {
			contents = append(contents, content.PlainText)
		}
	}
	return contents
}

func FetchNotion(BlockId string) []byte {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
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

func lineHandler(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	SecretToken := os.Getenv("LINE_CANNEL_SECRET_TOKEN")
	AccessToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
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
				var titles = FetchTitles()
				sendMessage := ShowTitles(titles)
				_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(sendMessage)).Do()
				if err != nil {
					log.Print(err)
				}
			}
		}
	}
}
