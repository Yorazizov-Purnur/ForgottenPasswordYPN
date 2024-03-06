package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"

	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type SignUpStruct struct {
	Name          string
	TelegramLogin string
	Password      string
}

var SignUpSlice = []SignUpStruct{} // ? empty
func main() {
	r := gin.Default()

	r.Use(Cors)
	r.POST("/signup", SignUp)
	go Recovering()

	r.Run(":3434")
}

func Recovering() {
	Readuser()
	BotResult, BotError := tgbotapi.NewBotAPI("7195811917:AAF_HGF45tavsok0yZoERSAczrYZpeG8gTk")

	if BotError != nil {
		fmt.Printf("BotError: %v\n", BotError)
	}

	updates := tgbotapi.NewUpdate(0)

	RetriveResult, err := BotResult.GetUpdatesChan(updates)

	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	fmt.Println("Connection done")

	IsEditingPassword := false

	for Item := range RetriveResult {
		if Item.Message.IsCommand() {
			if Item.Message.Command() == "reset" {

				for _, Item2 := range SignUpSlice {
					if Item2.TelegramLogin == Item.Message.Chat.UserName {
						msg := tgbotapi.NewMessage(Item.Message.Chat.ID, "Enter new password")
						BotResult.Send(msg)
					}
				}
			}
		}else{
			if IsEditingPassword {
				for i, Item2 := range SignUpSlice {
					if Item2.TelegramLogin == Item.Message.Chat.UserName {
						SignUpSlice[i].Password = Item.Message.Text
					}
				}
			}
		}
	}
}
func SignUp(c *gin.Context) {
	var SignUpTemp SignUpStruct
	c.ShouldBindJSON(&SignUpTemp)

	if SignUpTemp.Name == "" || SignUpTemp.Password == "" || SignUpTemp.TelegramLogin == "" {
		c.JSON(404, "Empty field")
	} else {
		Readuser()
		SignUpSlice = append(SignUpSlice, SignUpTemp)
		Writeuser()
	}
}

func Writeuser() {
	mardhalData, _ := json.Marshal(SignUpSlice)
	ioutil.WriteFile("app.json", mardhalData, 0644)
}

func Readuser() {
	readbyte, _ := ioutil.ReadFile("app.json")
	json.Unmarshal(readbyte, &SignUpSlice)
}

func Cors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://192.168.43.246:5500")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
	}

	c.Next()
}
