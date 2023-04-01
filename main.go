package LogAnalysis

import "sync"

func main() {
	var sy sync.WaitGroup
	filecontrol := &FileControl{
		Path: []string{"./test.log", "./redis.log", "./testA.log", "./testC.log"},
		//防止不够多的情况出现
		LineChan: make(chan []byte),
	}
	sy.Add(1)
	go filecontrol.Init()
	go LogGetChan(filecontrol.LineChan, len(filecontrol.Path), &sy)
	sy.Wait()
	//time.Sleep(10 * time.Second)
	analysis()
}
