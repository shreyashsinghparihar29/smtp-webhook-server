package main

import (
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

type Job struct {
	ID      int64
	Payload EmailMessage
}

var jobQueue = make(chan Job, 100)

// start worker pool
func startWorker() {
	for i := 0; i < 3; i++ {
		go worker(i)
	}
}

func worker(id int) {
	for job := range jobQueue {
		log.Println("Worker", id, "processing email ID:", job.ID)

		success := sendWithRetry(job)

		if success {
			updateStatus(job.ID, "SUCCESS", 0)
		}
	}
}

func sendWithRetry(job Job) bool {
	client := resty.New()

	maxRetry := 3
	delay := time.Second

	for i := 1; i <= maxRetry; i++ {
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(job.Payload).
			Post(*flagWebhook)

		if err == nil && resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
			log.Println("Webhook success for ID:", job.ID)
			return true
		}

		log.Println("Retry", i, "failed for ID:", job.ID)

		time.Sleep(delay)
		delay *= 2 // exponential backoff
	}

	// DEAD LETTER QUEUE
	log.Println("Moving to DLQ:", job.ID)
	updateStatus(job.ID, "FAILED", maxRetry)

	return false
}
