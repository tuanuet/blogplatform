package service_test

import (
	"context"
	"testing"

	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/aiagent/boilerplate/internal/domain/repository"
	"github.com/aiagent/boilerplate/internal/domain/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
type MockBlogRepo struct {
	mock.Mock
}

func (m *MockBlogRepo) Create(ctx context.Context, blog *entity.Blog) error {
	return m.Called(ctx, blog).Error(0)
}
func (m *MockBlogRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Blog, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Blog), args.Error(1)
}
func (m *MockBlogRepo) FindBySlug(ctx context.Context, authorID uuid.UUID, slug string) (*entity.Blog, error) {
	args := m.Called(ctx, authorID, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Blog), args.Error(1)
}
func (m *MockBlogRepo) FindAll(ctx context.Context, filter repository.BlogFilter, pagination repository.Pagination) (*repository.PaginatedResult[entity.Blog], error) {
	args := m.Called(ctx, filter, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[entity.Blog]), args.Error(1)
}
func (m *MockBlogRepo) Update(ctx context.Context, blog *entity.Blog) error {
	return m.Called(ctx, blog).Error(0)
}
func (m *MockBlogRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockBlogRepo) AddTags(ctx context.Context, blogID uuid.UUID, tagIDs []uuid.UUID) error {
	return m.Called(ctx, blogID, tagIDs).Error(0)
}
func (m *MockBlogRepo) RemoveTags(ctx context.Context, blogID uuid.UUID, tagIDs []uuid.UUID) error {
	return m.Called(ctx, blogID, tagIDs).Error(0)
}
func (m *MockBlogRepo) ReplaceTags(ctx context.Context, blogID uuid.UUID, tagIDs []uuid.UUID) error {
	return m.Called(ctx, blogID, tagIDs).Error(0)
}
func (m *MockBlogRepo) React(ctx context.Context, blogID, userID uuid.UUID, reactionType entity.ReactionType) (int, int, error) {
	args := m.Called(ctx, blogID, userID, reactionType)
	return args.Int(0), args.Int(1), args.Error(2)
}
func (m *MockBlogRepo) UpdateCounts(ctx context.Context, blogID uuid.UUID, upDelta, downDelta int) error {
	return m.Called(ctx, blogID, upDelta, downDelta).Error(0)
}
func (m *MockBlogRepo) FindRelated(ctx context.Context, blogID uuid.UUID, limit int) ([]entity.Blog, error) {
	args := m.Called(ctx, blogID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Blog), args.Error(1)
}

type MockTagRepo struct {
	mock.Mock
}

func (m *MockTagRepo) Create(ctx context.Context, tag *entity.Tag) error {
	return m.Called(ctx, tag).Error(0)
}
func (m *MockTagRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Tag, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Tag), args.Error(1)
}
func (m *MockTagRepo) FindBySlug(ctx context.Context, slug string) (*entity.Tag, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Tag), args.Error(1)
}
func (m *MockTagRepo) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.Tag, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Tag), args.Error(1)
}
func (m *MockTagRepo) FindAll(ctx context.Context, pagination repository.Pagination) (*repository.PaginatedResult[entity.Tag], error) {
	args := m.Called(ctx, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.PaginatedResult[entity.Tag]), args.Error(1)
}
func (m *MockTagRepo) FindOrCreate(ctx context.Context, name, slug string) (*entity.Tag, error) {
	args := m.Called(ctx, name, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Tag), args.Error(1)
}
func (m *MockTagRepo) Update(ctx context.Context, tag *entity.Tag) error {
	return m.Called(ctx, tag).Error(0)
}
func (m *MockTagRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockTagRepo) FindPopular(ctx context.Context, limit int) ([]entity.Tag, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Tag), args.Error(1)
}

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}
func (m *MockUserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}
func (m *MockUserRepo) Update(ctx context.Context, user *entity.User) error {
	return m.Called(ctx, user).Error(0)
}
func (m *MockUserRepo) UpdateProfile(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error {
	return m.Called(ctx, userID, updates).Error(0)
}
func (m *MockUserRepo) GetInterests(ctx context.Context, userID uuid.UUID) ([]entity.Tag, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Tag), args.Error(1)
}
func (m *MockUserRepo) ReplaceInterests(ctx context.Context, userID uuid.UUID, tagIDs []uuid.UUID) error {
	return m.Called(ctx, userID, tagIDs).Error(0)
}

