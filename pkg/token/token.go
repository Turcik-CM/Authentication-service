package token

import (
	"auth-service/pkg/config"
	"auth-service/pkg/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	ID       string `json:"user_id"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateAccessToken(in models.LoginResponse) (string, error) {
	claims := Claims{
		ID:       in.Id,
		Role:     in.Role,
		Email:    in.Email,
		Username: in.Username,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 10).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString([]byte(config.Load().ACCES_TOKEN))
	if err != nil {
		return "", err
	}

	return str, nil
}

func GenerateRefreshToken(in models.LoginResponse) (string, error) {
	claims := Claims{
		ID:       in.Id,
		Role:     in.Role,
		Email:    in.Email,
		Username: in.Username,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString([]byte(config.Load().REFRESH_TOKEN))
	if err != nil {
		return "", err
	}

	return str, nil
}

func ExtractClaimsRefresh(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Load().REFRESH_TOKEN), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, err
	}

	return claims, nil
}

func ExtractClaimsAccess(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Load().ACCES_TOKEN), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, err
	}

	return claims, nil
}
