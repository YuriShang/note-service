package user

import (
	"note_service/app/internal/client/user_client"
	"note_service/app/pkg/logging"

	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type UserClaims struct {
	userUUID uuid.UUID
}

func Authentication(c user_client.UserClient, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logging.GetLogger()
		token := user_client.Token{}
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			logger.Error("Malformed token")
		} else {
			token = user_client.Token{AccessToken: authHeader[1], TokenType: "Bearer"}
		}
		user, err := c.GetUserByToken(r.Context(), token)
		if err != nil {
			if err.Error() == "Unauthorized" {
				logger.Error("Token expired")
			}
		}
		var uc = UserClaims{
			userUUID: user.UUID,
		}
		ctx := context.WithValue(r.Context(), "userUUID", uc.userUUID)
		h(w, r.WithContext(ctx))
	}
}

func Authorization(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userUUID := r.Context().Value("userUUID").(uuid.UUID)
		if userUUID.String() == "00000000-0000-0000-0000-000000000000" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		ctx := context.WithValue(r.Context(), "userUUID", userUUID)
		h(w, r.WithContext(ctx))
	}
}
