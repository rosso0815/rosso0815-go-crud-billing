package services

import (
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func NewScheduler() {
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Println("error", err.Error())
		return
	}

	// add a job to the scheduler
	j, err := s.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			func(a string, b int) {
				// do things
				log.Println("scheduler", a, b)
			},
			"hello",
			1,
		),
	)
	if err != nil {
		log.Println("scheduler job error:", err)
		return
	}
	// each job has a unique id
	log.Println(j.ID())

	// start the scheduler
	s.Start()

	// block until you are ready to shut down
	// select {}
	// case <-time.After(time.Minute):
	//   log.Println("scheduler near shutdown")
	// }

	// when you're done, shut it down
	// err = s.Shutdown()
	// if err != nil {
	//   // handle error
	// }
}
