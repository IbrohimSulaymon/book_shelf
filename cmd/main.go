package main

import (
	"book_shelf/config"
	"book_shelf/repository"
	"book_shelf/router"
	"book_shelf/server"
	"log"
	"net"

	_ "book_shelf/domain"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	db, err := repository.Connect(cfg)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	s := server.New(repository.NewRepo(db))

	r := router.InitRouter(s)

	r.Run(net.JoinHostPort(cfg.Host, cfg.Port))
}
