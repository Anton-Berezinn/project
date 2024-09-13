package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func request_movie(c int, genre string) []byte {
	token := "E5RKRGN-RS3MSK4-GRPKWFA-QDZ0VH3"
	params := url.Values{} // Добавляем параметры к url
	params.Add("limit", strconv.Itoa(c+50))
	params.Add("selectFields", "year")
	params.Add("selectFields", "genres")
	params.Add("genres.name", genre)
	params.Add("selectFields", "genres")
	params.Add("selectFields", "name")
	baseURL := "https://api.kinopoisk.dev/v1.4/movie"
	queryURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	req, err := http.NewRequest("GET", queryURL, nil) // Составляем запрос
	if err != nil {
		panic(err)
	}
	req.Header.Add("X-API-KEY", token)
	client := http.Client{}     // Создаем клиента чтобы сделать запрос и получить ответ
	resp, err := client.Do(req) // делаем запрос и получаем ответ
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return body
}
