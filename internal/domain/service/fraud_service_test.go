package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/service"
	"github.com/aiagent/internal/domain/service/mocks"
	"github.com/aiagent/internal/domain/valueobject"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFraudDetectionService_GetUserRiskScore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFraudDetectionRepository(ctrl)
	mockNotif := mocks.NewMockNotificationService(ctrl)
	mockBatch := mocks.NewMockBatchJobService(ctrl)
	svc := service.NewFraudDetectionService(mockRepo, mockNotif, mockBatch)

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

	mockRepo.EXPECT().GetRiskScoreByUser(gomock.Any(), userID).Return(expectedScore, nil)
	mockRepo.EXPECT().GetBadgeStatusByUser(gomock.Any(), userID).Return(expectedBadge, nil)

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
}

func TestFraudDetectionService_GetUserRiskScore_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFraudDetectionRepository(ctrl)
	mockNotif := mocks.NewMockNotificationService(ctrl)
	mockBatch := mocks.NewMockBatchJobService(ctrl)
	svc := service.NewFraudDetectionService(mockRepo, mockNotif, mockBatch)

	userID := uuid.New()

	mockRepo.EXPECT().GetRiskScoreByUser(gomock.Any(), userID).Return(nil, nil)

	// Act
	result, err := svc.GetUserRiskScore(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestFraudDetectionService_GetFraudDashboard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFraudDetectionRepository(ctrl)
	mockNotif := mocks.NewMockNotificationService(ctrl)
	mockBatch := mocks.NewMockBatchJobService(ctrl)
	svc := service.NewFraudDetectionService(mockRepo, mockNotif, mockBatch)

	minScore := 70
	req := valueobject.FraudDashboardFilter{
		MinRiskScore: &minScore,
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

	mockRepo.EXPECT().GetUsersByRiskScoreRange(gomock.Any(), 70, 100, 1, 10).Return(riskScores, 2, nil)
	mockRepo.EXPECT().GetBotSignalsByUser(gomock.Any(), userID1, false).Return([]entity.BotDetectionSignal{}, nil)
	mockRepo.EXPECT().GetBotSignalsByUser(gomock.Any(), userID2, false).Return([]entity.BotDetectionSignal{}, nil)
	mockRepo.EXPECT().GetLastReviewByUser(gomock.Any(), userID1).Return(nil, nil)
	mockRepo.EXPECT().GetLastReviewByUser(gomock.Any(), userID2).Return(nil, nil)

	// Act
	result, err := svc.GetFraudDashboard(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.TotalCount)
	assert.Len(t, result.Users, 2)
}

func TestFraudDetectionService_ReviewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFraudDetectionRepository(ctrl)
	mockNotif := mocks.NewMockNotificationService(ctrl)
	mockBatch := mocks.NewMockBatchJobService(ctrl)
	svc := service.NewFraudDetectionService(mockRepo, mockNotif, mockBatch)

	adminID := uuid.New()
	userID := uuid.New()
	req := valueobject.ReviewUserCommand{
		Notes: "Reviewed user profile, looks suspicious",
	}

	riskScore := &entity.UserRiskScore{
		ID:           uuid.New(),
		UserID:       userID,
		OverallScore: 80,
	}

	mockRepo.EXPECT().GetRiskScoreByUser(gomock.Any(), userID).Return(riskScore, nil)
	mockRepo.EXPECT().CreateAdminReview(gomock.Any(), gomock.Any()).Return(nil)

	// Act
	result, err := svc.ReviewUser(context.Background(), adminID, userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, adminID, result.AdminID)
	assert.Equal(t, "reviewed", result.Action)
	assert.Equal(t, req.Notes, result.Notes)
}

func TestFraudDetectionService_BanUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFraudDetectionRepository(ctrl)
	mockNotif := mocks.NewMockNotificationService(ctrl)
	mockBatch := mocks.NewMockBatchJobService(ctrl)
	svc := service.NewFraudDetectionService(mockRepo, mockNotif, mockBatch)

	adminID := uuid.New()
	userID := uuid.New()
	req := valueobject.BanUserCommand{
		Reason: "Multiple bot detection signals confirmed",
		Notes:  "Banned after thorough investigation",
	}

	riskScore := &entity.UserRiskScore{
		ID:           uuid.New(),
		UserID:       userID,
		OverallScore: 90,
	}

	mockRepo.EXPECT().GetRiskScoreByUser(gomock.Any(), userID).Return(riskScore, nil)
	mockRepo.EXPECT().CreateAdminReview(gomock.Any(), gomock.Any()).Return(nil)

	// Act
	result, err := svc.BanUser(context.Background(), adminID, userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, adminID, result.AdminID)
	assert.Equal(t, "banned", result.Action)
	assert.Equal(t, req.Reason, result.Reason)
}

func TestFraudDetectionService_GetFraudTrends(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFraudDetectionRepository(ctrl)
	mockNotif := mocks.NewMockNotificationService(ctrl)
	mockBatch := mocks.NewMockBatchJobService(ctrl)
	svc := service.NewFraudDetectionService(mockRepo, mockNotif, mockBatch)

	req := valueobject.FraudTrendsFilter{
		Period: "7d",
	}

	now := time.Now()

	signalsByType := map[string]int{
		"rapid_follows": 15,
		"ip_cluster":    8,
	}

	dailyStats := []valueobject.DailyFraudStat{
		{Date: now.AddDate(0, 0, -1).Format("2006-01-02"), NewSignals: 5, NewSuspiciousAccounts: 2},
		{Date: now.AddDate(0, 0, -2).Format("2006-01-02"), NewSignals: 3, NewSuspiciousAccounts: 1},
	}

	mockRepo.EXPECT().GetSignalsCountByType(gomock.Any(), gomock.Any(), gomock.Any()).Return(signalsByType, nil)
	mockRepo.EXPECT().GetNewSuspiciousAccountsCount(gomock.Any(), gomock.Any(), gomock.Any()).Return(10, nil)
	mockRepo.EXPECT().GetBannedAccountsCount(gomock.Any(), gomock.Any(), gomock.Any()).Return(3, nil)
	mockRepo.EXPECT().GetReviewedAccountsCount(gomock.Any(), gomock.Any(), gomock.Any()).Return(8, nil)
	mockRepo.EXPECT().GetAverageRiskScore(gomock.Any(), gomock.Any(), gomock.Any()).Return(35.5, nil)
	mockRepo.EXPECT().GetRiskScoreDistribution(gomock.Any(), gomock.Any(), gomock.Any()).Return(map[string]int{"0-20": 100, "21-40": 50}, nil)
	mockRepo.EXPECT().GetDailyFraudStats(gomock.Any(), gomock.Any(), gomock.Any()).Return(dailyStats, nil)

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
}

func TestFraudDetectionService_TriggerBatchAnalysis(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockFraudDetectionRepository(ctrl)
	mockNotif := mocks.NewMockNotificationService(ctrl)
	mockBatch := mocks.NewMockBatchJobService(ctrl)
	svc := service.NewFraudDetectionService(mockRepo, mockNotif, mockBatch)

	req := valueobject.BatchAnalyzeCommand{}
	expectedJobID := uuid.New()

	mockBatch.EXPECT().StartBatchAnalysis(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedJobID, nil)

	// Act
	result, err := svc.TriggerBatchAnalysis(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "started", result.Status)
	assert.Equal(t, expectedJobID, result.JobID)
}
