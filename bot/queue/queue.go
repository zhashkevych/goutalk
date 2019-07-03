package queue

import (
	"github.com/zhashkevych/goutalk/bot/nlu"
	"sync"
)

type message struct {
	text   string
	roomID string
	userID string
}

func newMessage(text, roomID, userID string) *message {
	return &message{
		text:   text,
		roomID: roomID,
		userID: userID,
	}
}

type Result struct {
	RoomID      string
	UserID      string
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

func (q *Queue) Enqueue(text, roomID, userID string) {
	msg := newMessage(text, roomID, userID)
	q.queue <- msg
}
