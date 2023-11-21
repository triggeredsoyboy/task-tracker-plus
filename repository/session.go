package repository

import (
	"first-project/model"
	"time"

	"gorm.io/gorm"
)

type SessionRepository interface {
	AddSession(session model.Session) error
	UpdateSession(session model.Session) error
	DeleteSession(token string) error
	SessionAvailEmail(email string) (model.Session, error)
	SessionAvailToken(token string) (model.Session, error)
	TokenExpired(session model.Session) bool
}

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepo(db *gorm.DB) *sessionRepository {
	return &sessionRepository{db}
}

func (r *sessionRepository) AddSession(session model.Session) error {
	err := r.db.Create(&session).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *sessionRepository) UpdateSession(session model.Session) error {
	err := r.db.Where("email = ?", session.Email).Save(&session).Error
	if err != nil {
		return err
	}
	
	return nil
}

func (r *sessionRepository) DeleteSession(token string) error {
	var session model.Session
	err := r.db.Where("token = ?", token).Delete(&session).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *sessionRepository) SessionAvailEmail(email string) (model.Session, error) {
	var sessionEmail model.Session
	err := r.db.Where("email = ?", email).First(&sessionEmail).Error
	if err != nil {
		return model.Session{}, err
	}

	return sessionEmail, nil
}

func (r *sessionRepository) SessionAvailToken(token string) (model.Session, error) {
	var sessionToken model.Session
	err := r.db.Where("token = ?", sessionToken.Token).First(&sessionToken).Error
	if err != nil {
		return model.Session{}, err
	}

	return sessionToken, err
}

func (r *sessionRepository) TokenExpired(session model.Session) bool {
	return session.Expiry.Before(time.Now())
}

func (r *sessionRepository) TokenValidity(token string) (model.Session, error) {
	session, err := r.SessionAvailToken(token)
	if err != nil {
		return model.Session{}, err
	}

	if r.TokenExpired(session) {
		err := r.DeleteSession(token)
		if err != nil {
			return model.Session{}, err
		}
		return model.Session{}, err
	}

	return session, nil
}