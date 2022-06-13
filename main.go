package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type Verify struct {
	Mode      string `form:"hub.mode"`
	Token     string `form:"hub.verify_token"`
	Challenge int    `form:"hub.challenge"`
}

type Body struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type Entry struct {
	Id      string   `json:"id"`
	Changes []Change `json:"changes"`
}

type Value struct {
	MessagingProduct string    `json:"messaging_product"`
	Metadata         Metadata  `json:"metadata"`
	Contacts         []Contact `json:"contacts"`
	Messages         []Message `json:"messages"`
}

type Change struct {
	Field string `json:"field"`
	Value Value  `json:"value"`
}

type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberId      string `json:"phone_number_id"`
}

type Contact struct {
	WaId    string  `json:"wa_id"`
	Profile Profile `json:"profile"`
}

type Profile struct {
	Name string `json:"name"`
}

type Message struct {
	From      string `json:"from"`
	Id        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
	Text      Text   `json:"text"`
}

type Text struct {
	Body string `json:"body"`
}

type PostBody struct {
	MessagingProduct string `json:"messaging_product"`
	To               string `json:"to"`
	Text             Text   `json:"text"`
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	db, err := sql.Open("sqlite3", "./reminders.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlCreateStmt := `
	CREATE TABLE IF NOT EXISTS reminders (
	id INTEGER NOT NULL PRIMARY KEY,
	date_created DATETIME NOT NULL,
	user_number TEXT,
	user_name TEXT,
	reminder TEXT,
	reminder_date DATETIME,
	been_reminded BOOL
	); 
	`
	_, err = db.Exec(sqlCreateStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlCreateStmt)
		return
	}

	port := goDotEnvVariable("PORT")

	if port == "" {
		port = "8000"
	}
	router := gin.Default()
	router.GET("/webhook", func(c *gin.Context) {
		var verify Verify
		verifyToken := os.Getenv("VERIFY_TOKEN")
		if c.ShouldBindQuery(&verify) == nil {
			if verify.Mode != "" && verify.Token != "" {
				if verify.Mode == "subscribe" && verify.Token == verifyToken {
					fmt.Println("WEBHOOK VERIFIED")
					c.JSON(http.StatusOK, verify.Challenge)
				} else {
					c.JSON(http.StatusForbidden, "Forbidden")
				}
			}
		}
	})

	router.POST("/webhook", func(c *gin.Context) {
		var body Body
		if c.ShouldBindJSON(&body) == nil {
			if body.Object != "" {
				if body.Entry != nil && body.Entry[0].Changes != nil && body.Entry[0].Changes[0].Value.Messages != nil && (body.Entry[0].Changes[0].Value.Messages[0] != Message{}) {
					phoneNumberId := body.Entry[0].Changes[0].Value.Metadata.PhoneNumberId
					name := body.Entry[0].Changes[0].Value.Contacts[0].Profile.Name
					from := body.Entry[0].Changes[0].Value.Messages[0].From
					messBody := body.Entry[0].Changes[0].Value.Messages[0].Text.Body
					fmt.Println(phoneNumberId, name, from, messBody)

					_, err = db.Exec("INSERT INTO reminders VALUES(NULL,?,?,?,?,?,?);", time.Now(), from, name, messBody, time.Now(), false)

					if err != nil {
						log.Printf("%q: %s\n", err, sqlCreateStmt)
					}

					token := goDotEnvVariable("FACEBOOK_TOKEN")

					postBody, _ := json.Marshal(PostBody{
						MessagingProduct: "whatsapp",
						To:               from,
						Text: Text{
							Body: "I will remind you at the end of the week.",
						},
					})

					responseBody := bytes.NewBuffer(postBody)

					resp, err := http.Post(fmt.Sprintf("https://graph.facebook.com/v12.0/%s/messages?access_token=%s", phoneNumberId, token), "application/json", responseBody)

					if err != nil {
						log.Fatalf("An error occured %v", err)
					}

					defer resp.Body.Close()

					c.JSON(http.StatusOK, "Sent")
				}
			} else {
				c.JSON(http.StatusBadGateway, "Not found")
			}
		}
	})

	router.Run(fmt.Sprintf(":%s", port))
}
