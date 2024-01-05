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

func GetQuestinsPageInDateRange(start, stop time.Time, page int) (Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf(
			"https://api.stackexchange.com/2.3/questions?page=%d&order=desc&min=%d&max=%d&sort=creation&site=stackoverflow",
			page,
			start.Unix(),
			stop.Unix(),
		),
		nil,
	)
	if err != nil {
		return Response{}, err
	}

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return Response{}, err
	}

	var response Response
	err = json.Unmarshal(data, &response)
	if err != nil {
		return Response{}, err
	}

	// fmt.Printf("%+v\n", response)
	if len(response.Items) == 0 {
		fmt.Printf("%#v\n", response)
		return Response{}, fmt.Errorf("Response with 0 items")
	}

	fmt.Printf(
		"[GetQuestinsPageInDateRange] got %d questions (1st ID: %d) for time range %s - %s. quota remaining: %d\n",
		len(response.Items),
		response.Items[0].QuestionId,
		start.Format(time.DateOnly),
		stop.Format(time.DateOnly),
		response.QuotaRemaining,
	)
	return response, nil
}

func GetQuestinsByDateRange(minDate, maxDate time.Time) ([]Question, error) {
	questions := []Question{}

	// Get 1st page for that date range
	result, err := GetQuestinsPageInDateRange(minDate, maxDate, 1)
	if err != nil {
		return nil, err
	}

	// Get all other pages if any
	for page := 2; result.HasMore; page++ {
		result, err = GetQuestinsPageInDateRange(minDate, maxDate, page)
		if err != nil {
			return nil, err
		}

		questions = append(questions, result.Items...)
	}

	return questions, nil
}

func main() {
	maxDate := time.Now()
	minDate := maxDate.Add(-24 * time.Hour)
	questions, err := GetQuestinsByDateRange(minDate, maxDate)
	if err != nil {
		panic(err)
	}

	for i, q := range questions {
		fmt.Printf("%d\t %s", i+1, q.Title)
	}
}
