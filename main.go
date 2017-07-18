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
	"encoding/json"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client
type AQX struct {
SiteName string
County string
PSI int
MajorPollution string
Satus string
SO2 float32
CO float32
O3 int
PM10 int
PM2.5 int
NO2 float32
WindSpeed float32
WindDirection float32
FPMI int
NOx float32
NO float32
}

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)

	src_json := http://opendata2.epa.gov.tw/AQX.json
	u := AQX{}
	err  := json.Unmarshal(src_json, &u)
	if err != nil {
		panic(err)
	}
}

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
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(u["SiteName"]+"的 PM2.5 數值為 " + u["PM2.5"])).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}


