package service

import (
	"bytes"
	"context"
	"os"
	"text/template"
	"time"

	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/dto"
	"github.com/TEDxITS/website-backend-2024/entity"
	"github.com/TEDxITS/website-backend-2024/helpers"
	"github.com/TEDxITS/website-backend-2024/repository"
	"github.com/TEDxITS/website-backend-2024/utils"
)

type (
	UserService interface {
		RegisterUser(ctx context.Context, req dto.UserRequest) (dto.UserResponse, error)
		VerifyLogin(ctx context.Context, email string, password string) (entity.User, error)
		UpdateUser(ctx context.Context, req dto.UserRequest, userId string) (dto.UserResponse, error)
		Me(ctx context.Context, userId string, userRole string) (dto.UserResponse, error)
		GetAllPagination(ctx context.Context, req dto.PaginationQuery) (dto.UserPaginationResponse, error)

		SendVerificationEmail(userEmail string) error
	}

	userService struct {
		userRepo repository.UserRepository
	}
)

func NewUserService(ur repository.UserRepository) UserService {
	return &userService{
		userRepo: ur,
	}
}

func (s *userService) RegisterUser(ctx context.Context, req dto.UserRequest) (dto.UserResponse, error) {
	email, _ := s.userRepo.CheckEmailExist(req.Email)
	if email {
		return dto.UserResponse{}, dto.ErrEmailAlreadyExists
	}

	if !utils.ValidateEmail(req.Email) {
		return dto.UserResponse{}, dto.ErrEmailFormatInvalid
	}

	user := entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Verified: false,
	}

	userReg, err := s.userRepo.RegisterUser(user)
	if err != nil {
		return dto.UserResponse{}, dto.ErrCreateUser
	}

	return dto.UserResponse{
		ID:         userReg.ID.String(),
		Name:       userReg.Name,
		Role:       userReg.RoleID,
		Email:      userReg.Email,
		IsVerified: userReg.Verified,
	}, nil
}

func generateVerificationEmail(userEmail string) (utils.Email, error) {
	expired := time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05")
	token, err := utils.AESEncrypt(userEmail + "||" + expired)
	if err != nil {
		return utils.Email{}, err
	}

	verifyLink := constants.BASE_URL + "/api/user/verify/" + token
	readHtml, err := os.ReadFile("utils/template/base_mail.html")
	if err != nil {
		return utils.Email{}, err
	}

	data := struct {
		Email  string
		Verify string
	}{
		Email:  userEmail,
		Verify: verifyLink,
	}

	tmpl, err := template.New("custom").Parse(string(readHtml))
	if err != nil {
		return utils.Email{}, err
	}

	var strMail bytes.Buffer
	if err := tmpl.Execute(&strMail, data); err != nil {
		return utils.Email{}, err
	}

	return utils.Email{
		Email:   userEmail,
		Subject: "Verify Your TEDxITS Account",
		Body:    strMail.String(),
	}, nil
}

func (s *userService) SendVerificationEmail(userEmail string) error {
	email, err := generateVerificationEmail(userEmail)
	if err != nil {
		return dto.ErrVerifyEmailNotGenerated
	}

	err = utils.SendMail(email)
	if err != nil {
		return dto.ErrVerifyEmailNotSent
	}

	return nil
}

func (s *userService) VerifyLogin(ctx context.Context, email string, password string) (entity.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return entity.User{}, dto.ErrCredentialsNotMatched
	}

	if !user.Verified {
		return entity.User{}, dto.ErrAccountNotVerified
	}

	checkPassword, err := helpers.CheckPassword(user.Password, []byte(password))
	if err != nil || !checkPassword {
		return entity.User{}, dto.ErrCredentialsNotMatched
	}

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, req dto.UserRequest, userId string) (dto.UserResponse, error) {
	user, err := s.userRepo.GetUserById(userId)
	if err != nil {
		return dto.UserResponse{}, dto.ErrUserNotFound
	}

	data := entity.User{
		ID:       user.ID,
		Name:     req.Name,
		Role:     user.Role,
		Email:    req.Email,
		Password: req.Password,
	}

	userUpdate, err := s.userRepo.UpdateUser(data)
	if err != nil {
		return dto.UserResponse{}, dto.ErrUpdateUser
	}

	return dto.UserResponse{
		ID:         userUpdate.ID.String(),
		Name:       userUpdate.Name,
		Role:       ctx.Value("user_role").(string),
		Email:      userUpdate.Email,
		IsVerified: userUpdate.Verified,
	}, nil
}

func (s *userService) Me(ctx context.Context, userId string, userRole string) (dto.UserResponse, error) {
	user, err := s.userRepo.GetUserById(userId)
	if err != nil {
		return dto.UserResponse{}, dto.ErrGetUserById
	}

	return dto.UserResponse{
		ID:         user.ID.String(),
		Name:       user.Name,
		Role:       userRole,
		Email:      user.Email,
		IsVerified: user.Verified,
	}, nil
}

func (s *userService) GetAllPagination(ctx context.Context, req dto.PaginationQuery) (dto.UserPaginationResponse, error) {
	var limit int
	var page int

	limit = req.PerPage
	if limit <= 0 {
		limit = constants.ENUM_PAGINATION_LIMIT
	}

	page = req.Page
	if page <= 0 {
		page = constants.ENUM_PAGINATION_PAGE
	}

	users, maxPage, count, err := s.userRepo.GetAllUserPagination(req.Search, limit, page)
	if err != nil {
		return dto.UserPaginationResponse{}, err
	}

	var result []dto.UserResponse
	for _, user := range users {
		result = append(result, dto.UserResponse{
			ID:         user.ID.String(),
			Name:       user.Name,
			Email:      user.Email,
			RoleID:     user.RoleID,
			IsVerified: user.Verified,
		})
	}

	return dto.UserPaginationResponse{
		Data: result,
		PaginationMetadata: dto.PaginationMetadata{
			Page:    page,
			PerPage: limit,
			MaxPage: maxPage,
			Count:   count,
		},
	}, nil
}
