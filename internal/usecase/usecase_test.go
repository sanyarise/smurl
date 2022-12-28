package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	helpers "github.com/sanyarise/smurl/internal/helpers/mocks"
	"github.com/sanyarise/smurl/internal/models"
	"github.com/sanyarise/smurl/internal/repository/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type TestStatement struct {
	store   *mocks.MockSmurlStore
	helpers *helpers.MockHelper
	logger  *zap.Logger
	usecase *SmurlUsecase
}

func NewTestStatement(ctrl *gomock.Controller) *TestStatement {
	store := mocks.NewMockSmurlStore(ctrl)
	helpers := helpers.NewMockHelper(ctrl)
	logger := zap.L()
	usecase := NewSmurlUsecase(store, helpers, logger)
	return &TestStatement{
		store:   store,
		helpers: helpers,
		logger:  logger,
		usecase: usecase,
	}
}

var (
	ctx             = context.Background()
	testCreateSmurl = models.Smurl{
		LongURL:  "test",
		SmallURL: "test",
		AdminURL: "test",
	}
	testUpdateSmurl = models.Smurl{
		SmallURL: "test",
		AdminURL: "test",
		IPInfo:   []string{"test"},
	}
	testUpdatedSmurl = models.Smurl{
		SmallURL: "test",
		AdminURL: "test",
		IPInfo:   []string{"test"},
		Count:    1,
	}
	err = errors.New("test error")
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := NewTestStatement(ctrl)

	s.helpers.EXPECT().RandString().Return("test")
	s.helpers.EXPECT().RandString().Return("test")
	s.store.EXPECT().Create(ctx, testCreateSmurl).Return(nil, err)
	res, err := s.usecase.Create(ctx, "test")
	require.Error(t, err)
	require.Nil(t, res)

	s.helpers.EXPECT().RandString().Return("test")
	s.helpers.EXPECT().RandString().Return("test")
	s.store.EXPECT().Create(ctx, testCreateSmurl).Return(&testCreateSmurl, nil)
	res, err = s.usecase.Create(ctx, "test")
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, &testCreateSmurl)
}

func TestUpdateStat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := NewTestStatement(ctrl)

	s.store.EXPECT().UpdateStat(ctx, testUpdatedSmurl).Return(err)
	err := s.usecase.UpdateStat(ctx, testUpdateSmurl)
	require.Error(t, err)

	s.store.EXPECT().UpdateStat(ctx, testUpdatedSmurl).Return(nil)
	err = s.usecase.UpdateStat(ctx, testUpdateSmurl)
	require.NoError(t, err)
}

func TestFindURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := NewTestStatement(ctrl)

	s.store.EXPECT().FindURL(ctx, "test").Return(nil, err)
	res, err := s.usecase.FindURL(ctx, "test")
	require.Error(t, err)
	require.Nil(t, res)

	s.store.EXPECT().FindURL(ctx, "test").Return(&testCreateSmurl, nil)
	res, err = s.usecase.FindURL(ctx, "test")
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, &testCreateSmurl)
}

func TestReadStat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := NewTestStatement(ctrl)

	s.store.EXPECT().ReadStat(ctx, "test").Return(nil, err)
	res, err := s.usecase.ReadStat(ctx, "test")
	require.Error(t, err)
	require.Nil(t, res)

	s.store.EXPECT().ReadStat(ctx, "test").Return(&testCreateSmurl, nil)
	res, err = s.usecase.ReadStat(ctx, "test")
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, &testCreateSmurl)
}
