package queue

import (
	"github.com/zhashkevych/goutalk/nlu"
	"sync"
)

type message struct {
	text   string
	roomID string
}

func newMessage(text, roomID string) *message {
	return &message{
		text:   text,
		roomID: roomID,
	}
}

type Result struct {
	RoomID      string
	ResponseMsg string
	Err         error
}

type Queue struct {
	queue     chan *message
	doneWg    *sync.WaitGroup
	processor nlu.Processor

	maxProcessTime int64
}

func NewQueue(processor nlu.Processor, maxProcessTime int64) *Queue {
	return &Queue{
		queue:          make(chan *message),
		doneWg:         new(sync.WaitGroup),
		processor:      processor,
		maxProcessTime: maxProcessTime,
	}
}

func (q *Queue) Start(results chan *Result) {
	go q.startProcess(results)
}

func (q *Queue) Stop() {
	q.doneWg.Wait()
	close(q.queue)
}

func (q *Queue) Push(text, roomID string) {
	msg := newMessage(text, roomID)
	q.queue <- msg
}
