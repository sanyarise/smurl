package usecase

import (
	"context"
	"fmt"

	"github.com/sanyarise/smurl/internal/helpers"
	"github.com/sanyarise/smurl/internal/models"
	"go.uber.org/zap"
)

// Interface for communication with the database
type SmurlStore interface {
	Create(ctx context.Context, smurl models.Smurl) (*models.Smurl, error)
	UpdateStat(ctx context.Context, smurl models.Smurl) error
	ReadStat(ctx context.Context, adminUrl string) (*models.Smurl, error)
	FindURL(ctx context.Context, smallUrl string) (*models.Smurl, error)
}

var _ Usecase = SmurlUsecase{}

type SmurlUsecase struct {
	repository SmurlStore
	helpers    helpers.Helper
	logger     *zap.Logger
}

func NewSmurlUsecase(smurlStore SmurlStore, helpers helpers.Helper, logger *zap.Logger) *SmurlUsecase {
	logger.Debug("Enter in usecase NewSmurlUsecase()")
	return &SmurlUsecase{
		repository: smurlStore,
		helpers:    helpers,
		logger:     logger,
	}
}

func (usecase SmurlUsecase) Create(ctx context.Context, longUrl string) (*models.Smurl, error) {
	usecase.logger.Debug("Enter in usecase Create()")
	createdSmurl := models.Smurl{
		LongURL: longUrl,
	}
	createdSmurl.SmallURL = usecase.helpers.RandString()
	createdSmurl.AdminURL = usecase.helpers.RandString()

	smurl, err := usecase.repository.Create(ctx, createdSmurl)
	if err != nil {
		usecase.logger.Error("",
			zap.Error(err))
		return nil, fmt.Errorf("create url error: %w", err)
	}
	return smurl, nil
}

func (usecase SmurlUsecase) UpdateStat(ctx context.Context, updatedSmurl models.Smurl) error {
	usecase.logger.Debug("Enter in usecase UpdateStat()")
	// Update the hit counter field
	updatedSmurl.Count++
	// Call the database method to update statistics
	err := usecase.repository.UpdateStat(ctx, updatedSmurl)
	if err != nil {
		usecase.logger.Error("",
			zap.Error(err))
		return fmt.Errorf("update stat error: %w", err)
	}
	return nil
}

func (usecase SmurlUsecase) FindURL(ctx context.Context, smallUrl string) (*models.Smurl, error) {
	usecase.logger.Debug("Enter in usecase FindURL()")
	smurl, err := usecase.repository.FindURL(ctx, smallUrl)
	if err != nil {
		return nil, err
	}
	return smurl, nil
}

func (usecase SmurlUsecase) ReadStat(ctx context.Context, adminUrl string) (*models.Smurl, error) {
	usecase.logger.Debug("Enter in usecase ReadStat()")
	smurl, err := usecase.repository.ReadStat(ctx, adminUrl)
	if err != nil {
		return nil, err
	}
	return smurl, nil
}
