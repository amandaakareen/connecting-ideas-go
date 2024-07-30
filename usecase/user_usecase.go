package usecase

import (
	"github.com/example/dto"
	"github.com/example/entity"
	"github.com/example/repository"
	"github.com/example/service"
	"github.com/example/util"
)

type UserUseCase struct {
	repository *repository.UserRepository
}

func NewUserUsecase(ur *repository.UserRepository) *UserUseCase {
	return &UserUseCase{repository: ur}
}

func (us *UserUseCase) Save(userDto dto.UserCreate) (string, error) {
	err := us.repository.ExistEmail(userDto.Email)
	if err != nil {
		return "", err
	}

	err = us.repository.ExistName(userDto.Name)
	if err != nil {
		return "", err
	}
	code := "12345"

	user := entity.NewUser(userDto.Name, userDto.Email, userDto.Password, "PENDENTE", code)

	util.EncryptPassword(&user)

	err = us.repository.Save(user)
	if err != nil {
		return "", err
	}

	err = service.SendEmail(user.Email, "titulo", "codigo"+code)
	if err != nil {
		return "", err
	}

	return "Enviamos um código para seu email para validarmos sua conta!", nil
}

func (us *UserUseCase) Login(user dto.UserLogin) (string, error) {
	userResult, err := us.repository.FindByEmail(user.Email)
	if err != nil {
		return "", err
	}

	err = util.CheckPassword(userResult.Password, user.Password)
	if err != nil {
		return "", err
	}

	return "autenticado", nil
}
