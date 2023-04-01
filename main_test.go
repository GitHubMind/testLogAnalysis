package LogAnalysis

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestFile(t *testing.T) {
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

// 并测
func BenchmarkFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
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
		analysis()
	}
}

// 生成测试数据
func TestSetValue(t *testing.T) {
	var methods = []string{"POST", "GET", "PUT", "DELETE", "PATCH"}

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 1000; i++ {
		d := LogData{
			SessionID: rand.Intn(100),
			ID:        rand.Intn(1000),
			SeqNumber: rand.Intn(100),
			Method:    methods[rand.Intn(len(methods))],
			Timestamp: time.Now(),
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "test",
			},
			Payload: map[string]interface{}{
				"param1": rand.Intn(100),
				"param2": fmt.Sprintf("value-%d", i),
			},
			ExtraData: fmt.Sprintf("extra-data-%d", i),
		}

		b, _ := json.Marshal(d)
		fmt.Println(string(b))
	}

}
