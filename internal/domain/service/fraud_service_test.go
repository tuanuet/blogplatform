package service

import (
	"context"
	"testing"
	"time"

	"github.com/aiagent/boilerplate/internal/domain/entity"
	"github.com/aiagent/boilerplate/internal/interfaces/http/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFraudDetectionRepository is a mock implementation of FraudDetectionRepository
type MockFraudDetectionRepository struct {
	mock.Mock
}

func (m *MockFraudDetectionRepository) CreateFollowerEvent(ctx context.Context, event *entity.FollowerEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockFraudDetectionRepository) GetFollowerEventsByUser(ctx context.Context, userID uuid.UUID, from, to *time.Time) ([]entity.FollowerEvent, error) {
	args := m.Called(ctx, userID, from, to)
	return args.Get(0).([]entity.FollowerEvent), args.Error(1)
}

func (m *MockFraudDetectionRepository) GetFollowerEventsByIP(ctx context.Context, ipAddress string, from, to *time.Time) ([]entity.FollowerEvent, error) {
	args := m.Called(ctx, ipAddress, from, to)
	return args.Get(0).([]entity.FollowerEvent), args.Error(1)
}

func (m *MockFraudDetectionRepository) CreateBotSignal(ctx context.Context, signal *entity.BotDetectionSignal) error {
	args := m.Called(ctx, signal)
	return args.Error(0)
}

func (m *MockFraudDetectionRepository) GetBotSignalsByUser(ctx context.Context, userID uuid.UUID, processed bool) ([]entity.BotDetectionSignal, error) {
	args := m.Called(ctx, userID, processed)
	return args.Get(0).([]entity.BotDetectionSignal), args.Error(1)
}

func (m *MockFraudDetectionRepository) GetUnprocessedBotSignals(ctx context.Context, limit int) ([]entity.BotDetectionSignal, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]entity.BotDetectionSignal), args.Error(1)
}

func (m *MockFraudDetectionRepository) MarkBotSignalAsProcessed(ctx context.Context, signalID uuid.UUID) error {
	args := m.Called(ctx, signalID)
	return args.Error(0)
}

func (m *MockFraudDetectionRepository) CreateOrUpdateRiskScore(ctx context.Context, score *entity.UserRiskScore) error {
	args := m.Called(ctx, score)
	return args.Error(0)
}

func (m *MockFraudDetectionRepository) GetRiskScoreByUser(ctx context.Context, userID uuid.UUID) (*entity.UserRiskScore, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserRiskScore), args.Error(1)
}

func (m *MockFraudDetectionRepository) GetUsersByRiskScoreRange(ctx context.Context, minScore, maxScore int, page, pageSize int) ([]entity.UserRiskScore, int, error) {
	args := m.Called(ctx, minScore, maxScore, page, pageSize)
	return args.Get(0).([]entity.UserRiskScore), args.Int(1), args.Error(2)
}

func (m *MockFraudDetectionRepository) CreateOrUpdateBadgeStatus(ctx context.Context, status *entity.UserBadgeStatus) error {
	args := m.Called(ctx, status)
	return args.Error(0)
}

func (m *MockFraudDetectionRepository) GetBadgeStatusByUser(ctx context.Context, userID uuid.UUID) (*entity.UserBadgeStatus, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserBadgeStatus), args.Error(1)
}

func (m *MockFraudDetectionRepository) CreateAdminReview(ctx context.Context, review *entity.AdminReview) error {
	args := m.Called(ctx, review)
	return args.Error(0)
}

func (m *MockFraudDetectionRepository) GetAdminReviewsByUser(ctx context.Context, userID uuid.UUID, limit int) ([]entity.AdminReview, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]entity.AdminReview), args.Error(1)
}

func (m *MockFraudDetectionRepository) GetLastReviewByUser(ctx context.Context, userID uuid.UUID) (*entity.AdminReview, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AdminReview), args.Error(1)
}

