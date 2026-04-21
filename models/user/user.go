package user

// UserProfile represents user profile data
type UserProfile struct {
	ID        int    `json:"id"`
	UserUUID  string `json:"user_uuid,omitempty"`
	Username  string `json:"username"`
	Email     string `json:"email,omitempty"`
	RoleID    int    `json:"role_id"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at,omitempty"`
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

// UserUpdateResponse matches the contract consumed by ServerApp.
type UserUpdateResponse struct {
	Message  string `json:"message"`
	Username string `json:"username"`
	Valid    bool   `json:"valid"`
}
