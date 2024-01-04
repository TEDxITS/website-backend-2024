package service

import (
	"context"

	"github.com/TEDxITS/website-backend-2024/constants"
	"github.com/TEDxITS/website-backend-2024/dto"
	"github.com/TEDxITS/website-backend-2024/entity"
	"github.com/TEDxITS/website-backend-2024/repository"
	"github.com/TEDxITS/website-backend-2024/utils"
)

type (
	LinkShortenerService interface {
		RedirectByAlias(context.Context, string) (dto.LinkShortenerResponse, error)
		GetAllPagination(context.Context, dto.PaginationQuery) (dto.LinkShortenerPaginationResponse, error)
		CreateLinkShortener(context.Context, dto.LinkShortenerRequest) (dto.LinkShortenerResponse, error)
	}

	linkShortenerService struct {
		linkShortenRepo repository.LinkShortenerRepository
	}
)

func NewLinkShortenerService(repo repository.LinkShortenerRepository) LinkShortenerService {
	return &linkShortenerService{
		linkShortenRepo: repo,
	}
}

func (s *linkShortenerService) CreateLinkShortener(ctx context.Context, req dto.LinkShortenerRequest) (dto.LinkShortenerResponse, error) {
	aliasExist, err := s.linkShortenRepo.CheckAliasExist(req.Alias)
	if err != nil {
		return dto.LinkShortenerResponse{}, err
	}

	if aliasExist {
		return dto.LinkShortenerResponse{}, dto.ErrAliasHasBeenTaken
	}

	if !utils.ValidateLink(req.Link) {
		return dto.LinkShortenerResponse{}, dto.ErrLinkFormatInvalid
	}

	linkShorten := entity.LinkShortener{
		Alias: req.Alias,
		Link:  req.Link,
	}

	res, err := s.linkShortenRepo.CreateLinkShortener(linkShorten)
	if err != nil {
		return dto.LinkShortenerResponse{}, dto.ErrCreateLinkShortener
	}

	return dto.LinkShortenerResponse{
		Alias: res.Alias,
		Link:  res.Link,
	}, nil
}

func (s *linkShortenerService) RedirectByAlias(ctx context.Context, alias string) (dto.LinkShortenerResponse, error) {
	linkShorten, err := s.linkShortenRepo.GetLinkShortenerByAlias(alias)
	if err != nil {
		return dto.LinkShortenerResponse{}, dto.ErrAliasNotFound
	}

	return dto.LinkShortenerResponse{
		ID:    linkShorten.ID.String(),
		Alias: linkShorten.Alias,
		Link:  linkShorten.Link,
	}, nil
}

func (s *linkShortenerService) GetAllPagination(ctx context.Context, req dto.PaginationQuery) (dto.LinkShortenerPaginationResponse, error) {
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

	links, maxPage, count, err := s.linkShortenRepo.GetAllLinkShortenerPagination(req.Search, limit, page)
	if err != nil {
		return dto.LinkShortenerPaginationResponse{}, err
	}

	var result []dto.LinkShortenerResponse
	for _, link := range links {
		result = append(result, dto.LinkShortenerResponse{
			ID:    link.ID.String(),
			Alias: link.Alias,
			Link:  link.Link,
		})
	}

	return dto.LinkShortenerPaginationResponse{
		Data: result,
		PaginationMetadata: dto.PaginationMetadata{
			Page:    page,
			PerPage: limit,
			MaxPage: maxPage,
			Count:   count,
		},
	}, nil
}
