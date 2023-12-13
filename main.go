package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)



func getVideos() (int, error) {
	url := "https://www.tiktok.com/@dazangela"
	response, err := http.Get(url)
	if err != nil {
		return -1, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return -1, err
	}
	fbody := string(body)
	regex := `"videoCount"\s*:\s*(\d+)`
	r := regexp.MustCompile(regex)
	match := r.FindStringSubmatch(fbody)
	if len(match) >= 2 {
		videoCount := match[1]
		count := 0
		fmt.Sscan(videoCount, &count)
		return count, nil
	}
	return -1, fmt.Errorf("Não foi possível encontrar o número de vídeos")
}

func monitorVideoCount(ch chan<- int) {
	lastVideoCount := -1
	for {
		videoCount, err := getVideos()
		if err != nil {
			log.Println("Erro ao obter vídeos:", err)
			time.Sleep(30 * time.Second) // Tentar novamente após 30 segundos em caso de erro
			continue
		}

		if videoCount != -1 && videoCount != lastVideoCount {
			if lastVideoCount != -1 {
				ch <- videoCount
			}
			lastVideoCount = videoCount
		}
		time.Sleep(20 * time.Second) // Verificar a cada 5 minutos
	}
}

func main() {
	botToken := "6820528824:AAHx4HPkuOdZo2x0ACoBKg2qvZrCfaLfTM8" // Substitua pelo token do seu bot

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	videoChangeChannel := make(chan int)
	go monitorVideoCount(videoChangeChannel)

	for newCount := range videoChangeChannel {
		msg := tgbotapi.NewMessage(5952155531, fmt.Sprintf("Angela postou um vídeo novo! Total de vídeos: %d", newCount))
		_, err := bot.Send(msg)
		if err != nil {
			log.Println("Erro ao enviar mensagem:", err)
		}
	}
}
