package dto

type (
	UserPathRequest struct {
		UserID string `uri:"userId" binding:"required,uuid"`
	}

	AddUserRequest struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
		Age   int    `json:"age" validate:"required"`
	}

	GetUserRequest struct {
		UserID string
	}

	UpdateUserRequest struct {
		Name  *string `json:"name,omitempty"`
		Email *string `json:"email,omitempty"`
		Age   *int    `json:"age,omitempty"`
	}

	DeleteUserRequest struct {
		UserID string
	}
)

type (
	AddUserResponse struct {
		ID string `json:"id"`
	}

	GetUserResponse struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}

	UpdateUserResponse struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}

	DeleteUserResponse struct {
		BaseResponse
	}
)
