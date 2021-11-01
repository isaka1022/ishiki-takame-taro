package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Content struct {
	Content string
	Link    string
}

type Annotations struct {
	Bold          bool
	Italic        bool
	StrikeThrough bool
	Underline     bool
	Code          bool
	Color         string
}

type BulletedListItemContent struct {
	Type        string
	Text        Content
	Annotations Annotations
	PlainText   string
	href        string
}

type BulletedListItem struct {
	Text []*BulletedListItemContent
}

type Block struct {
	Object         string
	Id             string
	CreatedTime    string
	LastEditedTime string
	hasChildren    bool
	Archived       bool
	Type           string
	BulletedListItem
}

type Blocks []*Block

type Body struct {
	Object     string
	Results    Blocks
	NextCursor string
	hasMore    bool
}

func main() {
	err := godotenv.Load(".env")
	ApiKey := os.Getenv("NOTION_SECRET_KEY")
	DatabaseId := os.Getenv("NOTION_DATABASE_ID")

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.notion.com/v1/blocks/"+DatabaseId+"/children?page_size=5", nil)
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

	var body Body
	json.Unmarshal(b, &body)
	fmt.Println(body)
	fmt.Println(body.Results)

	for _, block := range body.Results {
		fmt.Println(block.Object)
		fmt.Println(block.Type)
		fmt.Println(block.BulletedListItem)
		for _, item := range block.BulletedListItem.Text {
			fmt.Println(item.Type)
		}
		// if block.Type == "heading_1" {
		// 	fmt.Println(block.BulletedListItems
		// 	)
		// }
	}

	var out bytes.Buffer
	json.Indent(&out, b, "", " ")
	out.WriteTo(os.Stdout)

}

// func main() {
// 	http.HandleFunc("/", helloHandler)
// 	http.HandleFunc("/callback", lineHandler)

// 	fmt.Println("https://localhost:8080 で起動中...")

// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

// func helloHandler(w http.ResponseWriter, r *http.Request) {
// 	msg := "Hello world!!!!"
// 	fmt.Fprintf(w, msg)
// }

// func lineHandler(w http.ResponseWriter, r *http.Request) {
// 	bot, err := linebot.New(
// 		"8d269eda5acfc678f48a4e9d0ea0fd55",
// 		"KHWIbLhXzgRLS6zqDclV6jFBPZODnu8jIYBmvGA3bWQFug6v0xdhgvVy+1ujPrN0rEWnrnbfX1bDpz4C7FmbWYBB19ne2Mm/F5uIL5WmXXFsGWXZa2qEp8yjP6SnKcTBszzK2SlTIqDPKpV9hqWtawdB04t89/1O/w1cDnyilFU=",
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	events, err := bot.ParseRequest(r)
// 	if err != nil {
// 		if err == linebot.ErrInvalidSignature {
// 			w.WriteHeader(400)
// 		} else {
// 			w.WriteHeader(500)
// 		}
// 	}

// 	for _, event := range events {
// 		if event.Type == linebot.EventTypeMessage {
// 			switch message := event.Message.(type) {
// 			case *linebot.TextMessage:
// 				replyMessage := message.Text
// 				_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()
// 				if err != nil {
// 					log.Print(err)
// 				}

// 			case *linebot.LocationMessage:
// 				sendRestoInfo(bot, event)
// 			}
// 		}
// 	}
// }

// func sendRestoInfo(bot *linebot.Client, e *linebot.Event) {
// 	msg := e.Message.(*linebot.LocationMessage)

// 	let := strconv.FormatFloat(msg.Latitude, 'f', 2, 64)
// 	lng := strconv.FormatFloat(msg.Longitude, 'f', 2, 64)

// 	replyMsg := fmt.Sprintf("緯度：%s\n経度： %s", let, lng)

// 	_, err := bot.ReplyMessage(e.ReplyToken, linebot.NewTextMessage(replyMsg)).Do()
// 	if err != nil {
// 		log.Print(err)
// 	}
