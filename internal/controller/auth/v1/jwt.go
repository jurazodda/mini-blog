package v1

import (
	"context"
	"mini-blog/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (h *AuthHandlerV1) GenerateToken(ctx context.Context, userID int) (string, error) {
	logger := h.Logger.With().Ctx(ctx).Str("method", "GenerateToken").Int("user_id", userID).Logger()

	tokenTTL := config.Get().Jwt.TokenTTL * time.Hour

	initialClaims := jwt.MapClaims{
		"expires_at": jwt.NewNumericDate(time.Now().UTC().Add(tokenTTL)),
		"issued_at":  jwt.NewNumericDate(time.Now().UTC()),
		"user_id":    userID,
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, initialClaims)
	logger.Info().Msg("token generated")

	return token.SignedString([]byte(config.Get().Jwt.SigningKey))
}

func (h *AuthHandlerV1) ParseToken(jwtString string) (claims jwt.MapClaims, err error) {
	logger := h.Logger.With().Str("method", "ParseToken").Logger()

	token, err := jwt.ParseWithClaims(jwtString, &claims, func(token *jwt.Token) (interface{}, error) {
		if sm, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || sm.Name != "HS256" {
			logger.Error().Err(err).Msg("token method not match")
			return nil, err
		}
		return []byte(config.Get().Jwt.SigningKey), nil
	})
	if err != nil || !token.Valid {
		logger.Error().Err(err).Msg("invalid token")
		return nil, err
	}
	return claims, nil
}
