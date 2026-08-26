package application

import (
	"context"
	"time"

	"github.com/google/uuid"

	"palomnik/internal/domain"
	"palomnik/internal/ports"
)

type NewsLikeState struct {
	Count      int
	LikedByYou bool
}

type NewsEngagementService struct {
	engagement ports.NewsEngagementRepository
	news       ports.NewsRepository
	users      ports.UserRepository
}

func NewNewsEngagementService(
	engagement ports.NewsEngagementRepository,
	news ports.NewsRepository,
	users ports.UserRepository,
) *NewsEngagementService {
	return &NewsEngagementService{
		engagement: engagement,
		news:       news,
		users:      users,
	}
}

func (s *NewsEngagementService) GetLikeState(ctx context.Context, slug, visitorRaw string) (NewsLikeState, error) {
	article, err := s.publishedNewsBySlug(ctx, slug)
	if err != nil {
		return NewsLikeState{}, err
	}

	count, err := s.engagement.CountNewsLikes(ctx, article.ID)
	if err != nil {
		return NewsLikeState{}, err
	}

	visitorID, err := domain.NormalizeVisitorID(visitorRaw)
	if err != nil || visitorID == "" {
		return NewsLikeState{Count: count, LikedByYou: false}, nil
	}

	liked, err := s.engagement.HasNewsLike(ctx, article.ID, visitorID)
	if err != nil {
		return NewsLikeState{}, err
	}

	return NewsLikeState{Count: count, LikedByYou: liked}, nil
}

func (s *NewsEngagementService) ToggleLike(ctx context.Context, slug, visitorRaw string) (NewsLikeState, error) {
	article, err := s.publishedNewsBySlug(ctx, slug)
	if err != nil {
		return NewsLikeState{}, err
	}

	visitorID, err := domain.NormalizeVisitorID(visitorRaw)
	if err != nil {
		return NewsLikeState{}, err
	}

	liked, err := s.engagement.HasNewsLike(ctx, article.ID, visitorID)
	if err != nil {
		return NewsLikeState{}, err
	}

	if liked {
		if err := s.engagement.RemoveNewsLike(ctx, article.ID, visitorID); err != nil {
			return NewsLikeState{}, err
		}
	} else if err := s.engagement.AddNewsLike(ctx, article.ID, visitorID); err != nil {
		return NewsLikeState{}, err
	}

	return s.GetLikeState(ctx, slug, visitorID)
}

func (s *NewsEngagementService) ListComments(ctx context.Context, slug string, pagination ports.Pagination) (ports.NewsCommentList, error) {
	article, err := s.publishedNewsBySlug(ctx, slug)
	if err != nil {
		return ports.NewsCommentList{}, err
	}

	pagination = ports.NormalizePagination(pagination.Page, pagination.Limit)
	total, err := s.engagement.CountNewsComments(ctx, article.ID)
	if err != nil {
		return ports.NewsCommentList{}, err
	}

	items, err := s.engagement.ListNewsComments(ctx, article.ID, pagination)
	if err != nil {
		return ports.NewsCommentList{}, err
	}

	start := (pagination.Page - 1) * pagination.Limit
	hasNext := start+len(items) < total

	return ports.NewsCommentList{
		Items: items,
		Meta: ports.PageMeta{
			Page:    pagination.Page,
			Limit:   pagination.Limit,
			Total:   total,
			HasNext: hasNext,
		},
	}, nil
}

func (s *NewsEngagementService) AddComment(ctx context.Context, slug string, userID uuid.UUID, bodyRaw string) (domain.NewsComment, error) {
	article, err := s.publishedNewsBySlug(ctx, slug)
	if err != nil {
		return domain.NewsComment{}, err
	}

	user, err := s.users.GetUserByID(ctx, userID)
	if err != nil {
		return domain.NewsComment{}, err
	}

	comment, err := domain.NewNewsComment(domain.NewNewsCommentInput{
		ID:     uuid.New(),
		NewsID: article.ID,
		UserID: userID,
		Body:   bodyRaw,
		Now:    time.Now().UTC(),
	})
	if err != nil {
		return domain.NewsComment{}, err
	}
	comment.Author = domain.DisplayNewsCommentAuthor(user.Name)

	return s.engagement.AddNewsComment(ctx, comment)
}

func (s *NewsEngagementService) publishedNewsBySlug(ctx context.Context, slug string) (domain.NewsArticle, error) {
	article, err := s.news.GetNewsBySlug(ctx, slug)
	if err != nil {
		return domain.NewsArticle{}, err
	}
	if !article.IsPublished {
		return domain.NewsArticle{}, domain.ErrNotFound
	}
	return article, nil
}
