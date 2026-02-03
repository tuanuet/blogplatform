package service_test
import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/repository/mocks"
	"github.com/aiagent/internal/domain/service"
	servicemocks "github.com/aiagent/internal/domain/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNotificationDispatcher_Notify_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := uuid.New()
	notifType := entity.NotificationTypeNewFollower
	targetID := uuid.New()
	actorID := uuid.New()
	
	data := map[string]interface{}{
		"target_id":  targetID.String(),
		"actor_id":   actorID.String(),
		"actor_name": "Test User",
	}

	mockNotifRepo := mocks.NewMockNotificationRepository(ctrl)
	mockPrefRepo := mocks.NewMockNotificationPreferenceRepository(ctrl)
	mockTokenRepo := mocks.NewMockDeviceTokenRepository(ctrl)
	mockAggregator := servicemocks.NewMockNotificationAggregator(ctrl)
	mockFirebase := servicemocks.NewMockFirebaseAdapter(ctrl)

	mockPrefRepo.EXPECT().IsEnabled(ctx, userID, notifType, "push").Return(true, nil)
	mockAggregator.EXPECT().CheckRateLimit(ctx, userID, notifType).Return(true, nil)
	mockAggregator.EXPECT().ShouldAggregate(ctx, userID, notifType, targetID).Return(nil, nil)
	mockNotifRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
	mockAggregator.EXPECT().IncrementRateLimit(ctx, userID, notifType).Return(nil)

	tokens := []*entity.UserDeviceToken{{DeviceToken: "token1", Platform: "ios"}}
	mockTokenRepo.EXPECT().FindByUserID(ctx, userID).Return(tokens, nil)
	mockFirebase.EXPECT().SendPushToUser(ctx, userID, gomock.Any(), gomock.Any(), data).Return(nil)

	dispatcher := service.NewNotificationDispatcher(mockNotifRepo, mockPrefRepo, mockTokenRepo, mockAggregator, mockFirebase)
	assert.NoError(t, dispatcher.Notify(ctx, userID, notifType, data))
}

func TestNotificationDispatcher_Notify_PreferenceDisabled_EarlyReturn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := uuid.New()
	notifType := entity.NotificationTypeNewFollower
	data := map[string]interface{}{}

	mockNotifRepo := mocks.NewMockNotificationRepository(ctrl)
	mockPrefRepo := mocks.NewMockNotificationPreferenceRepository(ctrl)
	mockTokenRepo := mocks.NewMockDeviceTokenRepository(ctrl)
	mockAggregator := servicemocks.NewMockNotificationAggregator(ctrl)
	mockFirebase := servicemocks.NewMockFirebaseAdapter(ctrl)

	mockPrefRepo.EXPECT().IsEnabled(ctx, userID, notifType, "push").Return(false, nil)
	dispatcher := service.NewNotificationDispatcher(mockNotifRepo, mockPrefRepo, mockTokenRepo, mockAggregator, mockFirebase)
	assert.NoError(t, dispatcher.Notify(ctx, userID, notifType, data))
}

func TestNotificationDispatcher_Notify_RateLimitExceeded_SkipsNotification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := uuid.New()
	notifType := entity.NotificationTypeNewFollower
	targetID := uuid.New()
	
	data := map[string]interface{}{"target_id": targetID.String()}

	mockNotifRepo := mocks.NewMockNotificationRepository(ctrl)
	mockPrefRepo := mocks.NewMockNotificationPreferenceRepository(ctrl)
	mockTokenRepo := mocks.NewMockDeviceTokenRepository(ctrl)
	mockAggregator := servicemocks.NewMockNotificationAggregator(ctrl)
	mockFirebase := servicemocks.NewMockFirebaseAdapter(ctrl)

	mockPrefRepo.EXPECT().IsEnabled(ctx, userID, notifType, "push").Return(true, nil)
	mockAggregator.EXPECT().CheckRateLimit(ctx, userID, notifType).Return(false, nil)

	dispatcher := service.NewNotificationDispatcher(mockNotifRepo, mockPrefRepo, mockTokenRepo, mockAggregator, mockFirebase)
	assert.NoError(t, dispatcher.Notify(ctx, userID, notifType, data))
}

