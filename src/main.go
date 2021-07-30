package main

import (
	"game-server/src/db"
	_ "game-server/src/db"
	"game-server/src/file"
	"game-server/src/socket"
	"game-server/src/test"
	"log"
)

func init() {
	log.Println("GameServer is init...")
}

func main() {
	log.Println("GameServer is starting...")

	log.Println(test.RandName())
	// testDB()
	// testFile()
	testSocket()
}

func testDB () {
	db.Start()
	db.Shutdown()
}

func testFile () {
	file.Start()
}

func testSocket () {
	socket.Start()
}