func TestGetPersonalizedFeed_WithInterests(t *testing.T) {
	mockBlogRepo := new(MockBlogRepo)
	mockTagRepo := new(MockTagRepo)
	mockUserRepo := new(MockUserRepo)

	svc := service.NewRecommendationService(mockBlogRepo, mockTagRepo, mockUserRepo)

	userID := uuid.New()
	tagID := uuid.New()

	// Expectations
	mockUserRepo.On("GetInterests", mock.Anything, userID).Return([]entity.Tag{
		{ID: tagID, Name: "Go", Slug: "go"},
	}, nil)

	mockBlogRepo.On("FindAll", mock.Anything, mock.MatchedBy(func(f repository.BlogFilter) bool {
		return len(f.TagIDs) == 1 && f.TagIDs[0] == tagID
	}), mock.Anything).Return(&repository.PaginatedResult[entity.Blog]{
		Data:  []entity.Blog{{Title: "Go Blog"}},
		Total: 1,
	}, nil)

	// Act
	result, err := svc.GetPersonalizedFeed(context.Background(), userID, repository.Pagination{Page: 1, PageSize: 10})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Go Blog", result.Data[0].Title)
	mockUserRepo.AssertExpectations(t)
	mockBlogRepo.AssertExpectations(t)
}

func TestGetPersonalizedFeed_NoInterests_Fallback(t *testing.T) {
	mockBlogRepo := new(MockBlogRepo)
	mockTagRepo := new(MockTagRepo)
	mockUserRepo := new(MockUserRepo)

	svc := service.NewRecommendationService(mockBlogRepo, mockTagRepo, mockUserRepo)

	userID := uuid.New()

	// Expectations
	mockUserRepo.On("GetInterests", mock.Anything, userID).Return([]entity.Tag{}, nil)

	mockBlogRepo.On("FindAll", mock.Anything, mock.MatchedBy(func(f repository.BlogFilter) bool {
		return len(f.TagIDs) == 0
	}), mock.Anything).Return(&repository.PaginatedResult[entity.Blog]{
		Data:  []entity.Blog{{Title: "Recent Blog"}},
		Total: 1,
	}, nil)

	// Act
	result, err := svc.GetPersonalizedFeed(context.Background(), userID, repository.Pagination{Page: 1, PageSize: 10})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Recent Blog", result.Data[0].Title)
}

func TestGetRelatedBlogs(t *testing.T) {
	mockBlogRepo := new(MockBlogRepo)
	mockTagRepo := new(MockTagRepo)
	mockUserRepo := new(MockUserRepo)

	svc := service.NewRecommendationService(mockBlogRepo, mockTagRepo, mockUserRepo)

	blogID := uuid.New()

	mockBlogRepo.On("FindRelated", mock.Anything, blogID, 3).Return([]entity.Blog{
		{Title: "Related 1"},
	}, nil)

	blogs, err := svc.GetRelatedBlogs(context.Background(), blogID, 3)

	assert.NoError(t, err)
	assert.Len(t, blogs, 1)
	assert.Equal(t, "Related 1", blogs[0].Title)
}

func TestUpdateInterests(t *testing.T) {
	mockBlogRepo := new(MockBlogRepo)
	mockTagRepo := new(MockTagRepo)
	mockUserRepo := new(MockUserRepo)

	svc := service.NewRecommendationService(mockBlogRepo, mockTagRepo, mockUserRepo)

	userID := uuid.New()
	tagID := uuid.New()
	tagIDs := []uuid.UUID{tagID}

	mockTagRepo.On("FindByIDs", mock.Anything, tagIDs).Return([]entity.Tag{{ID: tagID}}, nil)
	mockUserRepo.On("ReplaceInterests", mock.Anything, userID, tagIDs).Return(nil)

	err := svc.UpdateInterests(context.Background(), userID, tagIDs)

	assert.NoError(t, err)
	mockTagRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}
