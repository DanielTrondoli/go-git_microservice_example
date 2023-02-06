package account

import (
	"context"
	"encoding/json"
	"net/http"
)

type (
	CreateUserResponse struct {
		Ok string `json:"ok"`
	}
	GetUserResponse struct {
		Email string `json:"email"`
	}
)

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
