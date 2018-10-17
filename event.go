package event

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/hpcloud/tail"
	"github.com/teris-io/shortid"
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
var TotalEventWritten uint32

func init() {
	EventChannel = make(chan Event, 1000000)
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

func EventWriter() {

	for Event := range EventChannel {

		Event.Timestamp = time.Now().Unix()

		eventLog, err := json.Marshal(Event)
		if err != nil {
			log.Panic("Failed to save event")
		} else {
			TotalEventWritten++

			f, err := os.OpenFile(WriteFilename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				panic(err)
			}

			if _, err = f.WriteString(string(eventLog) + "\n"); err != nil {
				panic(err)
			}
			f.Close()

			if TotalEventWritten%10000 == 0 {
				sid, _ := shortid.New(1, shortid.DefaultABC, 1)
				uid, _ := sid.Generate()
				NewFilename = "storage/event-" + uid
				ChangeLogFile.ID = fmt.Sprintf("%d", uid)
				ChangeLogFile.Name = "ChangeLogFile"
				cPayload.CurrentFilename = WriteFilename
				cPayload.NextFilename = NewFilename
				payload, _ := json.Marshal(cPayload)
				ChangeLogFile.Payload = string(payload)
				eventLog, _ = json.Marshal(ChangeLogFile)
				f, err := os.OpenFile(WriteFilename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
				if err != nil {
					panic(err)
				}

				if _, err = f.WriteString(string(eventLog) + "\n"); err != nil {
					panic(err)
				}
				f.Close()
				WriteFilename = NewFilename

				err = ioutil.WriteFile(ConfigFilename, []byte(WriteFilename), 0600)
				if err != nil {
					panic(err)
				}
			}

		}

	}
}

func Loader() {
	var event_data Event
	var LogName string
	var ChangeLogPayload ChangeLogFilePayload
	LogName = "storage/event.log"
	for true {
		t, err := tail.TailFile(LogName, tail.Config{Follow: true})
		if err != nil {
			panic(err)
		}
		for line := range t.Lines {
			//fmt.Println(line.Text)
			err = json.Unmarshal([]byte(line.Text), &event_data)
			if err != nil {
				panic("error read event")
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
