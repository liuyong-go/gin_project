package main

import (
	"context"
	"fmt"

	"github.com/liuyong-go/gin_project/bootstrap"
	"github.com/liuyong-go/gin_project/libs/yelastic"
	"github.com/spf13/cast"
)

var ctx = context.Background()

func main() {

	bootstrap.TestInit()
	//createDoc()
	//updateDoc()
	getByIDs()
	//getByID()
	//search()
	//	wordSearch()
	//multiSearch()
	//delete()
}

type employee struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Age       int      `json:"age"`
	About     string   `json:"about"`
	Interests []string `json:"interests"`
}

func createDoc() {
	var docs = []employee{
		{
			FirstName: "John",
			LastName:  "Smith",
			Age:       25,
			About:     "I love to go rock climbing 1",
			Interests: []string{"sports", "music"},
		}, {
			FirstName: "Jane",
			LastName:  "Smith",
			Age:       32,
			About:     "I like to collect rock albums 1",
			Interests: []string{"music"},
		},
		{
			FirstName: "Douglas",
			LastName:  "Fir",
			Age:       35,
			About:     "I like to build cabinets 1",
			Interests: []string{"forestry"},
		},
	}

	for index, doc := range docs {
		body, err := yelastic.NewES().CreateDocument(ctx, "megacorp", "employee", cast.ToString(index+1), doc)
		fmt.Println(string(body), err)
	}
}
func updateDoc() {
	var docs = []employee{
		{
			FirstName: "John",
			LastName:  "Smith",
			Age:       25,
			About:     "I love to go rock climbing 3",
			Interests: []string{"sports", "music"},
		}, {
			FirstName: "Jane",
			LastName:  "Smith",
			Age:       32,
			About:     "I like to collect rock albums 3",
			Interests: []string{"music"},
		},
		{
			FirstName: "Douglas",
			LastName:  "Fir",
			Age:       35,
			About:     "I like to build cabinets 3",
			Interests: []string{"forestry"},
		},
	}

	for index, doc := range docs {
		body, err := yelastic.NewES().PutDocument(ctx, "megacorp", "employee", cast.ToString(index+5), doc)
		fmt.Println(string(body), err)
	}
}
func getByID() {
	body, err := yelastic.NewES().GetByID(ctx, "megacorp", "employee", "1")
	fmt.Println(body, err)
}
func getByIDs() {
	body, err := yelastic.NewES().MgetByIds(ctx, "megacorp", "employee", []string{"1", "2", "3"})
	fmt.Println(string(body), err)
}
func multiSearch() {
	body, err := yelastic.NewES().WordMultiSearch(ctx, "John", []string{"first_name"}, "megacorp", "employee")
	fmt.Println(string(body), err)
}
func wordSearch() {
	res, err := yelastic.NewES().PhrasseMatch(ctx, "about", "rock albums", "megacorp", "employee")
	fmt.Println(string(res), err)
}
func delete() {
	res, err := yelastic.NewES().Delete(ctx, "megacorp", "employee", "5")
	fmt.Println(string(res), err)
}
