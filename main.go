package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Response struct {
	Items          []Question `json:"items"`
	HasMore        bool       `json:"has_more"`
	QuotaMax       int        `json:"quota_max"`
	QuotaRemaining int        `json:"quota_remaining"`
}

type Question struct {
	Tags             []string `json:"tags"`
	Owner            Owner    `json:"owner"`
	IsAnswered       bool     `json:"is_answered"`
	ViewCount        int      `json:"view_count"`
	AnswerCount      int      `json:"answer_count"`
	Score            int      `json:"score"`
	LastActivityDate int      `json:"last_activity_date"`
	CreationDate     int      `json:"creation_date"`
	QuestionId       int      `json:"question_id"`
	ContentLicense   string   `json:"content_license"`
	Link             string   `json:"link"`
	Title            string   `json:"title"`
}

type Owner struct {
	AccountId    int    `json:"account_id"`
	Reputation   int    `json:"reputation"`
	UserId       int    `json:"user_id"`
	UserType     string `json:"user_type"`
	ProfileImage string `json:"profile_image"`
	DisplayName  string `json:"display_name"`
	Link         string `json:"link"`
}

func GetQuestinsByDateRange(start, stop time.Time) Response {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf(
			"https://api.stackexchange.com/2.3/questions?order=desc&min=%d&max=%d&sort=creation&site=stackoverflow",
			start.Unix(),
			stop.Unix(),
		),
		nil,
	)
	if err != nil {
		panic(err)
	}

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var response Response
	err = json.Unmarshal(data, &response)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%+v\n", response)
	if len(response.Items) == 0 {
		fmt.Printf("%#v\n", response)
	}

	fmt.Printf(
		"got %d questions (1st ID: %d) for time range %s - %s. quota remaining: %d\n",
		len(response.Items),
		response.Items[0].QuestionId,
		start.Format(time.DateOnly),
		stop.Format(time.DateOnly),
		response.QuotaRemaining,
	)
	return response
}

func main() {
	maxDate := time.Now()
	minDate := maxDate.Add(-24 * time.Hour)
	for {
		GetQuestinsByDateRange(minDate, maxDate)
		maxDate = minDate
		minDate = minDate.Add(-24 * time.Hour)
	}
}
