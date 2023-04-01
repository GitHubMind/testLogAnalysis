package LogAnalysis

import (
	"bufio"
	"io"
	"log"
	"os"
)

type FileControl struct {
	//path
	Path     []string
	LineChan chan []byte
}

// var file  FileControl
//
//	func init() {
//		//初始化他的chan
//	}
func (f FileControl) Init() {
	//获取文本，并且一行一行的读入数据
	for i := 0; i < len(f.Path); i++ {
		go f.transferChan(f.Path[i])
	}

	//最后才把通知报错的关闭 。或者不关也行,防止select没跑完就被关闭了

}
func (f FileControl) transferChan(path string) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		log.Println("open file is failure err :", err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()

		if err != nil {
			log.Println(err)
			if err == io.EOF {

				f.LineChan <- []byte{'E', 'O', 'F'}
			}
			return
		}
		//到末尾了

		if len(line) > 0 {
			//log.Println(line, end)
			f.LineChan <- line
		}
	}

}
