package application

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"autopark/internal/domain"
)

const (
	Read   = "read"
	Create = "create"
	Update = "update"
	Delete = "delete"
)

type AuthService struct {
	repo   domain.Repository
	secret []byte
}

type AuthUser struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	Role        string `json:"role"`
	EmployeeID  *int64 `json:"employee_id,omitempty"`
	ExpiresAt   int64  `json:"expires_at,omitempty"`
	Permissions RBAC   `json:"permissions"`
}

type AuthSession struct {
	Token       string   `json:"token"`
	User        AuthUser `json:"user"`
	Permissions RBAC     `json:"permissions"`
}

// RBAC - role based access control.
type RBAC struct {
	Entities map[string][]string `json:"entities"`
	Reports  map[string]bool     `json:"reports"`
}

func NewAuthService(repo domain.Repository, secret string) *AuthService {
	if secret == "" {
		secret = "dev-secret"
	}
	return &AuthService{repo: repo, secret: []byte(secret)}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (AuthSession, error) {
	user, err := s.repo.FindUserByUsername(ctx, strings.TrimSpace(username))
	if err != nil {
		return AuthSession{}, errors.New("invalid username or password")
	}
	if hashPassword(password) != user.PasswordHash {
		return AuthSession{}, errors.New("invalid username or password")
	}

	authUser := AuthUser{
		ID:          user.ID,
		Username:    user.Username,
		Role:        user.Role,
		EmployeeID:  user.EmployeeID,
		ExpiresAt:   time.Now().Add(24 * time.Hour).Unix(),
		Permissions: PermissionsForRole(user.Role),
	}
	token, err := s.sign(authUser)
	if err != nil {
		return AuthSession{}, err
	}
	return AuthSession{Token: token, User: authUser, Permissions: authUser.Permissions}, nil
}

func (s *AuthService) Parse(token string) (AuthUser, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return AuthUser{}, errors.New("invalid token")
	}
	expected := s.signature(parts[0])
	if !hmac.Equal([]byte(expected), []byte(parts[1])) {
		return AuthUser{}, errors.New("invalid token signature")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return AuthUser{}, err
	}
	var user AuthUser
	if err := json.Unmarshal(payload, &user); err != nil {
		return AuthUser{}, err
	}
	if user.ExpiresAt < time.Now().Unix() {
		return AuthUser{}, errors.New("token expired")
	}
	user.Permissions = PermissionsForRole(user.Role)
	return user, nil
}

func (s *AuthService) sign(user AuthUser) (string, error) {
	payload, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	encoded := base64.RawURLEncoding.EncodeToString(payload)
	return encoded + "." + s.signature(encoded), nil
}

func (s *AuthService) signature(encodedPayload string) string {
	mac := hmac.New(sha256.New, s.secret)
	_, _ = mac.Write([]byte(encodedPayload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func CanEntity(role, entity, action string) bool {
	for _, allowed := range PermissionsForRole(role).Entities[entity] {
		if allowed == action {
			return true
		}
	}
	return false
}

func CanReport(role, report string) bool {
	return PermissionsForRole(role).Reports[report]
}

func PermissionsForRole(role string) RBAC {
	entities := map[string][]string{}
	reports := map[string]bool{}

	addEntities := func(actions []string, names ...string) {
		for _, name := range names {
			entities[name] = mergeActions(entities[name], actions)
		}
	}
	addReports := func(names ...string) {
		for _, name := range names {
			reports[name] = true
		}
	}

	allEntities := []string{
		"positions", "facilities", "vehicle-categories", "parts", "routes", "employees", "drivers", "repairmen", "brigades",
		"vehicles", "buses", "route-taxis", "taxis", "trucks", "aux-vehicles", "route-assignments", "driver-vehicle-assignments",
		"transport-logs", "repairs", "repair-works", "replaced-parts", "part-requests", "part-request-items",
	}

	read := []string{Read}
	write := []string{Read, Create, Update}
	full := []string{Read, Create, Update, Delete}

	switch role {
	case "management":
		addEntities(read, allEntities...)
		addReports("autopark", "drivers", "driver-distribution", "passenger-routes", "mileage", "repairs-summary", "staff-hierarchy", "garage", "vehicle-distribution", "freight", "used-parts", "vehicle-movement", "manager-subordinates", "repairman-works")
	case "workshop_heads":
		addEntities(read, allEntities...)
		addEntities(write, "vehicles", "brigades", "repairs", "repair-works", "replaced-parts", "part-requests", "part-request-items")
		addReports("autopark", "drivers", "driver-distribution", "passenger-routes", "mileage", "repairs-summary", "staff-hierarchy", "garage", "vehicle-distribution", "freight", "used-parts", "vehicle-movement", "manager-subordinates", "repairman-works")
	case "foremen":
		addEntities(read, "employees", "brigades", "vehicles", "routes", "route-assignments", "transport-logs", "repairs", "repair-works", "replaced-parts", "part-requests", "part-request-items", "parts", "facilities")
		addEntities(write, "repairs", "repair-works", "replaced-parts", "part-requests", "part-request-items")
		addReports("drivers", "driver-distribution", "mileage", "repairs-summary", "staff-hierarchy", "vehicle-distribution", "used-parts", "manager-subordinates", "repairman-works")
	case "dispatchers":
		addEntities(read, "vehicles", "vehicle-categories", "drivers", "employees", "facilities", "routes", "route-assignments", "transport-logs")
		addEntities(full, "routes", "route-assignments")
		addEntities(write, "transport-logs", "vehicles")
		addReports("autopark", "drivers", "driver-distribution", "passenger-routes", "mileage", "garage", "vehicle-distribution", "freight")
	case "accounting":
		addEntities(read, "vehicles", "vehicle-categories", "employees", "positions", "repairs", "replaced-parts", "transport-logs", "part-requests", "parts")
		addEntities([]string{Read, Create}, "vehicles")
		addReports("autopark", "mileage", "repairs-summary", "garage", "freight", "used-parts", "vehicle-movement")
	case "hr":
		addEntities(full, "employees", "drivers", "repairmen", "positions")
		addEntities(write, "brigades", "driver-vehicle-assignments")
		addEntities(read, "vehicles", "vehicle-categories", "facilities")
		addReports("drivers", "driver-distribution", "staff-hierarchy", "manager-subordinates", "repairman-works")
	case "drivers_role":
		addEntities(read, "vehicles", "routes", "route-assignments", "transport-logs", "repairs")
		addEntities([]string{Read, Create}, "transport-logs")
		addReports("passenger-routes", "mileage", "repairs-summary", "freight")
	case "repairmen_role":
		addEntities(read, "vehicles", "repairs", "repair-works", "replaced-parts", "part-requests", "parts")
		addEntities([]string{Read, Create}, "repair-works", "replaced-parts")
		addReports("repairman-works")
	}

	return RBAC{Entities: entities, Reports: reports}
}

func mergeActions(existing, add []string) []string {
	seen := map[string]bool{}
	for _, action := range existing {
		seen[action] = true
	}
	for _, action := range add {
		if !seen[action] {
			existing = append(existing, action)
			seen[action] = true
		}
	}
	return existing
}
