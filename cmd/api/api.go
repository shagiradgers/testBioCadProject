package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"testBaioCadProject/pkg/database"
	"testBaioCadProject/pkg/server"
)

func main() {
	var (
		dbConfig   database.Config
		servConfig server.Config
	)

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(data, &dbConfig); err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(data, &servConfig); err != nil {
		panic(err)
	}

	db := database.NewDatabase(dbConfig)
	if err := db.Init(); err != nil {
		panic(err)
	}

	handl := server.NewHandler(db)

	fmt.Println("Running server at:", servConfig.Host+":"+servConfig.Port)

	serv := server.NewServer(servConfig, handl)
	if err := serv.Run(); err != nil {
		panic(err)
	}

}
