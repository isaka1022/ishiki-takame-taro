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
	hasMore    bool   `json:"has_more"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	DatabaseId := os.Getenv("NOTION_DATABASE_ID")

	var texts []string

	var body Body
	var b = Request(DatabaseId)
	json.Unmarshal(b, &body)

	for _, block := range body.Results {
		fmt.Println(block.Id)
		fmt.Println(block.HasChildren)
		if block.HasChildren == true {

		}
		for _, content := range block.BulletedListItem.Text {
			texts = append(texts, content.PlainText)
		}
	}

	// fmt.Println(texts)
	// rand.Seed(time.Now().UnixNano())
	// num := rand.Intn(len(texts))
	// fmt.Println(texts[num])

	var out bytes.Buffer
	json.Indent(&out, b, "", " ")
	out.WriteTo(os.Stdout)
}

func Request(BlockId string) []byte {
	ApiKey := os.Getenv("NOTION_SECRET_KEY")

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.notion.com/v1/blocks/"+BlockId+"/children?page_size=5", nil)
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
