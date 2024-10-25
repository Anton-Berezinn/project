package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	switch key {
	case "1":
		io.WriteString(w, `{"id":"1","name":"aa","last_name":"bb"}`)
	case "2":
		io.WriteString(w, `{"id":"2","name":"aaa","last_name":"bbb"}`)
	default:
		http.Error(w, "Key not found", http.StatusNotFound)
	}
}

type Test struct {
	url string
}

type User struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Last_name string `json:"last_name"`
}

func (t *Test) Users(id string) (*User, error) {
	url := t.url + "?key=" + id
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data User
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func TestHandle(t *testing.T) {
	cases := []struct {
		id       string
		wantName string
		wantLast string
	}{
		{"1", "aa", "bb"},
		{"2", "aaa", "bbb"},
	}
	server := httptest.NewServer(http.HandlerFunc(Handle))
	for _, val := range cases {
		a := &Test{
			url: server.URL,
		}
		an, err := a.Users(val.id)
		if err != nil {
			if val.id == "4" {
				// ожидаем ошибку для несуществующего ключа
				continue
			}
			t.Error(err.Error())
		}
		if an.Name != val.wantName {
			t.Errorf("want name %q, got %q", val.wantName, an.Name)
		}
		if an.Last_name != val.wantLast {
			t.Errorf("want last name %q, got %q", val.wantLast, an.Last_name)
		}
	}
}
