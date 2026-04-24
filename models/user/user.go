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

// TransformationStat representa el conteo de uso de una transformación específica
type TransformationStat struct {
	Name  string `json:"name" xml:"name"`
	Count int    `json:"count" xml:"count"`
}

// UserActivity representa una entrada en el historial de procesamiento (Timeline)
type UserActivity struct {
	EventType   string     `json:"event_type"`
	RefUUID     string     `json:"ref_uuid"`
	ParentUUID  string     `json:"parent_uuid"`
	Description string     `json:"description"`
	OccurredAt  *time.Time `json:"occurred_at"`
	
	// Deprecated: used for old activity format, keeping for compatibility if needed
	BatchID     string     `json:"batch_id,omitempty"`
	RequestTime *time.Time `json:"request_time,omitempty"`
	Status      string     `json:"status,omitempty"`
	ImageCount  int        `json:"image_count,omitempty"`
}

// UserStatistics representa los datos agregados de uso de un usuario
type UserStatistics struct {
	TotalBatches       int                  `json:"total_batches"`
	TotalImages        int                  `json:"total_images"`
	ImagesCompleted    int                  `json:"images_completed"`
	ImagesFailed       int                  `json:"images_failed"`
	TopTransformations []TransformationStat `json:"top_transformations"`
}

// UserUpdateResponse matches the contract consumed by ServerApp.
type UserUpdateResponse struct {
	Message  string `json:"message"`
	Username string `json:"username"`
	Valid    bool   `json:"valid"`
}
