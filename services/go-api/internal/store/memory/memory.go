package memory

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/orcs-toolkit/orcs/services/go-api/internal/domain"
)

type Store struct {
	mu         sync.RWMutex
	users      map[string]domain.User
	policies   map[string]domain.Policy
	logs       []domain.LogEntry
	favorites  []string
	idSequence int64
}

func New() *Store {
	return &Store{
		users:    make(map[string]domain.User),
		policies: make(map[string]domain.Policy),
		logs:     []domain.LogEntry{{Timestamp: time.Now().Format(time.RFC3339), Level: "info", Message: "go-api initialized", Meta: map[string]any{"source": "go-api"}}},
	}
}

func (s *Store) nextID(prefix string) string {
	s.idSequence++
	return fmt.Sprintf("%s-%d", prefix, s.idSequence)
}

func (s *Store) CreateUser(_ context.Context, user domain.User) (domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, u := range s.users {
		if strings.EqualFold(u.Email, user.Email) {
			return domain.User{}, errors.New("email already registered")
		}
	}
	user.ID = s.nextID("usr")
	if user.Date.IsZero() {
		user.Date = time.Now().UTC()
	}
	user.Role = strings.ToLower(strings.TrimSpace(user.Role))
	if user.Role == "" {
		user.Role = "default"
	}
	s.users[user.ID] = user
	return user, nil
}

func (s *Store) FindUserByEmail(_ context.Context, email string) (*domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.users {
		if strings.EqualFold(u.Email, email) {
			copy := u
			return &copy, nil
		}
	}
	return nil, nil
}

func (s *Store) FindUserByID(_ context.Context, id string) (*domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return nil, nil
	}
	copy := u
	return &copy, nil
}

func (s *Store) GetUsers(_ context.Context) ([]domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.User, 0, len(s.users))
	for _, u := range s.users {
		out = append(out, u)
	}
	return out, nil
}

func (s *Store) UpdateUser(_ context.Context, user domain.User) (*domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.users[user.ID]
	if !ok {
		return nil, nil
	}
	if user.Name != "" {
		existing.Name = user.Name
	}
	if user.Email != "" {
		existing.Email = user.Email
	}
	if user.Role != "" {
		existing.Role = user.Role
	}
	existing.IsAdmin = user.IsAdmin
	if user.Password != "" {
		existing.Password = user.Password
	}
	s.users[user.ID] = existing
	copy := existing
	return &copy, nil
}

func (s *Store) DeleteUser(_ context.Context, id string) (*domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.users[id]
	if !ok {
		return nil, nil
	}
	delete(s.users, id)
	copy := u
	return &copy, nil
}

func (s *Store) CreatePolicy(_ context.Context, policy domain.Policy) (domain.Policy, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	policy.ID = s.nextID("pol")
	policy.Role = strings.ToLower(strings.TrimSpace(policy.Role))
	if policy.CreatedAt.IsZero() {
		policy.CreatedAt = time.Now().UTC()
	}
	s.policies[policy.ID] = policy
	return policy, nil
}

func (s *Store) GetPolicies(_ context.Context) ([]domain.Policy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.Policy, 0, len(s.policies))
	for _, p := range s.policies {
		out = append(out, p)
	}
	return out, nil
}

func (s *Store) GetPolicy(_ context.Context, id string) (*domain.Policy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.policies[id]
	if !ok {
		return nil, nil
	}
	copy := p
	return &copy, nil
}

func (s *Store) GetPolicyByRole(_ context.Context, role string) (*domain.Policy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.policies {
		if strings.EqualFold(p.Role, role) {
			copy := p
			return &copy, nil
		}
	}
	return nil, nil
}

func (s *Store) UpdatePolicyBanList(_ context.Context, role string, banList []string) (*domain.Policy, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, p := range s.policies {
		if strings.EqualFold(p.Role, role) {
			p.BanList = unique(banList)
			s.policies[id] = p
			copy := p
			return &copy, nil
		}
	}
	return nil, nil
}

func (s *Store) AddSingleProcess(_ context.Context, roles []string, processName string, all bool, addToFavorite bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if all {
		for id, p := range s.policies {
			p.BanList = appendUnique(p.BanList, processName)
			s.policies[id] = p
		}
	} else {
		for id, p := range s.policies {
			for _, role := range roles {
				if strings.EqualFold(p.Role, role) {
					p.BanList = appendUnique(p.BanList, processName)
					s.policies[id] = p
				}
			}
		}
	}
	if addToFavorite {
		s.favorites = appendUnique(s.favorites, processName)
	}
	return nil
}

func (s *Store) UpdateSinglePolicy(_ context.Context, role string, processName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, p := range s.policies {
		if strings.EqualFold(p.Role, role) {
			p.BanList = appendUnique(p.BanList, processName)
			s.policies[id] = p
			return nil
		}
	}
	return nil
}

func (s *Store) DeletePolicy(_ context.Context, id string) (*domain.Policy, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.policies[id]
	if !ok {
		return nil, nil
	}
	delete(s.policies, id)
	copy := p
	return &copy, nil
}

func (s *Store) GetRoleWisePolicy(_ context.Context) ([]map[string]any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]map[string]any, 0, len(s.policies))
	for _, p := range s.policies {
		out = append(out, map[string]any{
			"_id":  p.Role,
			"list": [][]string{p.BanList},
		})
	}
	return out, nil
}

func (s *Store) GetFavoriteProcesses(_ context.Context) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, len(s.favorites))
	copy(out, s.favorites)
	return out, nil
}

func (s *Store) GetLogs(_ context.Context, limit int) ([]domain.LogEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if limit <= 0 || limit > len(s.logs) {
		limit = len(s.logs)
	}
	out := make([]domain.LogEntry, limit)
	copy(out, s.logs[:limit])
	return out, nil
}

func unique(in []string) []string {
	set := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, v := range in {
		k := strings.TrimSpace(strings.ToLower(v))
		if k == "" {
			continue
		}
		if _, ok := set[k]; ok {
			continue
		}
		set[k] = struct{}{}
		out = append(out, v)
	}
	return out
}

func appendUnique(in []string, value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return in
	}
	for _, v := range in {
		if strings.EqualFold(v, value) {
			return in
		}
	}
	return append(in, value)
}