func TestNotificationDispatcher_Notify_AggregationUpdatesExisting(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := uuid.New()
	notifType := entity.NotificationTypeBlogLike
	targetID := uuid.New()
	actorID := uuid.New()
	
	data := map[string]interface{}{
		"target_id":  targetID.String(),
		"actor_id":   actorID.String(),
		"actor_name": "Test User",
	}

	mockNotifRepo := mocks.NewMockNotificationRepository(ctrl)
	mockPrefRepo := mocks.NewMockNotificationPreferenceRepository(ctrl)
	mockTokenRepo := mocks.NewMockDeviceTokenRepository(ctrl)
	mockAggregator := servicemocks.NewMockNotificationAggregator(ctrl)
	mockFirebase := servicemocks.NewMockFirebaseAdapter(ctrl)

	existingNotif := &entity.Notification{
		ID:           uuid.New(),
		UserID:       userID,
		Type:         notifType,
		Body:         "Test body",
		Data:         data,
		IsRead:       false,
		GroupedCount: 1,
	}

	mockPrefRepo.EXPECT().IsEnabled(ctx, userID, notifType, "push").Return(true, nil)
	mockAggregator.EXPECT().CheckRateLimit(ctx, userID, notifType).Return(true, nil)
	mockAggregator.EXPECT().ShouldAggregate(ctx, userID, notifType, targetID).Return(existingNotif, nil)
	mockNotifRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
	mockAggregator.EXPECT().IncrementRateLimit(ctx, userID, notifType).Return(nil)

	tokens := []*entity.UserDeviceToken{{DeviceToken: "token1", Platform: "ios"}}
	mockTokenRepo.EXPECT().FindByUserID(ctx, userID).Return(tokens, nil)
	mockFirebase.EXPECT().SendPushToUser(ctx, userID, gomock.Any(), gomock.Any(), data).Return(nil)

	dispatcher := service.NewNotificationDispatcher(mockNotifRepo, mockPrefRepo, mockTokenRepo, mockAggregator, mockFirebase)
	assert.NoError(t, dispatcher.Notify(ctx, userID, notifType, data))
}

func TestNotificationDispatcher_Notify_NoDeviceTokens_NoPush(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := uuid.New()
	notifType := entity.NotificationTypeNewFollower
	targetID := uuid.New()
	
	data := map[string]interface{}{"target_id": targetID.String()}

	mockNotifRepo := mocks.NewMockNotificationRepository(ctrl)
	mockPrefRepo := mocks.NewMockNotificationPreferenceRepository(ctrl)
	mockTokenRepo := mocks.NewMockDeviceTokenRepository(ctrl)
	mockAggregator := servicemocks.NewMockNotificationAggregator(ctrl)
	mockFirebase := servicemocks.NewMockFirebaseAdapter(ctrl)

	mockPrefRepo.EXPECT().IsEnabled(ctx, userID, notifType, "push").Return(true, nil)
	mockAggregator.EXPECT().CheckRateLimit(ctx, userID, notifType).Return(true, nil)
	mockAggregator.EXPECT().ShouldAggregate(ctx, userID, notifType, targetID).Return(nil, nil)
	mockNotifRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
	mockAggregator.EXPECT().IncrementRateLimit(ctx, userID, notifType).Return(nil)
	mockTokenRepo.EXPECT().FindByUserID(ctx, userID).Return([]*entity.UserDeviceToken{}, nil)

	dispatcher := service.NewNotificationDispatcher(mockNotifRepo, mockPrefRepo, mockTokenRepo, mockAggregator, mockFirebase)
	assert.NoError(t, dispatcher.Notify(ctx, userID, notifType, data))
}

func TestNotificationDispatcher_Notify_SaveError_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := uuid.New()
	notifType := entity.NotificationTypeNewFollower
	targetID := uuid.New()
	
	data := map[string]interface{}{"target_id": targetID.String()}

	mockNotifRepo := mocks.NewMockNotificationRepository(ctrl)
	mockPrefRepo := mocks.NewMockNotificationPreferenceRepository(ctrl)
	mockTokenRepo := mocks.NewMockDeviceTokenRepository(ctrl)
	mockAggregator := servicemocks.NewMockNotificationAggregator(ctrl)
	mockFirebase := servicemocks.NewMockFirebaseAdapter(ctrl)

	mockPrefRepo.EXPECT().IsEnabled(ctx, userID, notifType, "push").Return(true, nil)
	mockAggregator.EXPECT().CheckRateLimit(ctx, userID, notifType).Return(true, nil)
	mockAggregator.EXPECT().ShouldAggregate(ctx, userID, notifType, targetID).Return(nil, nil)
	mockNotifRepo.EXPECT().Save(ctx, gomock.Any()).Return(errors.New("database error"))

	dispatcher := service.NewNotificationDispatcher(mockNotifRepo, mockPrefRepo, mockTokenRepo, mockAggregator, mockFirebase)
	err := dispatcher.Notify(ctx, userID, notifType, data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}

