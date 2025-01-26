package auth

import (
	"errors"
	"time"

	t "github.com/joshibbotson/gym-tracker-backend/internal/modules/auth/types"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	LoginOrCreateUser(config t.AuthData) (*t.Session, error)
	CreateLocalUser(name, email, password string) (*t.User, error)
	Login(email, password string) (*t.Session, error)
}

type authService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) LoginOrCreateUser(config t.AuthData) (*t.Session, error) {
	user, _ := s.repo.FindUserByEmail(config.Email)
	if user != nil {
		return s.createOrUpdateSession(user)
	}

	user = &t.User{
		Name:          config.Name,
		Email:         config.Email,
		Surname:       config.Surname,
		FirstName:     config.FirstName,
		AuthId:        config.AuthId,
		PictureUrl:    config.PictureUrl,
		VerifiedEmail: config.VerifiedEmail,
		AuthProvider:  config.AuthProvider,
	}

	user, err := s.repo.InsertUser(*user)
	if err != nil {
		return nil, err
	}

	return s.createOrUpdateSession(user)
}

func (s *authService) CreateLocalUser(name, email, password string) (*t.User, error) {
	user, _ := s.repo.FindUserByEmail(email)
	if user != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := s.hashPassword(password)
	if err != nil {
		return nil, err
	}

	user = &t.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	return s.repo.InsertUser(*user)
}

func (s *authService) Login(email, password string) (*t.Session, error) {
	user, err := s.repo.FindUserByEmail(email)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("incorrect password")
	}

	return s.createOrUpdateSession(user)
}

func (s *authService) createOrUpdateSession(user *t.User) (*t.Session, error) {
	expiresAt := time.Now().Add(24 * time.Hour)
	session, _ := s.repo.FindAndUpdateSession(user.ID, expiresAt)
	if session != nil {
		return session, nil
	}

	sessionID := uuid.New().String()
	newSession := t.Session{
		UserID:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		SessionID: sessionID,
		ExpiresAt: expiresAt,
	}

	return s.repo.CreateSession(newSession)
}

func (s *authService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
