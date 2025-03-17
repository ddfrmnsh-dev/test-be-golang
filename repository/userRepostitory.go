package repository

import (
	"api-backend-go/model"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll() ([]model.User, error)
	Save(user model.User) (model.User, error)
	FindById(id int) (model.User, error)
	FindBySingle(column, value string) (model.User, error)
	Delete(id int) (model.User, error)
	Update(user model.User) (model.User, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (u *UserRepositoryImpl) FindAll() ([]model.User, error) {
	var users []model.User

	res := u.db.Find(&users)

	if res.Error != nil {
		fmt.Println("Err db:", res.Error)
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		fmt.Println("No users found")
		return users, nil
	}

	return users, nil
}

func (u *UserRepositoryImpl) Save(user model.User) (model.User, error) {
	res := u.db.Create(&user)
	if res.Error != nil {
		fmt.Println("Err db:", res.Error)
		return user, res.Error
	}
	return user, nil
}

func (u *UserRepositoryImpl) FindById(id int) (model.User, error) {
	var user model.User

	res := u.db.First(&user, id)

	if res.Error != nil {
		fmt.Println("Err db:", res.Error)
		return user, res.Error
	}

	if res.RowsAffected == 0 {
		fmt.Println("No user found")
		return user, nil
	}

	return user, nil
}

func (u *UserRepositoryImpl) FindBySingle(column, value string) (model.User, error) {
	var user model.User

	res := u.db.Where(fmt.Sprintf("%s = ?", column), value).First(&user)
	if res.Error != nil {
		return user, res.Error
	}

	if res.RowsAffected == 0 {
		return user, gorm.ErrRecordNotFound
	}

	return user, nil
}

func (u *UserRepositoryImpl) Delete(id int) (model.User, error) {
	checkId, err := u.FindById(id)
	if err != nil {
		return checkId, err
	}

	var user model.User
	res := u.db.Delete(&user, id)
	if res.Error != nil {
		fmt.Println("Err db:", res.Error)
		return user, res.Error
	}
	return checkId, nil
}

func (u *UserRepositoryImpl) Update(user model.User) (model.User, error) {
	res := u.db.Updates(&user)

	if res.Error != nil {
		return user, res.Error
	}

	return user, nil
}
