package service
import (
	"errors"
	"time"
	"autosmm/internal/models"
	"autosmm/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)
type AuthService struct {
	userRepo *repository.UserRepository
	jwtKey   []byte
}
func NewAuthService(repo *repository.UserRepository, key string) *AuthService {
	return &AuthService{userRepo: repo, jwtKey: []byte(key)}
}
func (s *AuthService) Register(username, email, password, fullName string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil { return nil, err }
	user := &models.User{
		ID: uuid.New(), Username: username, Email: email,
		PasswordHash: string(hashedPassword), FullName: fullName,
		Role: "user", Plan: "freemium",
	}
	if err := s.userRepo.Create(user); err != nil { return nil, err }
	return user, nil
}
func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil { return "", errors.New("invalid credentials") }
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil { return "", errors.New("invalid credentials") }
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID, "role": user.Role, "exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString(s.jwtKey), nil
}
func (s *AuthService) ValidateToken(tokenString string) (uuid.UUID, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { return s.jwtKey, nil })
	if err != nil || !token.Valid { return uuid.Nil, "", errors.New("invalid token") }
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok { return uuid.Nil, "", errors.New("invalid claims") }
	userID, _ := uuid.Parse(claims["user_id"].(string))
	return userID, claims["role"].(string), nil
}