func (m *MockFraudDetectionRepository) CreateBotNotification(ctx context.Context, notification *entity.BotFollowerNotification) error {
	args := m.Called(ctx, notification)
	return args.Error(0)
}

func (m *MockFraudDetectionRepository) GetBotNotificationsByUser(ctx context.Context, userID uuid.UUID, unreadOnly bool) ([]entity.BotFollowerNotification, error) {
	args := m.Called(ctx, userID, unreadOnly)
	return args.Get(0).([]entity.BotFollowerNotification), args.Error(1)
}

func (m *MockFraudDetectionRepository) MarkNotificationAsRead(ctx context.Context, notificationID uuid.UUID) error {
	args := m.Called(ctx, notificationID)
	return args.Error(0)
}

func (m *MockFraudDetectionRepository) GetSignalsCountByType(ctx context.Context, from, to time.Time) (map[string]int, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *MockFraudDetectionRepository) GetNewSuspiciousAccountsCount(ctx context.Context, from, to time.Time) (int, error) {
	args := m.Called(ctx, from, to)
	return args.Int(0), args.Error(1)
}

func (m *MockFraudDetectionRepository) GetBannedAccountsCount(ctx context.Context, from, to time.Time) (int, error) {
	args := m.Called(ctx, from, to)
	return args.Int(0), args.Error(1)
}

func (m *MockFraudDetectionRepository) GetReviewedAccountsCount(ctx context.Context, from, to time.Time) (int, error) {
	args := m.Called(ctx, from, to)
	return args.Int(0), args.Error(1)
}

func (m *MockFraudDetectionRepository) GetAverageRiskScore(ctx context.Context, from, to time.Time) (float64, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockFraudDetectionRepository) GetRiskScoreDistribution(ctx context.Context, from, to time.Time) (map[string]int, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *MockFraudDetectionRepository) GetDailyFraudStats(ctx context.Context, from, to time.Time) ([]dto.DailyFraudStat, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).([]dto.DailyFraudStat), args.Error(1)
}

// MockNotificationService is a mock implementation of NotificationService
type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) SendBotFollowerNotification(ctx context.Context, userID uuid.UUID, notifications []dto.BotFollowerNotificationResponse) error {
	args := m.Called(ctx, userID, notifications)
	return args.Error(0)
}

func (m *MockNotificationService) SendBadgeStatusUpdate(ctx context.Context, userID uuid.UUID, status string, reason string) error {
	args := m.Called(ctx, userID, status, reason)
	return args.Error(0)
}

// MockBatchJobService is a mock implementation of BatchJobService
type MockBatchJobService struct {
	mock.Mock
}

func (m *MockBatchJobService) StartBatchAnalysis(ctx context.Context, dateFrom, dateTo *time.Time) (uuid.UUID, error) {
	args := m.Called(ctx, dateFrom, dateTo)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockBatchJobService) GetBatchJobStatus(ctx context.Context, jobID uuid.UUID) (*dto.BatchAnalyzeResponse, error) {
	args := m.Called(ctx, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.BatchAnalyzeResponse), args.Error(1)
}

func TestFraudDetectionService_GetUserRiskScore(t *testing.T) {
	// Arrange
	mockRepo := new(MockFraudDetectionRepository)
	mockNotif := new(MockNotificationService)
	mockBatch := new(MockBatchJobService)
	svc := NewFraudDetectionService(mockRepo, mockNotif, mockBatch)

	userID := uuid.New()
	expectedScore := &entity.UserRiskScore{
		ID:                        uuid.New(),
		UserID:                    userID,
		OverallScore:              45,
		FollowerAuthenticityScore: 60,
		EngagementQualityScore:    70,
		AccountAgeFactor:          0.8,
		CalculationVersion:        "v1.0",
		LastCalculatedAt:          time.Now(),
	}

	expectedBadge := &entity.UserBadgeStatus{
		ID:     uuid.New(),
		UserID: userID,
		Status: "active",
	}

	mockRepo.On("GetRiskScoreByUser", mock.Anything, userID).Return(expectedScore, nil)
	mockRepo.On("GetBadgeStatusByUser", mock.Anything, userID).Return(expectedBadge, nil)

	// Act
	result, err := svc.GetUserRiskScore(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, 45, result.OverallScore)
	assert.Equal(t, 60, result.FollowerAuthenticityScore)
	assert.Equal(t, 70, result.EngagementQualityScore)
	assert.Equal(t, "active", result.BadgeStatus)
	mockRepo.AssertExpectations(t)
}

