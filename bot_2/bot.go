package main

import (
	"encoding/json"
	"flag"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"strings"
)

type Doc struct {
	Name string `json:"name"`
	Year int    `json:"year"`
}

type Response struct {
	Docs []Doc `json:"docs"`
}

const (
	FindGenre     = 0
	EnterQuantity = 1
)

var userChatAction = make(map[int64]int)      // Инициализация карты действий пользователей
var userMessageGenre = make(map[int64]string) // Сохранение жанра для каждого пользователя

func main() {
	tgToken, apiToken := getTokens()
	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {

		log.Fatal("error creating bot:", err)
	}

	client := NewMovieClient(apiToken)

	updates := requestAnswer(bot)
	for update := range updates {
		if update.Message == nil {
			continue // Пропускаем, если нет сообщения
		}

		chatID := update.Message.Chat.ID
		text := update.Message.Text

		if action, ok := userChatAction[chatID]; ok {
			switch action {
			case FindGenre:
				genre := strings.ToLower(text)
				if genreFound(genre) {
					userMessageGenre[chatID] = genre
					sendMessage(bot, chatID, "Сколько фильмов показывать?\nВведи число: ")
					userChatAction[chatID] = EnterQuantity
				} else {
					sendMessage(bot, chatID, "Такого жанра, я к сожалению не нашел\n   Введи еще раз:")
				}
				continue
			case EnterQuantity:
				countMovies, err := strconv.Atoi(text)
				if err != nil || countMovies <= 0 {
					sendMessage(bot, chatID, "Пожалуйста, введи корректное число больше 0")
				} else {
					genre := userMessageGenre[chatID]
					requestMovie(bot, client, genre, countMovies, chatID)
					delete(userChatAction, chatID) // Сбрасываем состояние после завершения
				}
				continue
			}
		}

		// Обработка команд вне состояний
		switch text {
		case "/start":
			sendMessage(bot, chatID, "**Привет**\n, я помогу тебе выбрать фильм или сериал😁Напиши жанр: ")
			userChatAction[chatID] = FindGenre
		default:
			sendMessage(bot, chatID, "Привет\n\tПопробуй написать команду /start")
		}
	}
}

func getTokens() (string, string) {
	tgToken := flag.String("tg-token", "", "Telegram bot token")
	apiToken := flag.String("api-token", "", "Api token")
	flag.Parse()
	return *tgToken, *apiToken
}

func requestAnswer(bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Printf("Не удалось сделать запрос к update %v", err)
	}
	return updates
}

func requestMovie(bot *tgbotapi.BotAPI, client *MovieClient, genre string, moviesCount int, chatID int64) {
	if moviesCount < 0 {
		sendMessage(bot, chatID, "Введи число, которое больше 0")
	}

	sendMessage(bot, chatID, "Я уже работаю над твоим запросом")
	answer, err := client.GetMovies(moviesCount, genre)
	if err != nil {
		log.Printf("Не удалось успешно сделать запрос %v", err)
		sendMessage(bot, chatID, "Не удалось успешно сделать запрос в API")
		return
	}
	responseUser(bot, chatID, answer, moviesCount)
}

func responseUser(bot *tgbotapi.BotAPI, chatID int64, msg []byte, count int) {
	var response Response
	if err := json.Unmarshal(msg, &response); err != nil {
		panic(err)
	}
	countMovie := 0
	for _, movie := range response.Docs {
		if count > 0 {
			if len(movie.Name) == 0 {
				continue
			}
			count--
			countMovie++
			sendMessage(bot, chatID, fmt.Sprintf("%v: Фильм:\n\t%v\n\tYear:\n\t%v", countMovie, movie.Name, movie.Year))
		} else {
			sendMessage(bot, chatID, "Будем ждать тебя снова\nНапиши /start")
			break
		}
	}
}

func genreFound(msg string) bool {
	genres := map[string]bool{
		"аниме":           true,
		"биография":       true,
		"боевик":          true,
		"вестерн":         true,
		"военный":         true,
		"детектив":        true,
		"детский":         true,
		"для взрослых":    true,
		"документальный":  true,
		"драма":           true,
		"игра":            true,
		"история":         true,
		"комедия":         true,
		"концерт":         true,
		"короткометражка": true,
		"криминал":        true,
		"мелодрама":       true,
		"музыка":          true,
		"мультфильм":      true,
		"мюзикл":          true,
		"новости":         true,
		"приключения":     true,
		"реальное ТВ":     true,
		"семейный":        true,
		"спорт":           true,
		"ток-шоу":         true,
		"триллер":         true,
		"ужасы":           true,
		"фантастика":      true,
		"фильм-нуар":      true,
		"фэнтези":         true,
		"церемония":       true,
	}
	return genres[msg]
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("error in sendMessage %v", err)
	}
}
