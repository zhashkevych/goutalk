package queue

import (
	"context"
	"time"
)

func (q *Queue) startProcess(results chan *Result) {
	for {
		msg := <-q.queue
		go q.proccesMessageAndRespond(msg, results)
	}
}

func (q *Queue) proccesMessageAndRespond(m *message, results chan *Result) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(q.maxProcessTime))
	defer cancel()

	q.doneWg.Add(1)
	defer q.doneWg.Done()

	response, err := q.processor.Process(ctx, m.text, m.roomID, m.userID)

	results <- &Result{
		RoomID:      m.roomID,
		ResponseMsg: response,
		Err:         err,
	}
}
