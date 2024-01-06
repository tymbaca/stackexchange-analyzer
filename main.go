package main

import (
	"log"
	"os"

	"time"

	"github.com/tymbaca/stackexchange-analyzer/database"
	"github.com/tymbaca/stackexchange-analyzer/puller"
)

func MustGetenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("env var %s is empty: please set it", key)
	}
	return val
}

func main() {
	db := database.New()
	log.Printf("successfully connected to clickhouse")

	pullerToken := MustGetenv("STACK_TOKEN")
	pullerKey := MustGetenv("STACK_KEY")

	cfg := puller.NewConfig(pullerToken, pullerKey).
		WithInterval(0 * time.Second).
		WithTimeout(20 * time.Second)

	p := puller.NewPuller(cfg)

	maxDate := time.Now()
	minDate := maxDate.Add(-24 * time.Hour)
	questions, err := p.GetQuestinsByDateRange(minDate, maxDate)
	if err != nil {
		panic(err)
	}

	err = db.PushQuestions(questions)
	if err != nil {
		panic(err)
	}

	// file, err := os.Create("output.txt")
	// if err != nil {
	// 	panic(err)
	// }

	// for i, q := range questions {
	// 	fmt.Fprintf(file, "%d\t %s\n", i+1, q.Title)
	// }
}
