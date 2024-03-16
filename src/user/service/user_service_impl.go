package service

import (
	"errors"
	"fmt"
	"test-mkp/config"
	"test-mkp/exception"
	"test-mkp/src/user/entity"
	"test-mkp/src/user/model"
	"test-mkp/src/user/repository"
	"test-mkp/src/user/validation"
	"test-mkp/utils"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func NewUserService(
	userRepo repository.UserRepository,
	jwtProvider config.JwtConfig,
) UserService {
	validator := new(validation.UserValidation)
	return &userServiceImpl{userRepo, *validator, jwtProvider}
}

type userServiceImpl struct {
	userRepo    repository.UserRepository
	validator   validation.UserValidation
	JwtProvider config.JwtConfig
}

func (s *userServiceImpl) Login(request model.LoginRequest) (model.LoginResponse, error) {
	response := model.LoginResponse{}
	if err := s.validator.LoginValidate(request); err != nil {
		return response, err
	}

	user, err := s.userRepo.FindUserByPhoneOrEmail(request.Identity)
	if err != nil {
		return response, err
	}

	if user == nil {
		return response, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		fmt.Println(err)
		return response, errors.New("password not same")
	}
	jwtToken, err := s.generateJwt(*user)
	if err != nil {
		exception.PanicIfNeeded(exception.ServerError{
			Message: "generate session failed",
		})
	}

	response = model.LoginResponse{
		ID:          user.ID,
		Name:        utils.CheckPointerValue(user.Name),
		Email:       utils.CheckPointerValue(user.Email),
		Phone:       utils.CheckPointerValue(user.Phone),
		AccessToken: jwtToken,
	}
	return response, nil
}

func (s *userServiceImpl) Register(request model.RegisterRequest) error {
	if err := s.validator.RegisterValidate(request); err != nil {
		return err
	}

	checkEmail, err := s.userRepo.FindUserByPhoneOrEmail(request.Email)
	if err != nil {
		return err
	}

	if checkEmail != nil {
		return errors.New("email not available")
	}

	checkPhone, err := s.userRepo.FindUserByPhoneOrEmail(request.Phone)
	if err != nil {
		return err
	}

	if checkPhone != nil {
		return errors.New("phone not available")
	}

	// Hashing the password with the default cost
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	req := entity.User{
		ID:       uuid.NewString(),
		Name:     &request.Name,
		Email:    &request.Email,
		Phone:    &request.Phone,
		Password: string(hashedPassword),
		Default: entity.Default{
			CreatedAt: utils.JakartaTime(time.Now()),
			UpdatedAt: utils.JakartaTime(time.Now()),
		},
	}

	if err := s.userRepo.RegisterUser(req); err != nil {
		return err
	}

	return nil
}

func (s *userServiceImpl) generateJwt(user entity.User) (string, error) {
	jwtToken, err := s.JwtProvider.Generate(user.ID)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}
