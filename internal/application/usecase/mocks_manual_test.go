package usecase_test

import (
	"context"
	"reflect"

	"github.com/aiagent/internal/domain/entity"
	"github.com/aiagent/internal/domain/service"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

// MockSocialAccountRepository is a mock of SocialAccountRepository interface
type MockSocialAccountRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSocialAccountRepositoryMockRecorder
}

type MockSocialAccountRepositoryMockRecorder struct {
	mock *MockSocialAccountRepository
}

func NewMockSocialAccountRepository(ctrl *gomock.Controller) *MockSocialAccountRepository {
	mock := &MockSocialAccountRepository{ctrl: ctrl}
	mock.recorder = &MockSocialAccountRepositoryMockRecorder{mock}
	return mock
}

func (m *MockSocialAccountRepository) EXPECT() *MockSocialAccountRepositoryMockRecorder {
	return m.recorder
}

func (m *MockSocialAccountRepository) Create(ctx context.Context, socialAccount *entity.SocialAccount) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, socialAccount)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockSocialAccountRepositoryMockRecorder) Create(ctx, socialAccount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSocialAccountRepository)(nil).Create), ctx, socialAccount)
}

func (m *MockSocialAccountRepository) FindByProvider(ctx context.Context, provider, providerID string) (*entity.SocialAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByProvider", ctx, provider, providerID)
	ret0, _ := ret[0].(*entity.SocialAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockSocialAccountRepositoryMockRecorder) FindByProvider(ctx, provider, providerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByProvider", reflect.TypeOf((*MockSocialAccountRepository)(nil).FindByProvider), ctx, provider, providerID)
}

func (m *MockSocialAccountRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.SocialAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUserID", ctx, userID)
	ret0, _ := ret[0].([]*entity.SocialAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockSocialAccountRepositoryMockRecorder) GetByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUserID", reflect.TypeOf((*MockSocialAccountRepository)(nil).GetByUserID), ctx, userID)
}

// MockSocialAuthService is a mock of SocialAuthService interface
type MockSocialAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockSocialAuthServiceMockRecorder
}

type MockSocialAuthServiceMockRecorder struct {
	mock *MockSocialAuthService
}

func NewMockSocialAuthService(ctrl *gomock.Controller) *MockSocialAuthService {
	mock := &MockSocialAuthService{ctrl: ctrl}
	mock.recorder = &MockSocialAuthServiceMockRecorder{mock}
	return mock
}

func (m *MockSocialAuthService) EXPECT() *MockSocialAuthServiceMockRecorder {
	return m.recorder
}

func (m *MockSocialAuthService) GetUserInfo(ctx context.Context, provider, code string) (*service.SocialUserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserInfo", ctx, provider, code)
	ret0, _ := ret[0].(*service.SocialUserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockSocialAuthServiceMockRecorder) GetUserInfo(ctx, provider, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserInfo", reflect.TypeOf((*MockSocialAuthService)(nil).GetUserInfo), ctx, provider, code)
}

func (m *MockSocialAuthService) GetAuthURL(ctx context.Context, provider string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAuthURL", ctx, provider)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockSocialAuthServiceMockRecorder) GetAuthURL(ctx, provider interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAuthURL", reflect.TypeOf((*MockSocialAuthService)(nil).GetAuthURL), ctx, provider)
}
