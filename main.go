package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	."weatherBot/models"
)

func main() {
	botToken := "2106950597:AAE1mxNipcDVfK5k6s5gJg0hBnu1zLqKxqU"
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken
	offset := 0
	for ; ; {
		updates, err := getUpdates(botUrl, offset)
		if err != nil {
			log.Println("Smth went wrong: ", err.Error())
		}

		for _, update :=range updates {
			response(botUrl, update)
			offset = update.UpdateId + 1
		}
	}
}

func getUpdates(botUrl string, offset int) ([]Update, error) {
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()


	body, err := ioutil.ReadAll(resp.Body)
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

func FloatToString(inputNumber float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(inputNumber, 'f', 0, 64)
}


func response(botUrl string, update Update) (error) {
	var botMessage BotMessage
	var weatherInfo *WeatherResponse

	weatherInfo, err := getWeather(update.Message.Text)
	if err != nil {
		return err
	}
	botMessage.ChatId = update.Message.Chat.ChatId
	botMessage.Text = "Температура сейчас: " + FloatToString(weatherInfo.Current.TempC) + "C"

	buf, err := json.Marshal(botMessage)
	if err != nil {
		return nil
	}

	_, err = http.Post(botUrl + "/sendMessage", "application/json", bytes.NewBuffer(buf))

	if err != nil {
		return nil
	}

	return nil
}


func getWeather(city string) (*WeatherResponse, error) {
	weatherToken := "1df1d168d065480c8be154020220302"
	url := "https://api.weatherapi.com/v1/current.json?key="
	var weatherResponse *WeatherResponse

	response, err := http.Get(url + weatherToken + "&q=" + city)

	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &weatherResponse)

	return weatherResponse, nil
}