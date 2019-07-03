package queue

import (
	"context"
	"sync"
)

type HandleFunc func(ctx context.Context, url string) (string, error)

type message struct {
	url    string
	chatID int64
}

type Result struct {
	ChatID   int64
	Filename string
	Err      error
}

type Queue struct {
	queue   chan *message
	doneWg  *sync.WaitGroup
	handler HandleFunc

	maxProcessTime int64
}

func NewQueue(h HandleFunc, maxProcessTime int64) *Queue {
	return &Queue{
		queue:          make(chan *message),
		doneWg:         new(sync.WaitGroup),
		handler:        h,
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

func (q *Queue) Enqueue(m *tgbotapi.Message) {
	msg := q.toMessage(m)
	q.queue <- msg
}

func (q *Queue) toMessage(m *tgbotapi.Message) *message {
	return &message{
		chatID: m.Chat.ID,
		url:    m.Text,
	}
}