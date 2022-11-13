package main

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"testBaioCadProject/pkg/database"
	"testBaioCadProject/pkg/fileChecker"
	"time"
)

// durationSeconds this is the time after which checker is started
const durationSeconds = 60

func main() {
	var (
		dbConfig          database.Config
		fileCheckerConfig fileChecker.Config
	)

	logger := log.New(os.Stdout, "App: ", 0)

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(data, &dbConfig); err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(data, &fileCheckerConfig); err != nil {
		panic(err)
	}

	db := database.NewDatabase(dbConfig)
	if err := db.Init(); err != nil {
		panic(err)
	}

	checker := fileChecker.NewChecker(db, fileCheckerConfig)

	for {
		logger.Println("Запустли сканнирование")
		if err := checker.RunChecker(); err != nil {
			log.Fatal(err.Error())
		}
		logger.Println("Спим", durationSeconds, "секунд")
		time.Sleep(time.Second * durationSeconds)
	}
}
