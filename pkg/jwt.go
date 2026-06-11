package pkg

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewClaims(id int, email string, role string) *Claims {
	return &Claims{
		UserId: id,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    os.Getenv("JWT_ISSUER"),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		},
	}
}

func (c *Claims) GenJWT() (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("missing secret key")
	}
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return unsignedToken.SignedString([]byte(jwtSecret))
}

func (c *Claims) VerifyJWT(token string) error {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return errors.New("missing secret key")
	}
	jwtToken, err := jwt.ParseWithClaims(token, c, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return err
	}

	if !jwtToken.Valid {
		return jwt.ErrTokenExpired
	}

	issuer, err := jwtToken.Claims.GetIssuer()
	if err != nil {
		return err
	}

	if issuer != os.Getenv("JWT_ISSUER") {
		return jwt.ErrTokenInvalidIssuer
	}

	return nil
}
