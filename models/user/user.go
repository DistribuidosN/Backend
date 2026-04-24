package user

import "time"

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

// UserActivity representa una entrada en el historial de procesamiento
type UserActivity struct {
	BatchID     string    `json:"batch_id"`
	RequestTime time.Time `json:"request_time"`
	Status      string    `json:"status"`
	ImageCount  int       `json:"image_count"`
}

// UserStatistics representa los datos agregados de uso de un usuario
type UserStatistics struct {
	TotalBatches    int `json:"total_batches"`
	TotalImages     int `json:"total_images"`
	ImagesCompleted int `json:"images_completed"`
	ImagesFailed    int `json:"images_failed"`
}

// UserUpdateResponse matches the contract consumed by ServerApp.
type UserUpdateResponse struct {
	Message  string `json:"message"`
	Username string `json:"username"`
	Valid    bool   `json:"valid"`
}
