package service

import (
	"errors"
	"first-project/model"
	repo "first-project/repository"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetCurrentUser(c *gin.Context) *model.User
	Register(user *model.User) (model.User, error)
	Login(user *model.User) (token *string, err error)
	Logout(token string) error
}

type userService struct {
	userRepo repo.UserRepository
	sessionRepo repo.SessionRepository
	sessionService SessionService
}

func NewUserService(userRepo repo.UserRepository, sessionRepo repo.SessionRepository, sessionService SessionService) UserService {
	return &userService{userRepo, sessionRepo, sessionService}
}

func (s *userService) GetCurrentUser(c *gin.Context) *model.User {
	token, err := c.Cookie("session_token")
	if err != nil {
		return nil
	}

	claims := &model.Claims{}
	tokenKey, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return model.JwtKey, nil
	})

	if err != nil || !tokenKey.Valid {
		return nil
	}

	user, err := s.userRepo.GetUserByEmail(claims.Email)
	if err != nil {
		return nil
	}

	return &user
}

func (s *userService) Register(user *model.User) (model.User, error) {
	dbUser, err := s.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return *user, err
	}

	if dbUser.Email != "" || dbUser.ID != 0 {
		return *user, errors.New("email not available")
	}

	// hash the password
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return *user, nil
	}
	
	user.Password = string(hashPass)
	user.CreatedAt = time.Now()

	newUser, err := s.userRepo.CreateUser(*user)
	if err != nil {
		return *user, err
	}

	return newUser, nil
}

func (s *userService) Login(user *model.User) (token *string, err error) {
	dbUser, err := s.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	if dbUser.ID == 0 || dbUser.Email == "" {
		return nil, errors.New("user not found")
	}

	if user.Email != dbUser.Email {
		return nil, errors.New("wrong email")
	}

	// compare sent-in pass with saved hash pass
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return nil, errors.New("wrong password")
	}

	// generate a jwt token
	tokenString, err := s.sessionService.GenerateToken(&dbUser)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func (s *userService) Logout(token string) error {
	err := s.sessionRepo.DeleteSession(token)
	if err != nil {
		return err
	}

	return nil
}