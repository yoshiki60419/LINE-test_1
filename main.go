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
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/line/line-bot-sdk-go/linebot"
)

type exrate struct {
	inCashRate, outCashRate, inRate, outRate string
}

var exRates map[string]exrate
var bot *linebot.Client

func main() {
	// crawler
	exRates = make(map[string]exrate, 0)

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

	doc, err := goquery.NewDocument("http://rate.bot.com.tw/Pages/Static/UIP003.zh-TW.htm")
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("[class=\"titleLeft\"]").Each(func(i int, s *goquery.Selection) {
		currency := strings.TrimSpace(s.Text())
		pos := strings.Index(currency, " ")
		currCut := currency[0:pos]
		fmt.Printf("%s: ", currCut)
		//		fmt.Printf("%d", len(currCut))
		//		if currCut == "美金" {
		//			fmt.Println("美金")
		//		}
		inCashRate := s.Next().Text()
		outCashRate := s.Next().Next().Text()
		inRate := s.Next().Next().Next().Text()
		outRate := s.Next().Next().Next().Next().Text()
		var rate exrate
		rate.inCashRate = inCashRate
		rate.outCashRate = outCashRate
		rate.inRate = inRate
		rate.outRate = outRate
		exRates[currCut] = rate
		// fmt.Printf("%s %s %s %s\n", inCashRate, outCashRate, inRate, outRate)
		// fmt.Printf("%s %s %s %s\n", exRates["美金"].inCashRate, exRates["美金"].outCashRate, exRates["美金"].inRate, exRates["美金"].outRate)
	})

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
				if message.Text == "USD" {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("現在美金匯率為: "+exRates["美金"].inCashRate)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	}
}
