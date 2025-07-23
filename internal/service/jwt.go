package service

import (
	"context"
	"mini-blog/config"
	"time"

	"mini-blog/pkg/logger"

	"sync"

	"github.com/golang-jwt/jwt/v5"
)

var (
	logInstance *logger.Logger
	logOnce     sync.Once
)

func getLogger() *logger.Logger {
	logOnce.Do(func() {
		var err error
		logInstance, err = logger.NewLogger()
		if err != nil {
			panic(err)
		}
	})
	return logInstance
}

func GenerateToken(ctx context.Context, userID int) (string, error) {
	log := getLogger()
	logger := log.Info().Str("method", "GenerateToken").Int("user_id", userID)

	tokenTTL := config.Get().Jwt.TokenTTL * time.Hour

	initialClaims := jwt.MapClaims{
		"expires_at": jwt.NewNumericDate(time.Now().UTC().Add(tokenTTL)),
		"issued_at":  jwt.NewNumericDate(time.Now().UTC()),
		"user_id":    userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, initialClaims)
	logger.Msg("token generated")

	return token.SignedString([]byte(config.GetEnv().Jwt.SigningKey))
}

func ParseToken(jwtString string) (claims jwt.MapClaims, err error) {
	log := getLogger()

	token, err := jwt.ParseWithClaims(jwtString, &claims, func(token *jwt.Token) (interface{}, error) {
		if sm, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || sm.Name != "HS256" {
			log.Error().Err(err).Msg("token method not match")
			return nil, err
		}
		return []byte(config.GetEnv().Jwt.SigningKey), nil
	})
	if err != nil || !token.Valid {
		log.Error().Err(err).Msg("invalid token")
		return nil, err
	}
	return claims, nil
}
