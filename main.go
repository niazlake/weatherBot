package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	."weatherBot/models"
)


const botToken = "2106950597:AAE1mxNipcDVfK5k6s5gJg0hBnu1zLqKxqU"
const botApi = "https://api.telegram.org/bot"
const botUrl = botApi + botToken
var isStringAlphabetic = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

func main() {

	offset := 0
	for ; ; {
		updates, err := getUpdates(botUrl, offset)
		if err != nil {
			log.Println("Smth went wrong: ", err.Error())
		}

		for _, update :=range updates {
			if update.Message != nil {
				response(update)
				offset = update.UpdateId + 1
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



func response(update Update) error {
	var botMessage BotMessage
	var weatherInfo *WeatherResponse
	botMessage.ChatId = update.Message.Chat.ChatId

	if !isStringAlphabetic(update.Message.Text) {
		botMessage.Text = "Please enter in English"
		sendMessageBot(botMessage)
		return nil
	}

	weatherInfo, err := getWeather(update.Message.Text)
	if err != nil {
		return err
	}

	if weatherInfo == nil  || weatherInfo.Current == nil {
		return err
	}

	temperature := "Температура сейчас: " + FloatToString(weatherInfo.Current.TempC) + "C"
	city := "Город: " + weatherInfo.Location.Country + ", " + weatherInfo.Location.Name + "\n"
	botMessage.Text = city + temperature

	sendMessageBot(botMessage)

	return nil
}

func sendMessageBot(message BotMessage) error {
	buf, err := json.Marshal(message)
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