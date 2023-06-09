package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/romankravchuk/muerta/internal/api/router/params"
	"github.com/romankravchuk/muerta/internal/pkg/auth"
	"github.com/romankravchuk/muerta/internal/pkg/config"
	"github.com/romankravchuk/muerta/internal/pkg/jwt"
	"github.com/romankravchuk/muerta/internal/storage/postgres/models"
	"github.com/romankravchuk/muerta/internal/storage/postgres/role"
	"github.com/romankravchuk/muerta/internal/storage/postgres/user"
	"github.com/romankravchuk/muerta/internal/storage/redis"
)

type JWTCredential struct {
	PrivateKey []byte
	PublicKey  []byte
	TTL        time.Duration
}

type AuthServicer interface {
	SignUpUser(ctx context.Context, payload *params.SignUp) error
	LoginUser(
		ctx context.Context,
		payload *params.Login,
	) (*params.TokenDetails, *params.TokenDetails, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (*params.TokenDetails, error)
	LogoutUser(ctx context.Context, refreshToken, accessTokenUUID string) error
}

type AuthService struct {
	cache        redis.Client
	usrStorage   user.UserStorage
	rlStorage    role.RoleRepositorer
	refreshCreds JWTCredential
	accessCreds  JWTCredential
}

// LogoutUser implements AuthServicer
func (s *AuthService) LogoutUser(ctx context.Context, refreshToken, accessTokenUUID string) error {
	if _, err := s.cache.Del(ctx, accessTokenUUID).Result(); err != nil {
		return fmt.Errorf("failed to delete access token: %w", err)
	}
	return nil
}

// RefreshAccessToken implements AuthServicer
func (s *AuthService) RefreshAccessToken(
	ctx context.Context,
	refreshToken string,
) (*params.TokenDetails, error) {
	tokenPayload, err := jwt.ValidateToken(refreshToken, s.refreshCreds.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}
	if err := s.cache.Get(ctx, tokenPayload.UUID).Err(); err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}
	access, err := jwt.CreateToken(tokenPayload, s.accessCreds.TTL, s.accessCreds.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}
	if err := s.cache.Set(ctx, access.UUID, tokenPayload.UserID, time.Until(time.Unix(access.ExpiresIn, 0))).Err(); err != nil {
		return nil, fmt.Errorf("failed to set access token in redis: %w", err)
	}
	return access, nil
}

// LoginUser implements AuthServicer
func (s *AuthService) LoginUser(
	ctx context.Context,
	payload *params.Login,
) (*params.TokenDetails, *params.TokenDetails, error) {
	model, err := s.usrStorage.FindByName(ctx, payload.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("user not found: %w", err)
	}
	hash := auth.GenerateHashFromPassword(payload.Password, model.Salt)
	if ok := auth.CompareHashAndPassword(payload.Password, model.Salt, hash); !ok {
		return nil, nil, fmt.Errorf("invalid name or password")
	}
	tokenPayload := &params.TokenPayload{
		UserID:   model.ID,
		Username: payload.Name,
		Roles:    []string{},
	}
	for _, role := range model.Roles {
		tokenPayload.Roles = append(tokenPayload.Roles, role.Name)
	}
	access, err := jwt.CreateToken(tokenPayload, s.accessCreds.TTL, s.accessCreds.PrivateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create access token: %w", err)
	}
	refresh, err := jwt.CreateToken(tokenPayload, s.refreshCreds.TTL, s.refreshCreds.PrivateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create refresh token: %w", err)
	}
	now := time.Now()
	if err := s.cache.Set(ctx, access.UUID, model.ID, time.Unix(access.ExpiresIn, 0).Sub(now)).Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to set access token in redis: %w", err)
	}
	if err := s.cache.Set(ctx, refresh.UUID, model.ID, time.Unix(refresh.ExpiresIn, 0).Sub(now)).Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to set refresh token in redis: %w", err)
	}
	return access, refresh, nil
}

// SignUpUser implements AuthServicer
func (s *AuthService) SignUpUser(ctx context.Context, payload *params.SignUp) error {
	if _, err := s.usrStorage.FindByName(ctx, payload.Name); err == nil {
		return fmt.Errorf("user already exists")
	}
	role, err := s.rlStorage.FindByName(ctx, "user")
	if err != nil {
		return fmt.Errorf("failed to find roles: %w", err)
	}
	salt := uuid.New().String()
	hash := auth.GenerateHashFromPassword(payload.Password, salt)
	model := models.User{
		Name:  payload.Name,
		Salt:  salt,
		Roles: []models.Role{role},
		Password: models.Password{
			Hash: hash,
		},
	}
	if err := s.usrStorage.Create(ctx, model); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func New(
	cfg *config.Config,
	repo user.UserStorage,
	roleRepository role.RoleRepositorer,
	redis redis.Client,
) AuthServicer {
	return &AuthService{
		cache:      redis,
		usrStorage: repo,
		rlStorage:  roleRepository,
		refreshCreds: JWTCredential{
			PrivateKey: cfg.RefreshTokenPrivateKey,
			PublicKey:  cfg.RefreshTokenPublicKey,
			TTL:        cfg.RefreshTokenExpiresIn,
		},
		accessCreds: JWTCredential{
			PrivateKey: cfg.AccessTokenPrivateKey,
			PublicKey:  cfg.AccessTokenPublicKey,
			TTL:        cfg.AccessTokenExpiresIn,
		},
	}
}
