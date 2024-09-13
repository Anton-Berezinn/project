package main

import (
	"encoding/json"
	"flag"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	User_message_genre string
)

type Doc struct {
	Name string `json:"name"`
	Year int    `json:"year"`
}

type Response struct {
	Docs []Doc `json:"docs"`
}

func main() {
	bot, err := tgbotapi.NewBotAPI(tok())
	if err != nil {
		log.Fatal("error creating bot:", err)
	}
	updates := request_answer(bot)
	for {
		for update := range updates {
			chatID := update.Message.Chat.ID
			switch update.Message.Text {
			case "/start":
				sendmessage(bot, chatID, "**Привет**\n, я помогу тебе выбрать фильм или сериал😁Напиши жанр: ")
				start(bot, chatID)
			default:
				sendmessage(bot, chatID, "Привет\n\tПопробуй написать команду /start")
			}

		}
	}
}

func tok() string {
	token := flag.String("token", "", "Telegram bot token")
	flag.Parse()
	return *token
}
func request_answer(bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	time.Sleep(5 * time.Second)
	u := tgbotapi.NewUpdate(0)

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}
	return updates
}

func start(bot *tgbotapi.BotAPI, chatID int64) {
	updates := request_answer(bot)
	for {
		for update := range updates {
			if update.Message.Chat.ID == chatID {
				User_message_genre = strings.ToLower(update.Message.Text)
				if FindMovie(User_message_genre) {
					Req_movie(bot, chatID)
				} else {
					sendmessage(bot, chatID, "Такого жанра, я к сожалению не нашел\n   Введи еще раз:")
				}
			} else {
				sendmessage(bot, chatID, "Упс, something wrong.... ")
			}
		}

	}
}

func Req_movie(bot *tgbotapi.BotAPI, chatID int64) {
	sendmessage(bot, chatID, "Сколько фильмов показывать?\nВведи число: ")
	updates := request_answer(bot)
	for {
		for update := range updates {
			if update.Message.Chat.ID == chatID {
				count_movies := update.Message.Text
				if C, err := strconv.Atoi(count_movies); err == nil {
					if C > 0 {
						sendmessage(bot, chatID, "Я уже работаю над твоим запросом")
						answer := request_movie(C, User_message_genre)
						Response_user(bot, chatID, answer, C)
					} else {
						sendmessage(bot, chatID, "Введи число которое больше 0")
					}
				} else {
					m := "Убедись что ты правильно написал число"
					sendmessage(bot, chatID, m)
				}

			} else {
				sendmessage(bot, chatID, "something wrong")
			}
		}
	}

}

func Response_user(bot *tgbotapi.BotAPI, chatID int64, msg []byte, count int) {
	updates := request_answer(bot)
	var response Response
	if err := json.Unmarshal(msg, &response); err != nil {
		panic(err)
	}
	for update := range updates {
		if update.Message.Chat.ID == chatID {
			count_movie := 0
			for _, movie := range response.Docs {
				if count > 0 {
					fmt.Println(len(movie.Name))
					if len(movie.Name) == 0 {
						continue
					} else {
						fmt.Println("dsfe")
						count--
						count_movie += 1
						sendmessage(bot, chatID, fmt.Sprintf("%v: Фильм:\n\t%v\n\tYear:\n\t%v", count_movie, movie.Name, movie.Year))
					}
				} else {
					sendmessage(bot, chatID, "Будем ждать тебя снова")
					break
				}
			}
		}
	}

}

func FindMovie(msg string) bool {
	genre := map[string]bool{
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
	return genre[msg]
}

func sendmessage(bot *tgbotapi.BotAPI, chatId int64, message string) {
	m := tgbotapi.NewMessage(chatId, message)
	if _, err := bot.Send(m); err != nil {
		log.Fatal(err)
	}
}
