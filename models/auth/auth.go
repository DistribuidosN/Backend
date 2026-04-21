package auth

// UserCredentials represents the input for a login operation
type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents the result of a successful authentication
type AuthResponse struct {
	Token  string `json:"token" xml:"token"`
	RoleId int    `json:"role_id" xml:"roleId"`
}

// RegisterRequest represents the user registration data
type RegisterRequest struct {
	Username string `json:"username" xml:"username"`
	Password string `json:"password" xml:"password"`
	Email    string `json:"email" xml:"email"`
	RoleId   int    `json:"role_id" xml:"role_id"`
}

type RegisterResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// ValidateResponse represents the detailed token validation result
type ValidateResponse struct {
	Valid    bool   `json:"valid" xml:"valid"`
	Role     string `json:"role" xml:"role"`
	Username string `json:"username" xml:"username"`
	UserUuid string `json:"user_uuid" xml:"userUuid"`
	RoleId   int    `json:"role_id" xml:"roleId"`
}

// ForgetPwdRequest represents a password recovery request
type ForgetPwdRequest struct {
	Email       string `json:"email"`
	NewPassword string `json:"newPassword"`
}

// ResetPasswordRequest represents a password reset request
type ResetPasswordRequest struct {
	NewPassword string `json:"newPassword"`
}
