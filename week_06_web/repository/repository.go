package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"week_06_web/model"
)

var ErrNotFound = errors.New("record not found")

type FileUserRepository struct {
	mu       sync.RWMutex
	users    map[string]model.User
	filepath string
}

func (r *FileUserRepository) load() error {
	data, err := os.ReadFile(r.filepath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read file %s: %w", r.filepath, err)
	}
	if len(data) == 0 {
		return nil
	}

	if err := json.Unmarshal(data, &r.users); err != nil {
		return fmt.Errorf("read file %s: %w", r.filepath, err)
	}
	if r.users == nil {
		r.users = make(map[string]model.User)
	}
	return nil
}

func NewFileUserRepository(filepath string) (*FileUserRepository, error) {
	var repo = &FileUserRepository{
		filepath: filepath,
		users:    make(map[string]model.User),
	}
	if err := repo.load(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *FileUserRepository) saveLocked() error {
	if err := os.MkdirAll(filepath.Dir(r.filepath), 0o755); err != nil {
		return fmt.Errorf("create directory err %w", err)
	}

	tmp, err := os.CreateTemp(filepath.Dir(r.filepath), "users-*.temp")

	if err != nil {
		return fmt.Errorf("create tmp file err %w", err)
	}

	tmpPath := tmp.Name()
	defer os.Remove(tmp.Name())

	if err := json.NewEncoder(tmp).Encode(r.users); err != nil {
		tmp.Close()
		return fmt.Errorf("encode users err %w", err)
	}

	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return fmt.Errorf("sync file err %w", err)
	}

	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close tmp file err %w", err)
	}

	if err := os.Rename(tmpPath, r.filepath); err != nil {
		return fmt.Errorf("replace user file err %w", err)
	}

	return nil
}

func (r *FileUserRepository) Create(_ context.Context, user model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.Id] = user
	return r.saveLocked()
}

func (r *FileUserRepository) FindById(_ context.Context, Id string) (model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[Id]
	if !ok {
		return model.User{}, ErrNotFound
	}

	return user, nil
}

func (r *FileUserRepository) FindByUsername(_ context.Context, username string) (model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return model.User{}, ErrNotFound
}
func (r *FileUserRepository) UpdatePassword(_ context.Context, id string, passwordHash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return ErrNotFound
	}

	user.PasswordHash = passwordHash
	r.users[id] = user
	return r.saveLocked()
}
