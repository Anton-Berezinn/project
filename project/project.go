package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strconv"
	"testing"
)

type Root struct {
	XMLName xml.Name `xml:"root"`
	Rows    []Row    `xml:"row"`
}

type Row struct {
	Id        string `xml:"id"`
	Age       string `xml:"age"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	About     string `xml:"about"`
}

type Data struct {
	query       string
	order_field string
	order_by    string
	limit       int
	offset      int
}

type Answer struct {
	Id    string
	Age   string
	Name  string
	About string
}

func (r *Data) GetUser() ([]Answer, error) {
	f, err := os.Open("dataset.xml")
	if err != nil {
		log.Println("error in open file", err)
	}
	defer f.Close()
	var root Root
	decoder := xml.NewDecoder(f)
	err = decoder.Decode(&root)
	if err != nil {
		return nil, err
	}
	var answer []Answer
	if r.query == "" {
		for _, row := range root.Rows {
			answer = append(answer, Answer{
				row.Id, row.Age, row.FirstName + row.LastName, row.About,
			})
		}
		return answer, nil
	}
	if r.order_by == "asc" {
		switch r.order_field {
		case "id":
			for _, row := range root.Rows[r.offset:r.limit] {
				answer = append(answer, Answer{
					row.Id, row.Age, row.FirstName + row.LastName, row.About,
				})
			}
			sort.Slice(answer, func(i, j int) bool {
				return answer[i].Id > answer[j].Id
			})
		case "age":
			for _, row := range root.Rows[r.offset:r.limit] {
				answer = append(answer, Answer{
					row.Id, row.Age, row.FirstName + row.LastName, row.About,
				})
			}
			sort.Slice(answer, func(i, j int) bool {
				return answer[i].Age > answer[j].Age
			})
		case "name":
			for _, row := range root.Rows[r.offset:r.limit] {
				answer = append(answer, Answer{
					row.Id, row.Age, row.FirstName + row.LastName, row.About,
				})
			}
			sort.Slice(answer, func(i, j int) bool {
				return answer[i].Name > answer[j].Name
			})
		}
	} else {
		switch r.order_field {
		case "id":
			for _, row := range root.Rows[r.offset:r.limit] {
				answer = append(answer, Answer{
					row.Id, row.Age, row.FirstName + row.LastName, row.About,
				})
			}
			sort.Slice(answer, func(i, j int) bool {
				return answer[i].Id < answer[j].Id
			})
		case "age":
			for _, row := range root.Rows[r.offset:r.limit] {
				answer = append(answer, Answer{
					row.Id, row.Age, row.FirstName + row.LastName, row.About,
				})
			}
			sort.Slice(answer, func(i, j int) bool {
				return answer[i].Age < answer[j].Age
			})
		case "name":
			for _, row := range root.Rows[r.offset:r.limit] {
				answer = append(answer, Answer{
					row.Id, row.Age, row.FirstName + row.LastName, row.About,
				})
			}
			sort.Slice(answer, func(i, j int) bool {
				return answer[i].Name < answer[j].Name
			})
		}
	}
	return answer, nil
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	switch query {
	case "name":
		fmt.Fprintf(w, "Hello World")
	case "about":
		fmt.Fprintf(w, "About")
	}
	order_field := r.URL.Query().Get("order_field")
	switch order_field {
	case "id":
		fmt.Fprintf(w, "Id")
	case "age":
		fmt.Fprintf(w, "Age")
	case "name":
		fmt.Fprintf(w, "Name")
	case "":
		fmt.Fprintf(w, "Name")
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
	order_by := r.URL.Query().Get("order_by")
	switch order_by {
	case "asc":
		fmt.Fprintf(w, "Asc")
	case "desc":
		fmt.Fprintf(w, "Desc")
	case "":
		fmt.Fprintf(w, "You need to point order_by...")
	}
	limit := r.URL.Query().Get("limit")
	v, err := strconv.Atoi(limit)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	offset := r.URL.Query().Get("offset")
	offs, err := strconv.Atoi(offset)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	data := &Data{
		query:       query,
		order_field: order_field,
		order_by:    order_by,
		limit:       v,
		offset:      offs,
	}
	dat, err := data.GetUser()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Println(dat)

}

func TestHandle(t *testing.T) {
	//test:= []Data {
	//	Data {
	//		"name","age","asc",1,3,
	//	},
	//}
	server := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer server.Close()
	resp, err := http.Get(server.URL + "?query=name&order_field=age&order_by=desc&limit=1&offset=0")
	if err != nil {
		t.Error(err.Error())
	}
	defer resp.Body.Close()
}
