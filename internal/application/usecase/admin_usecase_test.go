package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/aiagent/internal/application/usecase"
	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetDashboardStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockBlogRepo := mocks.NewMockBlogRepository(ctrl)
	mockCommentRepo := mocks.NewMockCommentRepository(ctrl)

	uc := usecase.NewAdminUseCase(mockUserRepo, mockBlogRepo, mockCommentRepo)

	// Mock Data
	now := time.Now()
	// Current month (index 11)
	currentMonth := now.Format("2006-01")
	// Last month (index 10)
	lastMonth := now.AddDate(0, -1, 0).Format("2006-01")

	userCounts := []entity.MonthlyCount{
		{Month: currentMonth, Count: 100},
		{Month: lastMonth, Count: 10},
	}
	blogCounts := []entity.MonthlyCount{
		{Month: lastMonth, Count: 5},
	}
	commentCounts := []entity.MonthlyCount{
		{Month: currentMonth, Count: 50},
	}

	mockUserRepo.EXPECT().CountByMonth(gomock.Any(), 12).Return(userCounts, nil)
	mockBlogRepo.EXPECT().CountByMonth(gomock.Any(), 12).Return(blogCounts, nil)
	mockCommentRepo.EXPECT().CountByMonth(gomock.Any(), 12).Return(commentCounts, nil)

	stats, err := uc.GetDashboardStats(context.Background())
	assert.NoError(t, err)
	assert.Len(t, stats.Stats, 12)

	// Verify Current Month (Newest, last in list)
	newestStat := stats.Stats[11]
	assert.Equal(t, currentMonth, newestStat.Month)
	assert.Equal(t, int64(100), newestStat.NewUsers)
	assert.Equal(t, int64(0), newestStat.NewBlogs)
	assert.Equal(t, int64(50), newestStat.NewComments)

	// Verify Last Month
	prevStat := stats.Stats[10]
	assert.Equal(t, lastMonth, prevStat.Month)
	assert.Equal(t, int64(10), prevStat.NewUsers)
	assert.Equal(t, int64(5), prevStat.NewBlogs)
	assert.Equal(t, int64(0), prevStat.NewComments)
}
