package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	addrConv   = "https://api.apilayer.com/exchangerates_data"
	apiKeyConv = "rTCs2Rp0bYfFMU9BgTQqdhezF6myAg5t"
)

func main() {
	botToken := "5990663565:AAFaPCQHVX8z93H1NG4CvH6FZuRa_JBGoZw"
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken
	offset := 0
	currencyConverter := NewCurrencyConverter(addrConv, apiKeyConv)
	currencies, err := currencyConverter.GetCurrencies()
	if err != nil {
		log.Fatal(err)
	}

	for {
		updates, err := getUpdates(botUrl, offset)
		if err != nil {
			log.Println("Something went wrong: ", err)
		}

		for _, update := range updates {
			offset = update.UpdateID + 1

			curs := strings.Split(update.Message.Text, " ")

			if len(curs) != 2 {
				err = respond(botUrl, update, "Передайте две валюты через пробел, например: USD RUB")
				if err != nil {
					log.Println(err)
				}
				continue
			}

			curs[0] = strings.ToUpper(curs[0])
			curs[1] = strings.ToUpper(curs[1])

			_, ok1 := currencies[curs[0]]
			_, ok2 := currencies[curs[1]]

			if !ok1 || !ok2 {
				err = respond(
					botUrl,
					update,
					"Неизвестные валюты. Передайте две валюты через пробел, например: USD RUB",
				)
				if err != nil {
					log.Println(err)
				}
				continue
			}

			result, err := currencyConverter.Convert(1, curs[0], curs[1])
			if err != nil {
				log.Println(err)
				continue
			}

			text := fmt.Sprintf("1 %s = %f %s", curs[0], result, curs[1])
			err = respond(botUrl, update, text)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

func getUpdates(botUrl string, offset int) ([]Update, error) {
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}

func respond(botUrl string, update Update, text string) error {
	var botMessage BotMessage
	botMessage.ChatID = update.Message.Chat.ChatID
	botMessage.Text = text
	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil

}
