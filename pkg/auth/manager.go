package auth

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Config interface {
	GetSecret() string
}

type TokenManager interface {
	NewJWT(userId int, ttl time.Duration) (string, error)
	ParseToken(accessToken string) (string, error)
	NewRefreshToken() (string, error)
}

type Manager struct {
	secret string
}

func NewManager(cfg Config) *Manager {
	return &Manager{
		secret: cfg.GetSecret(),
	}
}

func (m *Manager) NewJWT(userId int, ttl time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		Subject:   strconv.Itoa(userId),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.secret))
}

func (m *Manager) ParseToken(accessToken string) (int, error) {
	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return m.secret, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	id, err := strconv.Atoi(claims["sub"].(string))
	if err != nil {
		return 0, errors.New("invalid subject")
	}

	return id, nil
}

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	_, err := r.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
