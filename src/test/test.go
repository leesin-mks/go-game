package test

import (
	"fmt"
	"math/rand"
	"time"
)

var nameDB = [3]string{"leesin", "root", "robot"}

func RandName() string {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var index = r.Intn(len(nameDB))
	fmt.Println("index: ", index)
	return nameDB[index]
}
