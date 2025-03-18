package usecase

import (
	"api-backend-go/model"
	"api-backend-go/repository"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	FindAllUser(page, limit int) ([]model.User, int64, error)
	CreateUser(input model.User) (model.User, error)
	FindUserById(id int) (model.User, error)
	FindUserByEmail(email string) (model.User, error)
	UpdateUser(id model.GetCustomerDetailInput, user model.User) (model.User, error)
	DeleteUserById(id int) (model.User, error)
	StartWorkerPool(workerCount int)
}

type userUseCaseImpl struct {
	userRepository repository.UserRepository
	resultChannel  chan UserResult
	jobQueue       chan model.User
}

type UserResult struct {
	User  model.User
	Error error
}

func NewUserUseCase(userRepository repository.UserRepository) UserUseCase {
	return &userUseCaseImpl{userRepository: userRepository, resultChannel: make(chan UserResult, 10), jobQueue: make(chan model.User, 10)}
}

func (uc *userUseCaseImpl) FindAllUser(page, limit int) ([]model.User, int64, error) {
	users, total, err := uc.userRepository.FindAll(page, limit)
	if err != nil {
		return users, 0, err
	}

	return users, total, nil
}

func (uc *userUseCaseImpl) StartWorkerPool(workerCount int) {
	resultChannel := make(chan UserResult, workerCount)
	for i := 1; i <= workerCount; i++ {
		go uc.worker(i, resultChannel)
	}
}

func (uc *userUseCaseImpl) worker(workerID int, resultChannel chan<- UserResult) {
	for user := range uc.jobQueue {
		savedUser, err := uc.userRepository.Save(user)
		if err != nil {
			fmt.Printf("Worker %d failed to save user: %s, error: %v\n", workerID, user.Username, err)
			continue
		}

		uc.resultChannel <- UserResult{User: savedUser, Error: err}
	}
}

func (uc *userUseCaseImpl) CreateUser(input model.User) (model.User, error) {
	user := model.User{}

	_, err := uc.userRepository.FindBySingle("email", input.Email)

	if err == nil {
		return user, errors.New("email sudah digunakan")
	}

	if !strings.Contains(input.Email, "@") {
		return user, fmt.Errorf("email tidak valid")
	}

	if len(input.Password) < 8 {
		return user, fmt.Errorf("password minimal 8 karakter")
	}

	var missingFields []string

	if strings.TrimSpace(input.Username) == "" {
		missingFields = append(missingFields, "Username")
	}

	if strings.TrimSpace(input.Password) == "" {
		missingFields = append(missingFields, "Password")
	}

	if strings.TrimSpace(input.Email) == "" {
		missingFields = append(missingFields, "Email")
	}

	if len(missingFields) > 0 {
		return user, fmt.Errorf("inputan %s tidak boleh string kosong", strings.Join(missingFields, ", "))
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}

	user.Username = strings.ToLower(input.Username)
	user.Email = strings.ToLower(input.Email)
	user.Password = string(hashPassword)
	user.Role = input.Role
	user.IsActive = input.IsActive

	uc.jobQueue <- user

	result := <-uc.resultChannel

	return result.User, result.Error
	// saveUser, err := uc.userRepository.Save(user)

	// if err != nil {
	// 	return user, err
	// }

	// return saveUser, nil
}

func (uc *userUseCaseImpl) FindUserById(id int) (model.User, error) {
	user, err := uc.userRepository.FindById(id)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (uc *userUseCaseImpl) FindUserByEmail(email string) (model.User, error) {
	user, err := uc.userRepository.FindBySingle("email", email)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (uc *userUseCaseImpl) UpdateUser(inputId model.GetCustomerDetailInput, user model.User) (model.User, error) {
	updateId, _ := strconv.Atoi(inputId.Id)

	checkUser, err := uc.userRepository.FindById(updateId)
	if err != nil {
		return checkUser, err
	}

	if strings.TrimSpace(user.Username) != "" {
		checkUser.Username = user.Username
	}

	if strings.TrimSpace(user.Email) != "" {
		if !strings.Contains(user.Email, "@") {
			return user, fmt.Errorf("email tidak valid")
		}

		checkUser.Email = user.Email
	}

	if strings.TrimSpace(user.Password) != "" {
		if len(user.Password) < 8 {
			return user, fmt.Errorf("password minimal 8 karakter")
		}

		if strings.TrimSpace(user.Password) == "" {
			return user, fmt.Errorf("password tidak boleh string kosong")
		}

		hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return user, err
		}

		checkUser.Password = string(hashPassword)
	}

	if strings.TrimSpace(user.Role) != "" {
		checkUser.Role = user.Role
	}

	if user.IsActive != nil {
		checkUser.IsActive = user.IsActive
	}

	checkUser.UpdatedAt = time.Now()
	updatedUser, err := uc.userRepository.Update(checkUser)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

func (uc *userUseCaseImpl) DeleteUserById(id int) (model.User, error) {
	checkUser, err := uc.userRepository.FindById(id)
	if err != nil {
		return checkUser, fmt.Errorf("user tidak ditemukan")
	}

	deleteUser, err := uc.userRepository.Delete(id)
	if err != nil {
		return deleteUser, err
	}

	return deleteUser, nil
}
