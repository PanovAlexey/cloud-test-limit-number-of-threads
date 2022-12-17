package domain

import (
	"log"
	"time"
)

const MOCK_WORKER_TIME = 15

type Task struct {
}

func (t Task) Do() {
	time.Sleep(time.Second * MOCK_WORKER_TIME)

	log.Println("Task handled", time.Now().String())
}
