package user

// UserProfile represents user profile data
type UserProfile struct {
	Username string `json:"username"`
	Status   int    `json:"status"`
}

// UserActivity represents a single user activity record
type UserActivity struct {
	ID        string `json:"id"`
	Action    string `json:"action"`
	Timestamp string `json:"timestamp"`
}

// UserStats represents user usage statistics
type UserStats struct {
	ImagesUploaded int `json:"imagesUploaded"`
	TotalLogins    int `json:"totalLogins"`
}

// UserSearchResponse represents the result of a user search
type UserSearchResponse struct {
	UID      string `json:"uid"`
	Username string `json:"username"`
}
