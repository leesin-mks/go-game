package file

import (
	"bufio"
	"fmt"
	"os"
)

func Start() {
	fileName :="D:/phone.txt"
	_, err:= os.Stat(fileName)
	if err != nil {
		fmt.Println("打开文件失败", err)
		return
	}
	var file *os.File
	if os.IsNotExist(err) {
		file,_ = os.Create(fileName)
	} else {
		file,_ = os.OpenFile(fileName,os.O_RDWR|os.O_APPEND,0666)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for  i :=0 ;scanner.Scan();i++ {
		str := scanner.Text()
		fmt.Println(i," : " ,str)
	}

}