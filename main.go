package main

import (
	"fmt"
	"log"
	"net/http"
	// "ticketing-booking/internal/db"
	// "ticketing-booking/internal/kafka"
)

func main() {
	// Load environment variables
	// dbHost := os.Getenv("DB_HOST")
	// dbPort := os.Getenv("DB_PORT")
	// dbUser := os.Getenv("DB_USER")
	// dbPassword := os.Getenv("DB_PASSWORD")
	// dbName := os.Getenv("DB_DATABASE")
	// redisHost := os.Getenv("REDIS_HOST")
	// kafkaBroker := os.Getenv("KAFKA_BROKER")

	// ----------------------
	// Connect to Postgres
	// ----------------------
	// pg, err := db.NewPostgres(dbHost, dbPort, dbUser, dbPassword, dbName)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to Postgres: %v", err)
	// }
	// defer pg.Close()

	// ----------------------
	// Connect to Redis
	// ----------------------
	// rdb := cache.NewRedis(redisHost)

	// ----------------------
	// Connect to Kafka
	// ----------------------
	// producer, err := kafka.NewProducer(kafkaBroker)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to Kafka: %v", err)
	// }
	// defer producer.Close()

	// ----------------------
	// Initialize Booking Service
	// ----------------------
	// service := booking.NewService(pg, rdb, producer)

	// ----------------------
	// Define HTTP routes
	// ----------------------
	// http.HandleFunc("/reserve", service.ReserveTicketHandler)
	// http.HandleFunc("/cancel", service.CancelTicketHandler)

	port := "8080"
	fmt.Printf("Booking service running on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
