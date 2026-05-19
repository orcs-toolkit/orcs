package store

import (
	"context"

	"github.com/orcs-toolkit/orcs/services/go-api/internal/domain"
)

type Store interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
	FindUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUsers(ctx context.Context) ([]domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) (*domain.User, error)

	CreatePolicy(ctx context.Context, policy domain.Policy) (domain.Policy, error)
	GetPolicies(ctx context.Context) ([]domain.Policy, error)
	GetPolicy(ctx context.Context, id string) (*domain.Policy, error)
	GetPolicyByRole(ctx context.Context, role string) (*domain.Policy, error)
	UpdatePolicyBanList(ctx context.Context, role string, banList []string) (*domain.Policy, error)
	AddSingleProcess(ctx context.Context, roles []string, processName string, all bool, addToFavorite bool) error
	UpdateSinglePolicy(ctx context.Context, role string, processName string) error
	DeletePolicy(ctx context.Context, id string) (*domain.Policy, error)
	GetRoleWisePolicy(ctx context.Context) ([]map[string]any, error)
	GetFavoriteProcesses(ctx context.Context) ([]string, error)

	GetLogs(ctx context.Context, limit int) ([]domain.LogEntry, error)
}
