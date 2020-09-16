package event

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"sync/atomic"

	"github.com/hpcloud/tail"
)

type Event struct {
	ID        string
	Name      string
	Payload   string
	Timestamp int64
}

type EventHandler func(Event)

var EventListener map[string][]EventHandler
var EventChannel chan Event
var TotalEventWritten uint64

func init() {
	EventChannel = make(chan Event, 1)
	EventListener = make(map[string][]EventHandler)
	read, err := ioutil.ReadFile(ConfigFilename)
	if err != nil {
		WriteFilename = "storage/event.log"
	} else {
		fmt.Println(string(read))
		//fmt.Println(path)
		if string(read) != "" {
			WriteFilename = string(read)
		} else {
			WriteFilename = "storage/event.log"
		}
	}
}

//One way register. Cannot Unregister. change code and restart if don't want to register listener
func RegisterListener(EventName string, Handler EventHandler) {
	EventListener[EventName] = append(EventListener[EventName], Handler)
}

func Dispatch(Event Event) {
	EventChannel <- Event
}

var WriteFilename string
var NewFilename string
var ChangeLogFile Event
var cPayload ChangeLogFilePayload

const ConfigFilename string = "storage/filename.cfg"

func init() {
	log.Println(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func EventWriter() {
	// rl := rate.New(1000000, time.Millisecond)
	f, err := os.OpenFile(WriteFilename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	defer func() {
		log.Println("Udahan")
	}()
	for Event := range EventChannel {
		// rl.Wait()

		Event.Timestamp = time.Now().UnixNano()

		eventLog, err := json.Marshal(Event)
		if err != nil {
			log.Println("Failed to save event")
		} else {
			atomic.AddUint64(&TotalEventWritten, 1)

			if _, err = f.WriteString(string(eventLog) + "\n"); err != nil {
				log.Println(err)
			}
			// log.Println(eventLog)
			// log.Println("Total Event Written", TotalEventWritten)

			// if TotalEventWritten%10000 == 0 {

			// 	wg.Add(1)
			// 	go func() {
			// 		uid := primitive.NewObjectID().Hex()
			// 		NewFilename = "storage/event-" + uid
			// 		ChangeLogFile.ID = fmt.Sprintf("%s", uid)
			// 		ChangeLogFile.Name = "ChangeLogFile"
			// 		cPayload.CurrentFilename = WriteFilename
			// 		cPayload.NextFilename = NewFilename
			// 		payload, _ := json.Marshal(cPayload)
			// 		ChangeLogFile.Payload = string(payload)
			// 		eventLog, _ = json.Marshal(ChangeLogFile)
			// 		f, err := os.OpenFile(WriteFilename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			// 		if err != nil {
			// 			panic(err)
			// 		}

			// 		if _, err = f.WriteString(string(eventLog) + "\n"); err != nil {
			// 			panic(err)
			// 		}
			// 		f.Close()
			// 		WriteFilename = NewFilename

			// 		err = ioutil.WriteFile(ConfigFilename, []byte(WriteFilename), 0600)
			// 		if err != nil {
			// 			panic(err)
			// 		}
			// 		wg.Done()
			// 	}()
			// }

		}

	}

}

func Loader() {
	var event_data Event
	var LogName string
	var ChangeLogPayload ChangeLogFilePayload
	LogName = "storage/event.log"
	var lineCount int32
	for true {
		t, err := tail.TailFile(LogName, tail.Config{Follow: true})
		if err != nil {
			panic(err)
		}
		for line := range t.Lines {
			atomic.AddInt32(&lineCount, 1)
			err = json.Unmarshal([]byte(line.Text), &event_data)
			if err != nil {
				log.Println("error read event")
			} else {
				for _, value := range EventListener[event_data.Name] {
					value(event_data)
				}
				if event_data.Name == "ChangeLogFile" {
					err = json.Unmarshal([]byte(event_data.Payload), &ChangeLogPayload)
					LogName = ChangeLogPayload.NextFilename
					t.Stop()
				}
			}
		}
	}
}
