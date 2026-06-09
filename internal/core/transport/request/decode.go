package core_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
)

var validate = validator.New()

func DecodeAndValidate(r *http.Request, body any) error {
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return fmt.Errorf("failed to decode json: %w", err)
	}

	if err := validate.Struct(body); err != nil {
		return fmt.Errorf("failed to validate json: %w", core_errors.ErrInvalidArg)
	}

	return nil
}