func TestNotificationDispatcher_Notify_TokenRepoError_LogsAndContinues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := uuid.New()
	notifType := entity.NotificationTypeNewFollower
	targetID := uuid.New()
	actorID := uuid.New()
	
	data := map[string]interface{}{
		"target_id":  targetID.String(),
		"actor_id":   actorID.String(),
		"actor_name": "Test User",
	}

	mockNotifRepo := mocks.NewMockNotificationRepository(ctrl)
	mockPrefRepo := mocks.NewMockNotificationPreferenceRepository(ctrl)
	mockTokenRepo := mocks.NewMockDeviceTokenRepository(ctrl)
	mockAggregator := servicemocks.NewMockNotificationAggregator(ctrl)
	mockFirebase := servicemocks.NewMockFirebaseAdapter(ctrl)

	mockPrefRepo.EXPECT().IsEnabled(ctx, userID, notifType, "push").Return(true, nil)
	mockAggregator.EXPECT().CheckRateLimit(ctx, userID, notifType).Return(true, nil)
	mockAggregator.EXPECT().ShouldAggregate(ctx, userID, notifType, targetID).Return(nil, nil)
	mockNotifRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
	mockAggregator.EXPECT().IncrementRateLimit(ctx, userID, notifType).Return(nil)
	mockTokenRepo.EXPECT().FindByUserID(ctx, userID).Return(nil, errors.New("token repo error"))

	dispatcher := service.NewNotificationDispatcher(mockNotifRepo, mockPrefRepo, mockTokenRepo, mockAggregator, mockFirebase)
	assert.NoError(t, dispatcher.Notify(ctx, userID, notifType, data))
}

func TestNotificationDispatcher_Notify_FirebaseError_LogsAndContinues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := uuid.New()
	notifType := entity.NotificationTypeNewFollower
	targetID := uuid.New()
	actorID := uuid.New()
	
	data := map[string]interface{}{
		"target_id":  targetID.String(),
		"actor_id":   actorID.String(),
		"actor_name": "Test User",
	}

	mockNotifRepo := mocks.NewMockNotificationRepository(ctrl)
	mockPrefRepo := mocks.NewMockNotificationPreferenceRepository(ctrl)
	mockTokenRepo := mocks.NewMockDeviceTokenRepository(ctrl)
	mockAggregator := servicemocks.NewMockNotificationAggregator(ctrl)
	mockFirebase := servicemocks.NewMockFirebaseAdapter(ctrl)

	mockPrefRepo.EXPECT().IsEnabled(ctx, userID, notifType, "push").Return(true, nil)
	mockAggregator.EXPECT().CheckRateLimit(ctx, userID, notifType).Return(true, nil)
	mockAggregator.EXPECT().ShouldAggregate(ctx, userID, notifType, targetID).Return(nil, nil)
	mockNotifRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
	mockAggregator.EXPECT().IncrementRateLimit(ctx, userID, notifType).Return(nil)

	tokens := []*entity.UserDeviceToken{{DeviceToken: "token1", Platform: "ios"}}
	mockTokenRepo.EXPECT().FindByUserID(ctx, userID).Return(tokens, nil)
	mockFirebase.EXPECT().SendPushToUser(ctx, userID, gomock.Any(), gomock.Any(), data).Return(errors.New("firebase error"))

	dispatcher := service.NewNotificationDispatcher(mockNotifRepo, mockPrefRepo, mockTokenRepo, mockAggregator, mockFirebase)
	assert.NoError(t, dispatcher.Notify(ctx, userID, notifType, data))
}

func TestNotificationDispatcher_Notify_AggregatorError_ContinuesAnyway(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := uuid.New()
	notifType := entity.NotificationTypeNewFollower
	targetID := uuid.New()
	actorID := uuid.New()
	
	data := map[string]interface{}{
		"target_id":  targetID.String(),
		"actor_id":   actorID.String(),
		"actor_name": "Test User",
	}

	mockNotifRepo := mocks.NewMockNotificationRepository(ctrl)
	mockPrefRepo := mocks.NewMockNotificationPreferenceRepository(ctrl)
	mockTokenRepo := mocks.NewMockDeviceTokenRepository(ctrl)
	mockAggregator := servicemocks.NewMockNotificationAggregator(ctrl)
	mockFirebase := servicemocks.NewMockFirebaseAdapter(ctrl)

	mockPrefRepo.EXPECT().IsEnabled(ctx, userID, notifType, "push").Return(true, nil)
	mockAggregator.EXPECT().CheckRateLimit(ctx, userID, notifType).Return(true, nil)
	mockAggregator.EXPECT().ShouldAggregate(ctx, userID, notifType, targetID).Return(nil, errors.New("aggregator error"))
	mockNotifRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
	mockAggregator.EXPECT().IncrementRateLimit(ctx, userID, notifType).Return(nil)

	tokens := []*entity.UserDeviceToken{{DeviceToken: "token1", Platform: "ios"}}
	mockTokenRepo.EXPECT().FindByUserID(ctx, userID).Return(tokens, nil)
	mockFirebase.EXPECT().SendPushToUser(ctx, userID, gomock.Any(), gomock.Any(), data).Return(nil)

	dispatcher := service.NewNotificationDispatcher(mockNotifRepo, mockPrefRepo, mockTokenRepo, mockAggregator, mockFirebase)
	assert.NoError(t, dispatcher.Notify(ctx, userID, notifType, data))
}

