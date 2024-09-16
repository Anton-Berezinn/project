package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// MovieClient представляет клиента для обращения к API Kinopoisk
type MovieClient struct {
	BaseURL string
	Token   string
	Client  *http.Client
}

// NewMovieClient создает новый экземпляр клиента для работы с API
func NewMovieClient(token string) *MovieClient {
	return &MovieClient{
		BaseURL: "https://api.kinopoisk.dev/v1.4/movie",
		Token:   token,
		Client:  &http.Client{}, // инициализация стандартного HTTP-клиента
	}
}

// GetMovies отправляет запрос к API для получения фильмов по жанру и количеству
func (mc *MovieClient) GetMovies(count int, genre string) ([]byte, error) {
	// Создаем параметры запроса
	params := url.Values{}
	params.Add("limit", strconv.Itoa(count+50))
	params.Add("selectFields", "year")
	params.Add("selectFields", "genres")
	params.Add("genres.name", genre)
	params.Add("selectFields", "name")

	// Формируем URL с параметрами
	queryURL := fmt.Sprintf("%s?%s", mc.BaseURL, params.Encode())

	// Создаем HTTP-запрос
	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Добавляем токен в заголовки
	req.Header.Add("X-API-KEY", mc.Token)

	// Выполняем запрос
	resp, err := mc.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	// Чтение ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return body, nil
}
