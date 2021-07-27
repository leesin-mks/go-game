package main

import (
	"game-server/src/db"
	"game-server/src/test"
	"log"
)

func init() {
	log.Println("GameServer is init...")
}

func main() {
	log.Println("GameServer is starting...")

	log.Println(test.RandName())
	db.Start()

}
