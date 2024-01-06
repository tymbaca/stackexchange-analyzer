package puller

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
