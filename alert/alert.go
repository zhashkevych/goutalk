package alert

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/zhashkevych/scheduler"
	"time"
)

type ExecFunc func(roomID, message string) error

type Alerter struct {
	executor  ExecFunc
	scheduler *scheduler.Scheduler
}

type task struct {
	roomID     string
	userID     string
	bookingID  string
	timeDiff   time.Duration
	bookedTime time.Time
}

func NewAlerter(executor ExecFunc, scheduler *scheduler.Scheduler) *Alerter {
	return &Alerter{
		executor:  executor,
		scheduler: scheduler,
	}
}

func (a *Alerter) AddTask(roomID, userID, bookingID string, timeDiff time.Duration, bookedTime time.Time) {
	t := &task{
		roomID:     roomID,
		userID:     userID,
		bookingID:  bookingID,
		timeDiff:   timeDiff,
		bookedTime: bookedTime,
	}

	logrus.Printf("Adding task: %+v", t)

	a.scheduler.Add(context.Background(), func(ctx context.Context) {
		if t.checkTime() {
			logrus.Printf("%d - %d", time.Now().Unix(), t.bookedTime.Add(-t.timeDiff).Unix())
			a.executor(t.roomID, t.newMessage())
		}
	}, time.Second*5)

}

func (t *task) checkTime() bool {
	return time.Now().Unix() >= t.bookedTime.Add(t.timeDiff).Unix()
}

func (t *task) newMessage() string {
	return "Reminder!!! You have a booking " + t.bookingID + " in an " + t.bookedTime.Sub(time.Now()).String()
}
