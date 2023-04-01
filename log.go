package LogAnalysis

import (
	"bytes"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type LogData struct {
	SessionID int               `json:"session_id"`
	ID        int               `json:"id"`
	SeqNumber int               `json:"seq_number"`
	Method    string            `json:"method"`
	Timestamp time.Time         `json:"timestamp"`
	Headers   map[string]string `json:"headers"`
	Payload   interface{}       `json:"payload"`
	ExtraData string            `json:"extra_data"`
}

// 模式
type Mode struct {
	//序列
	OrderList string
	//统计有几个Session用了这个模式
	SessionList []int
	//通过相同的Mode
	Sum int
}

// 方法调用序列模式包含的方法个数可以调整
const step = 2

var (
	//统计
	SessionMark map[int][]LogData
	//作为模式计算的下表
	MethodName []string
	ModeMark   map[string]Mode
	Mutex      sync.Mutex
)

func init() {
	SessionMark = make(map[int][]LogData, 0)
	ModeMark = make(map[string]Mode, 0)
}

func LogGetChan(lineChan chan []byte, endIndex int, sy *sync.WaitGroup) {
	defer sy.Done()

	for {
		select {
		case val, _ := <-lineChan:
			//这个line记得这样写深复制，不能再翻同样的低级错误
			line := make([]byte, len(val))
			copy(line, val)
			if bytes.Equal(line, []byte{'E', 'O', 'F'}) {
				endIndex--
				if endIndex == 0 {
					return
				}
				continue
			}
			var _log LogData
			err := json.Unmarshal(line, &_log)
			if err != nil {
				//log.Println(string(line))
				log.Println("Unmarshal is fault : ", err)

				continue
			}
			Mutex.Lock()

			//数据装载
			if _, ok := SessionMark[_log.SessionID]; ok {
				SessionMark[_log.SessionID] = append(SessionMark[_log.SessionID], _log)
			} else {
				SessionMark[_log.SessionID] = []LogData{_log}
			}
			if !isHavedMethod(_log.Method) {
				MethodName = append(MethodName, _log.Method)
			}
			Mutex.Unlock()
		}
	}
	//数据导入结束，开始做数据汇总

}
func analysis() {
	//先map一下方法

	for SessionId, data := range SessionMark {
		//最后一位也算进step 里面 所以要-1
		for i := 0; i < len(data)-step+1; i++ {
			//开始具体走
			orderList := ""
			for j := 0; j < step; j++ {
				//如果用数字来优化感觉如果 method>0 要处理一下，肯定比字符串要好不少的
				orderList += strconv.Itoa(getMethodId(data[i+j].Method))
			}
			if val, ok := ModeMark[orderList]; ok {
				val.Sum++
				//不要重复的
				var repeat bool
				for i := 0; i < len(val.SessionList); i++ {
					if val.SessionList[i] == SessionId {
						repeat = true
						break
					}
				}
				if !repeat {
					val.SessionList = append(val.SessionList, SessionId)

				}
				ModeMark[orderList] = val
			} else {
				ModeMark[orderList] = Mode{OrderList: orderList, Sum: 1, SessionList: []int{SessionId}}
			}

		}
	}
	Max := 0
	var MaxResult []string
	log.Println("一共会话有:", len(SessionMark))
	for _, mode := range ModeMark {
		result := make([]string, 0)
		for _, i := range mode.OrderList {
			result = append(result, MethodName[i-'0'])
		}
		str := strings.Join(result, "->")
		if mode.Sum > Max {
			MaxResult = []string{str}
			Max = mode.Sum
		} else if mode.Sum == Max {
			MaxResult = append(MaxResult, str)
		}
		log.Println("该模式:", str, "出现了:", mode.Sum, "次", "并且再 session id :", mode.SessionList, "出现过")
	}

	log.Println("一共会话有:", len(SessionMark))
	log.Println("一共匹配模式有:", len(ModeMark))
	log.Println("最大匹配的模式是:", strings.Join(MaxResult, ","), "次数为:", Max)
}

func getMethodId(method string) int {
	for i := 0; i < len(MethodName); i++ {
		if MethodName[i] == method {
			return i
		}
	}
	return -1
}

func isHavedMethod(method string) (Pass bool) {
	for _, i2 := range MethodName {
		if i2 == method {
			Pass = true
			break
		}
	}
	return
}