func TestNotificationDispatcher_Notify_RateLimitError_ContinuesAnyway(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := uuid.New()
	notifType := entity.NotificationTypeNewFollower
	targetID := uuid.New()
	actorID := uuid.New()
	
	data := map[string]interface{}{
		"target_id":  targetID.String(),
		"actor_id":   actorID.String(),
		"actor_name": "Test User",
	}

	mockNotifRepo := mocks.NewMockNotificationRepository(ctrl)
	mockPrefRepo := mocks.NewMockNotificationPreferenceRepository(ctrl)
	mockTokenRepo := mocks.NewMockDeviceTokenRepository(ctrl)
	mockAggregator := servicemocks.NewMockNotificationAggregator(ctrl)
	mockFirebase := servicemocks.NewMockFirebaseAdapter(ctrl)

	mockPrefRepo.EXPECT().IsEnabled(ctx, userID, notifType, "push").Return(true, nil)
	mockAggregator.EXPECT().CheckRateLimit(ctx, userID, notifType).Return(false, errors.New("rate limit error"))
	mockAggregator.EXPECT().ShouldAggregate(ctx, userID, notifType, targetID).Return(nil, nil)
	mockNotifRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil)
	mockAggregator.EXPECT().IncrementRateLimit(ctx, userID, notifType).Return(nil)

	tokens := []*entity.UserDeviceToken{{DeviceToken: "token1", Platform: "ios"}}
	mockTokenRepo.EXPECT().FindByUserID(ctx, userID).Return(tokens, nil)
	mockFirebase.EXPECT().SendPushToUser(ctx, userID, gomock.Any(), gomock.Any(), data).Return(nil)

	dispatcher := service.NewNotificationDispatcher(mockNotifRepo, mockPrefRepo, mockTokenRepo, mockAggregator, mockFirebase)
	assert.NoError(t, dispatcher.Notify(ctx, userID, notifType, data))
}

func TestNotificationDispatcher_Notify_GeneratesProperNotification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	userID := uuid.New()
	notifType := entity.NotificationTypeBlogLike
	actorID := uuid.New()
	targetID := uuid.New()
	actorName := "Test User"
	
	data := map[string]interface{}{
		"target_id":  targetID.String(),
		"actor_id":   actorID.String(),
		"actor_name": actorName,
		"blog_title": "Test Blog",
	}

	mockNotifRepo := mocks.NewMockNotificationRepository(ctrl)
	mockPrefRepo := mocks.NewMockNotificationPreferenceRepository(ctrl)
	mockTokenRepo := mocks.NewMockDeviceTokenRepository(ctrl)
	mockAggregator := servicemocks.NewMockNotificationAggregator(ctrl)
	mockFirebase := servicemocks.NewMockFirebaseAdapter(ctrl)

	var savedNotif *entity.Notification

	mockPrefRepo.EXPECT().IsEnabled(ctx, userID, notifType, "push").Return(true, nil)
	mockAggregator.EXPECT().CheckRateLimit(ctx, userID, notifType).Return(true, nil)
	mockAggregator.EXPECT().ShouldAggregate(ctx, userID, notifType, targetID).Return(nil, nil)
	
	mockNotifRepo.EXPECT().Save(ctx, gomock.Any()).Do(func(_ context.Context, notif *entity.Notification) { savedNotif = notif }).Return(nil)
	mockAggregator.EXPECT().IncrementRateLimit(ctx, userID, notifType).Return(nil)

	tokens := []*entity.UserDeviceToken{{DeviceToken: "token1", Platform: "ios"}}
	mockTokenRepo.EXPECT().FindByUserID(ctx, userID).Return(tokens, nil)
	mockFirebase.EXPECT().SendPushToUser(ctx, userID, gomock.Any(), gomock.Any(), data).Return(nil)

	dispatcher := service.NewNotificationDispatcher(mockNotifRepo, mockPrefRepo, mockTokenRepo, mockAggregator, mockFirebase)
	assert.NoError(t, dispatcher.Notify(ctx, userID, notifType, data))
	
	assert.NotNil(t, savedNotif)
	assert.Equal(t, userID, savedNotif.UserID)
	assert.Equal(t, notifType, savedNotif.Type)
	assert.Equal(t, actorName, savedNotif.Data["actor_name"])
	assert.Equal(t, actorID.String(), savedNotif.Data["actor_id"])
	assert.Equal(t, targetID.String(), savedNotif.Data["target_id"])
	assert.Equal(t, 1, savedNotif.GroupedCount)
	assert.False(t, savedNotif.IsRead)
	assert.True(t, savedNotif.ExpiresAt.After(time.Now()))
}
