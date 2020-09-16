package event_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ariefsam/event"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDispatch(t *testing.T) {
	// LogName := "storage/event.log"
	// f, err := fls.OpenFile(LogName, os.O_CREATE|os.O_WRONLY, 0600)
	// if err != nil {
	// 	log.Println(err)
	// }
	// x, err := f.SeekLine(1, io.SeekStart)
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(x)
	// defer f.Close()
	// event.Loader()
	go event.EventWriter()

	var ev event.Event

	go func() {
		var i int
		for i = 0; i <= 500000; i++ {
			ev.ID = primitive.NewObjectID().Hex()
			ev.Name = "ACreated"
			ev.Payload = "anuani 1 " + fmt.Sprintln(i)
			ev.Timestamp = time.Now().UnixNano()
			event.Dispatch(ev)
		}
	}()
	go func() {
		var i int
		for i = 0; i <= 500000; i++ {
			ev.ID = primitive.NewObjectID().Hex()
			ev.Name = "ACreated"
			ev.Payload = "anuani 2 " + fmt.Sprintln(i)
			ev.Timestamp = time.Now().UnixNano()
			event.Dispatch(ev)
		}
	}()
	go func() {
		var i int
		for i = 0; i <= 500000; i++ {
			ev.ID = primitive.NewObjectID().Hex()
			ev.Name = "ACreated"
			ev.Payload = "anuani 3" + fmt.Sprintln(i)
			ev.Timestamp = time.Now().UnixNano()
			event.Dispatch(ev)
		}
	}()
	go func() {
		var i int
		for i = 0; i <= 500000; i++ {
			ev.ID = primitive.NewObjectID().Hex()
			ev.Name = "ACreated"
			ev.Payload = "anuani 4" + fmt.Sprintln(i)
			ev.Timestamp = time.Now().UnixNano()
			event.Dispatch(ev)
		}
	}()
	event.Loader()

}
