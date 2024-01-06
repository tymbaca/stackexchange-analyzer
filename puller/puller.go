package puller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Config struct {
	apiToken string
	apiKey   string
	interval time.Duration
	timeout  time.Duration
}

func NewConfig(apiToken, apiKey string) Config {
	return Config{
		apiToken: apiToken,
		apiKey:   apiKey,
		interval: 500 * time.Millisecond,
		timeout:  20 * time.Second,
	}
}

func (c Config) WithTimeout(timeout time.Duration) Config {
	c.timeout = timeout
	return c
}

func (c Config) WithInterval(interval time.Duration) Config {
	c.interval = interval
	return c
}

type Puller struct {
	Config
}

func NewPuller(cfg Config) *Puller {
	return &Puller{cfg}
}

func (p *Puller) GetQuestinsPageInDateRange(minDate, maxDate time.Time, page int) (Response, error) {
	// Create request with timeout
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://api.stackexchange.com/2.3/questions",
		nil,
	)
	if err != nil {
		return Response{}, err
	}

	// Fill query params
	q := req.URL.Query()
	q.Add("site", "stackoverflow")
	q.Add("pagesize", fmt.Sprint(100)) // This is a max for API
	q.Add("page", fmt.Sprint(page))
	q.Add("order", "desc")
	q.Add("min", fmt.Sprint(minDate.Unix()))
	q.Add("max", fmt.Sprint(maxDate.Unix()))
	q.Add("access_token", p.apiToken)
	q.Add("key", p.apiKey)
	req.URL.RawQuery = q.Encode()

	// Make request
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return Response{}, err
	}

	// Check status code
	if res.StatusCode != 200 {
		return Response{}, fmt.Errorf("got status code %d with message: %s", res.StatusCode, data)
	}

	// Parse result
	var response Response
	err = json.Unmarshal(data, &response)
	if err != nil {
		return Response{}, err
	}

	if len(response.Items) == 0 {
		fmt.Printf("%#v\n", response)
		return Response{}, fmt.Errorf("Response with 0 items")
	}

	fmt.Printf(
		"[GetQuestinsPageInDateRange] got %d questions (1st ID: %d) for time range %s - %s. quota remaining: %d\n",
		len(response.Items),
		response.Items[0].QuestionId,
		minDate.Format(time.DateOnly),
		maxDate.Format(time.DateOnly),
		response.QuotaRemaining,
	)
	return response, nil
}

func (p *Puller) GetQuestinsByDateRange(minDate, maxDate time.Time) ([]Question, error) {
	questions := []Question{}

	// Get 1st page for that date range
	result, err := p.GetQuestinsPageInDateRange(minDate, maxDate, 1)
	if err != nil {
		return nil, err
	}

	// Get all other pages if any
	for page := 2; result.HasMore; page++ {
		time.Sleep(p.interval)
		result, err = p.GetQuestinsPageInDateRange(minDate, maxDate, page)
		if err != nil {
			return nil, err
		}

		questions = append(questions, result.Items...)
	}

	return questions, nil
}
