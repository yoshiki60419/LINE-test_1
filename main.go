// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

type exrate struct {
	inCashRate, outCashRate, inRate, outRate string
}

var exRates map[string]exrate
var bot *linebot.Client

func main() {
	// Line bot
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

// Line bot
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if message.Text == "用戶進場" {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("您的進場資訊：\n車牌：AB-1234(一般用戶) 停於：南港停車場。\n進場時間：2017/7/25 10:05。\n出場代碼：5595。\n提醒您，若出場時車辨失敗，請輸入出場代碼即可靠卡繳費，謝謝您的使用。")).Do(); err != nil {
						log.Print(err)
					}
				} else if message.Text == "用戶出場" {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("您的出場資訊：\n車牌：4Q-5678(一般用戶) 停於：南港停車場。\n出場時間：2017/7/25 12:05；出場代碼：5595。\n本次停車時間：2 小時 0 分。\n本次停車費用：60 元\n本次優惠折抵：0 元\n祝您路途順利，謝謝您的使用。")).Do(); err != nil {
						log.Print(err)
					}
				} else if message.Text == "用戶資訊" {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("您的用戶資訊：\n車牌：4Q-5678(一般用戶)。\n您的付費方式：信用卡(詳細情形請查詢\"信用卡設定\")。\n本月停車時間：8 小時 30　分。\n本月停車費用：270 元。\n本月停車次數：4 次。\n祝您路途順利，謝謝您的使用。")).Do(); err != nil {
						log.Print(err)
					}
				} else if message.Text == "停車資訊" {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("您的停車資訊：\n車牌：3Q-1245(月租用戶) 停於：南港停車場。\n進場時間：2017/7/25 10:05。\n出場時間：2017/7/25 12:05。\n出場代碼：5595。\n提醒您，若出場時車辨失敗，請輸入出場代碼即可靠卡繳費，謝謝您的使用。")).Do(); err != nil {
						log.Print(err)
					}
				} else {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text+" 請輸入停車相關服務，謝謝您的合作。")).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	}
}
