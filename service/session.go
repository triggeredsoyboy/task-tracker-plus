package service

import (
	"first-project/model"
	repo "first-project/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type SessionService interface {
	GetSessionByEmail(email string) (model.Session, error)
	GenerateToken(user *model.User) (string, error)
}

type sessionService struct {
	sessionRepo repo.SessionRepository
}

func NewSessionService(sessionRepo repo.SessionRepository) *sessionService {
	return &sessionService{sessionRepo}
}

func (s *sessionService) GetSessionByEmail(email string) (model.Session, error) {
	result, err := s.sessionRepo.SessionAvailEmail(email)
	if err != nil {
		return model.Session{}, err
	}

	return result, nil
}

func (s *sessionService) GenerateToken(user *model.User) (string, error) {
	expTime := time.Now().Add(20 * time.Minute)
	claims := &model.Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString(model.JwtKey)
	if err != nil {
		return "", err
	}

	session := model.Session{
		Token:  tokenString,
		Email:  user.Email,
		Expiry: expTime,
	}

	existingSession, err := s.sessionRepo.SessionAvailEmail(session.Email)
	if err != nil {
		s.sessionRepo.AddSession(session)
	} else {
		existingSession.Token = session.Token
		existingSession.Expiry = session.Expiry
		s.sessionRepo.UpdateSession(existingSession)
	}

	return tokenString, nil
}