func TestFraudDetectionService_GetUserRiskScore_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockFraudDetectionRepository)
	mockNotif := new(MockNotificationService)
	mockBatch := new(MockBatchJobService)
	svc := NewFraudDetectionService(mockRepo, mockNotif, mockBatch)

	userID := uuid.New()

	mockRepo.On("GetRiskScoreByUser", mock.Anything, userID).Return(nil, nil)

	// Act
	result, err := svc.GetUserRiskScore(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestFraudDetectionService_GetFraudDashboard(t *testing.T) {
	// Arrange
	mockRepo := new(MockFraudDetectionRepository)
	mockNotif := new(MockNotificationService)
	mockBatch := new(MockBatchJobService)
	svc := NewFraudDetectionService(mockRepo, mockNotif, mockBatch)

	req := dto.FraudDashboardRequest{
		MinRiskScore: 70,
		Page:         1,
		PageSize:     10,
	}

	userID1 := uuid.New()
	userID2 := uuid.New()

	riskScores := []entity.UserRiskScore{
		{
			ID:               uuid.New(),
			UserID:           userID1,
			OverallScore:     85,
			LastCalculatedAt: time.Now(),
		},
		{
			ID:               uuid.New(),
			UserID:           userID2,
			OverallScore:     75,
			LastCalculatedAt: time.Now(),
		},
	}

	mockRepo.On("GetUsersByRiskScoreRange", mock.Anything, 70, 100, 1, 10).Return(riskScores, 2, nil)
	mockRepo.On("GetBotSignalsByUser", mock.Anything, userID1, false).Return([]entity.BotDetectionSignal{}, nil)
	mockRepo.On("GetBotSignalsByUser", mock.Anything, userID2, false).Return([]entity.BotDetectionSignal{}, nil)
	mockRepo.On("GetLastReviewByUser", mock.Anything, userID1).Return(nil, nil)
	mockRepo.On("GetLastReviewByUser", mock.Anything, userID2).Return(nil, nil)

	// Act
	result, err := svc.GetFraudDashboard(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.TotalCount)
	assert.Len(t, result.Users, 2)
	mockRepo.AssertExpectations(t)
}

func TestFraudDetectionService_ReviewUser(t *testing.T) {
	// Arrange
	mockRepo := new(MockFraudDetectionRepository)
	mockNotif := new(MockNotificationService)
	mockBatch := new(MockBatchJobService)
	svc := NewFraudDetectionService(mockRepo, mockNotif, mockBatch)

	adminID := uuid.New()
	userID := uuid.New()
	req := dto.ReviewUserRequest{
		Notes: "Reviewed user profile, looks suspicious",
	}

	riskScore := &entity.UserRiskScore{
		ID:           uuid.New(),
		UserID:       userID,
		OverallScore: 80,
	}

	mockRepo.On("GetRiskScoreByUser", mock.Anything, userID).Return(riskScore, nil)
	mockRepo.On("CreateAdminReview", mock.Anything, mock.AnythingOfType("*entity.AdminReview")).Return(nil)

	// Act
	result, err := svc.ReviewUser(context.Background(), adminID, userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, adminID, result.AdminID)
	assert.Equal(t, "reviewed", result.Action)
	assert.Equal(t, req.Notes, result.Notes)
	mockRepo.AssertExpectations(t)
}

func TestFraudDetectionService_BanUser(t *testing.T) {
	// Arrange
	mockRepo := new(MockFraudDetectionRepository)
	mockNotif := new(MockNotificationService)
	mockBatch := new(MockBatchJobService)
	svc := NewFraudDetectionService(mockRepo, mockNotif, mockBatch)

	adminID := uuid.New()
	userID := uuid.New()
	req := dto.BanUserRequest{
		Reason: "Multiple bot detection signals confirmed",
		Notes:  "Banned after thorough investigation",
	}

	riskScore := &entity.UserRiskScore{
		ID:           uuid.New(),
		UserID:       userID,
		OverallScore: 90,
	}

	mockRepo.On("GetRiskScoreByUser", mock.Anything, userID).Return(riskScore, nil)
	mockRepo.On("CreateAdminReview", mock.Anything, mock.AnythingOfType("*entity.AdminReview")).Return(nil)

	// Act
	result, err := svc.BanUser(context.Background(), adminID, userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, adminID, result.AdminID)
	assert.Equal(t, "banned", result.Action)
	assert.Equal(t, req.Reason, result.Reason)
	mockRepo.AssertExpectations(t)
}

func TestFraudDetectionService_GetFraudTrends(t *testing.T) {
	// Arrange
	mockRepo := new(MockFraudDetectionRepository)
	mockNotif := new(MockNotificationService)
	mockBatch := new(MockBatchJobService)
	svc := NewFraudDetectionService(mockRepo, mockNotif, mockBatch)

	req := dto.FraudTrendsRequest{
		Period: "7d",
	}

	now := time.Now()

	signalsByType := map[string]int{
		"rapid_follows": 15,
		"ip_cluster":    8,
	}

	dailyStats := []dto.DailyFraudStat{
		{Date: now.AddDate(0, 0, -1).Format("2006-01-02"), NewSignals: 5, NewSuspiciousAccounts: 2},
		{Date: now.AddDate(0, 0, -2).Format("2006-01-02"), NewSignals: 3, NewSuspiciousAccounts: 1},
	}

	mockRepo.On("GetSignalsCountByType", mock.Anything, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(signalsByType, nil)
	mockRepo.On("GetNewSuspiciousAccountsCount", mock.Anything, mock.Anything, mock.Anything).Return(10, nil)
	mockRepo.On("GetBannedAccountsCount", mock.Anything, mock.Anything, mock.Anything).Return(3, nil)
	mockRepo.On("GetReviewedAccountsCount", mock.Anything, mock.Anything, mock.Anything).Return(8, nil)
	mockRepo.On("GetAverageRiskScore", mock.Anything, mock.Anything, mock.Anything).Return(35.5, nil)
	mockRepo.On("GetRiskScoreDistribution", mock.Anything, mock.Anything, mock.Anything).Return(map[string]int{"0-20": 100, "21-40": 50}, nil)
	mockRepo.On("GetDailyFraudStats", mock.Anything, mock.Anything, mock.Anything).Return(dailyStats, nil)

	// Act
	result, err := svc.GetFraudTrends(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "7d", result.Period)
	assert.Equal(t, 23, result.TotalBotSignals) // 15 + 8
	assert.Equal(t, 10, result.NewSuspiciousAccounts)
	assert.Equal(t, 3, result.BannedAccounts)
	assert.Equal(t, 35.5, result.AverageRiskScore)
	mockRepo.AssertExpectations(t)
}

func TestFraudDetectionService_TriggerBatchAnalysis(t *testing.T) {
	// Arrange
	mockRepo := new(MockFraudDetectionRepository)
	mockNotif := new(MockNotificationService)
	mockBatch := new(MockBatchJobService)
	svc := NewFraudDetectionService(mockRepo, mockNotif, mockBatch)

	req := dto.BatchAnalyzeRequest{}
	expectedJobID := uuid.New()

	mockBatch.On("StartBatchAnalysis", mock.Anything, mock.Anything, mock.Anything).Return(expectedJobID, nil)

	// Act
	result, err := svc.TriggerBatchAnalysis(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "started", result.Status)
	assert.Equal(t, expectedJobID, result.JobID)
	mockBatch.AssertExpectations(t)
}
