package repository

import (
	"context"
	"errors"

	"github.com/DanielTrondoli/go-kit_microservice_example/account"
	"github.com/go-kit/log"
)

var ErrRepo = errors.New("unable to handle repo request")

type memoryRepo struct {
	db     map[string][]string
	logger log.Logger
}

func (repo *memoryRepo) add(new map[string]string) error {

	for key, elem := range new {

		_, ok := repo.db[key]
		if !ok {
			return errors.New("collumn name dont mach")
		}

		repo.db[key] = append(repo.db[key], elem)
	}

	return nil
}

func (repo memoryRepo) getById(id string) (map[string]string, error) {
	elem := make(map[string]string)
	var elemIndex int
	var find bool
	for i, v := range repo.db["id"] {
		if v == id {
			elemIndex = i
			find = true
			break
		}
	}

	if !find {
		return nil, errors.New("element not find")
	}

	for k, e := range repo.db {
		elem[k] = e[elemIndex]
	}

	return elem, nil
}

func NewMemoryRepo(columns []string, logger log.Logger) account.UserRepository {
	db := make(map[string][]string)

	for _, column := range columns {
		db[column] = append(db[column], "")
	}

	return &memoryRepo{
		db:     db,
		logger: log.With(logger, "repo", "In Memory"),
	}
}

func (repo *memoryRepo) CreateUser(ctx context.Context, user account.User) error {
	if user.Email == "" || user.Password == "" {
		return ErrRepo
	}

	userMap := map[string]string{
		"id":       user.ID,
		"email":    user.Email,
		"password": user.Password,
	}

	err := repo.add(userMap)
	if err != nil {
		return err
	}

	return nil
}

func (repo memoryRepo) GetUser(ctx context.Context, id string) (string, error) {

	user, err := repo.getById(id)
	if err != nil {
		return "", err
	}

	email, ok := user["email"]
	if ok {
		return email, nil
	}

	return "", ErrRepo
}
