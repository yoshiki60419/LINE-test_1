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
				if message.Text == "早安" {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text+" 好哥有ㄛ飽嗎~")).Do(); err != nil {
						log.Print(err)
					}
				}
				if message.Text == "好哥智商多少" {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text+" 87 不能再高~~")).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	}
}
