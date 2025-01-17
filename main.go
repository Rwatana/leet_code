package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"log_service/internal/server/infrastructure/mysql/db"
	"log_service/internal/server/infrastructure/mysql/repository"
	"log_service/internal/server/infrastructure/rabbitmq"
	"log_service/internal/server/presentation"
	"log_service/internal/server/usecase"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	logRepo := repository.NewLogRepository(db)
	logUseCase := usecase.NewInsertLogUseCase(logRepo)
	logListUseCase := usecase.NewListLogsUseCase(logRepo)
	HttpLogHandler := presentation.NewHttpLogHandler(logListUseCase)

	amqpConn, ch, msgs, err := rabbitmq.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	defer amqpConn.Close()
	defer ch.Close()

	amqpLogHandler := presentation.NewAMQPLogHandler(logUseCase, ch)

	done := make(chan bool)
	go func() {
		for d := range msgs {
			amqpLogHandler.HandleLog(d)
			if err := d.Ack(false); err != nil {
				log.Fatalf("Failed to ack message: %v", err)
			}
		}
		done <- true
	}()

	http.HandleFunc("/logs", HttpLogHandler.HandleLogList)

	go func() {
		log.Println("HTTP server is running on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()


	log.Printf("Waiting for messages. To exit press CTRL^C")

	<-ctx.Done()
	stop()

	log.Println("received sigint/sigterm, shutting down...")
	log.Println("press Ctrl^C again to force shutdown")

	if err := ch.Cancel(rabbitmq.QUEUE_NAME, false); err != nil {
		log.Panic(err)
	}
	if err := ch.Close(); err != nil {
		log.Panic(err)
	}

	select {
	case <-done:
		log.Println("finished processing all jobs")
	case <-time.After(5 * time.Second):
		log.Println("timed out waiting for jobs to finish")
	}
}





