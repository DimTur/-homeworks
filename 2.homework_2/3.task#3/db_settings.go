package main

import (
	"encoding/json"
	"log"
	"os"
)

type Students struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Grade int    `json:"grade"`
}

type Objects struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Results struct {
	Object_id  int `json:"object_id"`
	Student_id int `json:"student_id"`
	Result     int `json:"result"`
}

type Database struct {
	Students []Students `json:"students"`
	Objects  []Objects  `json:"objects"`
	Results  []Results  `json:"results"`
}

func getDb(dbUrl string) Database {
	file, err := os.Open(dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var database Database
	err = decoder.Decode(&database)
	if err != nil {
		log.Fatal(err)
	}
	return database
}
