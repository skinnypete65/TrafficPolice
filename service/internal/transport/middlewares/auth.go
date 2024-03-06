package middlewares

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/services"
	"TrafficPolice/internal/tokens"
	"errors"
	"log"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
)

type AuthMiddleware struct {
	tokenManager  tokens.TokenManager
	expertService services.ExpertService
}

func NewAuthMiddleware(tokenManager tokens.TokenManager, expertService services.ExpertService) *AuthMiddleware {
	return &AuthMiddleware{
		tokenManager:  tokenManager,
		expertService: expertService,
	}
}

func (h *AuthMiddleware) IdentifyRole(next http.Handler, roles ...domain.Role) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(authorizationHeader)
		tokenInfo, err := h.parseAuthHeader(authHeader)

		if err != nil {
			authError(w)
			return
		}

		hasPermission := false
		for _, role := range roles {
			if tokenInfo.UserRole == role {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			authError(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *AuthMiddleware) IsExpertConfirmed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(authorizationHeader)
		tokenInfo, err := h.parseAuthHeader(authHeader)

		if err != nil {
			authError(w)
			return
		}

		if tokenInfo.UserRole != domain.ExpertRole {
			next.ServeHTTP(w, r)
		}

		expert, err := h.expertService.GetExpertByUserID(tokenInfo.UserID)
		if err != nil {
			log.Println(err)
			authError(w)
			return
		}
		if !expert.IsConfirmed {
			notConfirmedError(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *AuthMiddleware) parseAuthHeader(authHeader string) (tokens.TokenInfo, error) {
	if authHeader == "" {
		return tokens.TokenInfo{}, errors.New("empty auth header")
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return tokens.TokenInfo{}, errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return tokens.TokenInfo{}, errors.New("token is empty")
	}
	accessToken := headerParts[1]

	return h.tokenManager.Parse(accessToken)
}
