package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type Claims struct {
	UserID    uint
	UserEmail string
	Role      string
	SubscriptionExpiry int64
	jwt.StandardClaims
}

func GenerateJwtToken(Email string, ID uint, Role string,Expiry int64) (string, error) {

	claims := Claims{
		UserID:    ID,
		UserEmail: Email,
		Role:      Role,
		SubscriptionExpiry: Expiry,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := viper.GetString("TOKEN_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT secret is not configured")
	}

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %v", err)
	}

	return signedToken, nil
}
