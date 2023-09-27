package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ornen/go-sbs1"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

type air struct {
	R string `json:"r"`
	T string `json:"t"`
	F string `json:"f"`
	D string `json:"d"`
}
type airName struct {
	D string `json:"d"`
}
type airType struct {
	T string `json:"t"`
}
type airTag struct {
	R string `json:"r"`
}

type airDesc struct {
	Desc string `json:"desc"`
}

func main() {

	aircraftsjson, err := os.ReadFile("./aircrafts.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	aircraftsjson2, err := os.ReadFile("./aircrafts.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	typesjson, err := os.ReadFile("./types.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	var planes map[string]air
	var planeName map[string]airName
	var planeType map[string]airType
	var planeDesc map[string]airDesc
	var planeTag map[string]airTag
	var planeTypeStr string
	var planeDescStr string
	var planeTagStr string
	var HexIdStr string
	var lastplanes [10]string
	var planenum int
	planenum = 0
	var planewashere bool
	var checkplane bool
	var tenmintimer bool
	var ispolice bool
	var planetimeAdded [10]time.Time
	var planetimeAddedTimer [10]time.Time
	var planetimeAddedTimerHour [10]time.Time
	
	for i := 0; i <= 9; i++ {
		planetimeAdded[i] = time.Now().Add(-61 * time.Minute)
		planetimeAddedTimer[i] = time.Now().Add(-61 * time.Minute)
		planetimeAddedTimerHour[i] = time.Now().Add(-61 * time.Minute)
	}
	var apicheck time.Time
	apicheck = time.Now().Add(-80 * time.Second)
	var counter int
	var eofcounter int
	var lasterror time.Time
	lasterror = time.Now()

	err = json.Unmarshal([]byte(aircraftsjson), &planes)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal([]byte(aircraftsjson), &planeName)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal([]byte(aircraftsjson), &planeType)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal([]byte(aircraftsjson2), &planeTag)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal([]byte(typesjson), &planeDesc)
	if err != nil {
		fmt.Println(err)
	}
	bot, err := tgbotapi.NewBotAPI("<bot-key>")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	msg := tgbotapi.NewMessage(int64(<channel-id>), "")
	msg.Text = "startup successful"
	bot.Send(msg)
	for true {

		conn, err := net.Dial("tcp", "10.0.0.20:30003")
		if err != nil {
			msg := tgbotapi.NewMessage(int64(<channel-id>), "")
			msg.Text = "ERROR \"critical error\" shutting down"
			bot.Send(msg)
			log.Fatal(err)
		}

		//defer conn.Close()

		var reader = sbs1.NewReader(conn)
		apicheck = time.Now()
		for {
			var message, err = reader.Read()

			if err != nil {
				if err == io.EOF {
					if time.Now().After(lasterror.Add(5 * time.Second)) {
						eofcounter = 0
					}
					if eofcounter >= 10 {
						msg := tgbotapi.NewMessage(int64(<channel-id>), "")
						msg.Text = "ERROR \"io.EOF\" API might not be working properly"
						bot.Send(msg)
						eofcounter = 0
						lasterror = time.Now()
					}
					eofcounter++
					break
				} else {
					log.Println(err)
					continue
				}
			}
			planewashere = false
			checkplane = false
			HexIdStr = fmt.Sprint(message.HexId)
			planeTypeStr = fmt.Sprint(planeType[message.HexId])
			planeTypeStr = strings.Trim(planeTypeStr, "{}")
			planeDescStr = fmt.Sprint(planeDesc[planeTypeStr])
			planeDescStr = strings.Trim(planeDescStr, "{}")
			isheli, _ := regexp.MatchString(`^H`, planeDescStr)

			planeTagStr = fmt.Sprint(planeTag[message.HexId])
			planeTagStr = strings.Trim(planeTagStr, "{}")
			ispolice, _ = regexp.MatchString(`OE-B`, planeTagStr)
			if time.Now().After(apicheck.Add(300 * time.Second)) {
				msg := tgbotapi.NewMessage(int64(<channel-id>), "")
				msg.Text = "ERROR \"no updates\" API might not be working properly or no aircafts are getting detected in the airspace"
				bot.Send(msg)
				if counter >= 5 {
					time.Sleep(1 * time.Minute)
					counter = -1
				}
				counter++
			}
			apicheck = time.Now()

			if isheli && ispolice {
				tenmintimer = false

				for i := 0; i <= 9; i++ {

					if time.Now().After(planetimeAdded[i].Add(10 * time.Minute)) {
						planetimeAdded[i] = time.Now().Add(-61 * time.Minute)
						lastplanes[i] = "0"
					}
					checkplane, _ = regexp.MatchString(HexIdStr, lastplanes[i])
					if checkplane == true {
						planewashere = true
						if time.Now().After(planetimeAdded[i].Add(25 * time.Minute)) {
							lastplanes[i] = "0"
							planetimeAdded[i] = time.Now().Add(-61 * time.Minute)
							tenmintimer = false
							planewashere = false
						}
						if planewashere == true && time.Now().After(planetimeAdded[i].Add(7*time.Minute)) {
							lastplanes[i] = "0"
							planetimeAdded[i] = time.Now().Add(-61 * time.Minute)
							tenmintimer = true
							planewashere = false
						}

					}
				}

				if planewashere == false {
					lastplanes[planenum] = HexIdStr
					planetimeAdded[planenum] = time.Now()

					planenum++
					if planenum > 9 {
						planenum = 0
					}

					if tenmintimer == true {
						//send telegram still in the air
						fmt.Print(HexIdStr)
						fmt.Println(" is still in the air")
						msg := tgbotapi.NewMessage(int64(<channel-id>), "")
						msg.ParseMode = "html"
						msg.DisableWebPagePreview = true
						msg.Text = fmt.Sprintf("%s is still in the air\n<a href=\"https://globe.adsbexchange.com/?icao=%s\">track flight</a>", planeTagStr, HexIdStr)
						bot.Send(msg)
					}
					if tenmintimer == false {
						//send telegram just enetered airspace
						fmt.Print(HexIdStr)
						fmt.Println(" has just entered airspace")
						msg := tgbotapi.NewMessage(int64(<channel-id>), "")
						msg.ParseMode = "html"
						msg.DisableWebPagePreview = true
						msg.Text = fmt.Sprintf("ALERT\nPOLICE Helicopter\n%s has just entered airspace\n<a href=\"https://globe.adsbexchange.com/?icao=%s\">track flight</a>", planeTagStr, HexIdStr)
						bot.Send(msg)
					}
					fmt.Println(HexIdStr)
					fmt.Print(planeTypeStr)
					fmt.Println(" is new Heli")
				}
			}
		}
	}
}
