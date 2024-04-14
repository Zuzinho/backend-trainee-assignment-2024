package session

import (
	"avito_hr/pkg/user"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

type CustomClaims struct {
	Role user.Role `json:"role"`
	jwt.StandardClaims
}

type SessionsManager struct {
	Config *JWTConfig
}

func NewSessionsManager(config *JWTConfig) *SessionsManager {
	return &SessionsManager{
		Config: config,
	}
}

func (manager *SessionsManager) Pack(role user.Role) (*string, error) {
	expirationTime := time.Now().Add(72 * time.Hour)
	claims := &CustomClaims{
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(manager.Config.Method, claims)

	log.Println(token.Claims)

	tokenString, err := token.SignedString(manager.Config.TokenSecret)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func (manager *SessionsManager) Unpack(inToken string) (*Session, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(inToken, claims, func(token *jwt.Token) (interface{}, error) {
		return manager.Config.TokenSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return &Session{
			Sub: claims.Role,
			Exp: time.Unix(claims.ExpiresAt, 0),
		}, nil
	} else {
		return nil, newInvalidTokenError(inToken)
	}
}
