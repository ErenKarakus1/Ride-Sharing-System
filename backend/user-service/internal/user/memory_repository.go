package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sort"
	"strings"
	"sync"
	"time"
)

type MemoryRepository struct {
	mu    sync.RWMutex
	users map[string]User
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		users: make(map[string]User),
	}
}

func (r *MemoryRepository) Create(ctx context.Context, user User) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	user.ID = newID()
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.CreatedAt = now
	user.UpdatedAt = now

	r.users[user.ID] = user
	return user, nil
}

func (r *MemoryRepository) List(ctx context.Context) ([]User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	sort.Slice(users, func(i int, j int) bool {
		return users[i].CreatedAt.Before(users[j].CreatedAt)
	})

	return users, nil
}

func (r *MemoryRepository) Get(ctx context.Context, id string) (User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return User{}, ErrUserNotFound
	}

	return user, nil
}

func (r *MemoryRepository) Update(ctx context.Context, id string, request UpdateUserRequest) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return User{}, ErrUserNotFound
	}

	if request.DisplayName != nil {
		user.DisplayName = strings.TrimSpace(*request.DisplayName)
	}
	if request.PhoneNumber != nil {
		user.PhoneNumber = strings.TrimSpace(*request.PhoneNumber)
	}
	user.UpdatedAt = time.Now().UTC()

	r.users[id] = user
	return user, nil
}

func newID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return time.Now().UTC().Format("20060102150405.000000000")
	}

	return hex.EncodeToString(bytes)
}
