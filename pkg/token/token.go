package token

import (
	pb "auth-service/genproto/user"
	"auth-service/pkg/config"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func ValidateToken(tokenstr string) (bool, error) {
	_, err := ExtractClaims(tokenstr)
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractClaims(tokenstr string) (jwt.MapClaims, error) {
	tokenstr = strings.TrimPrefix(tokenstr, "\"")
	tokenstr = strings.TrimSuffix(tokenstr, "\"")
	token, err := jwt.ParseWithClaims(tokenstr, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Load().ACCES_TOKEN), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token: %s", tokenstr)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("failed to parse token claims")
	}
	return claims, nil
}

type Claims struct {
	ID       string `json:"user_id"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Country  string `json:"nationality"`
	jwt.StandardClaims
}

func GenerateAccessToken(in *pb.LoginResponse1) (string, error) {
	claims := Claims{
		ID:       in.Id,
		Role:     in.Role,
		Email:    in.Email,
		Username: in.UserName,
		Country:  in.Country,
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

func GenerateRefreshToken(in *pb.LoginResponse1) (string, error) {
	claims := Claims{
		ID:       in.Id,
		Role:     in.Role,
		Email:    in.Email,
		Username: in.UserName,
		Country:  in.Country,
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
