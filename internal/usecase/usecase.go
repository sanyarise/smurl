package usecase

import (
	"context"
	"fmt"
	"log"

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
	logger     *zap.Logger
}

func NewSmurlUsecase(smurlStore SmurlStore, logger *zap.Logger) *SmurlUsecase {
	logger.Debug("Enter in usecase NewSmurlUsecase()")
	return &SmurlUsecase{
		repository: smurlStore,
		logger:     logger,
	}
}

func (usecase *SmurlUsecase) Create(ctx context.Context, longUrl string) (*models.Smurl, error) {
	usecase.logger.Debug("Enter in smurlrepo func CreateURL()")
	createdSmurl.SmallURL = helpers.RandString(usecase.logger)
	createdSmurl.AdminURL = helpers.RandString(usecase.logger)

	smurl, err := usecase.repository.Create(ctx, createdSmurl)
	if err != nil {
		usecase.logger.Error("",
			zap.Error(err))
		return nil, fmt.Errorf("create url error: %w", err)
	}
	return smurl, nil
}

func (usecase *SmurlUsecase) UpdateStat(ctx context.Context, updatedSmurl models.Smurl) error {
	usecase.logger.Debug("Enter in smurlrepo func CreateStat()")
	// Search for a small url in the database
	smurl, err := usecase.repository.FindURL(ctx, updatedSmurl.smallURL)
	if err != nil {
		usecase.logger.Error("",
			zap.Error(err))
		return fmt.Errorf("error find small url: %s", err)
	}
	// Update the hit counter field
	smurl.Count, err = helpers.CountUses(smurl.Count, usecase.logger)
	if err != nil {
		usecase.logger.Error("",
			zap.Error(err))
	}
	usecase.logger.Debug("Smurlrepo after find and count ++",
		zap.Any(":", smurl))
	// Update the field with IP information
	smurl.IPInfo = append(smurl.IPInfo, updatedSmurl.IPInfo)
	usecase.logger.Debug("IP info",
		zap.Any(":", smurl.IPInfo))
	// Call the database method to update statistics
	err = usecase.repository.UpdateStat(ctx, smurl)
	if err != nil {
		usecase.logger.Error("",
			zap.Error(err))
		return fmt.Errorf("update stat error: %w", err)
	}
	return nil
}

func (usecase *SmurlUsecase) FindURL(ctx context.Context, smallUrl string) (*models.Smurl, error) {
	usecase.logger.Debug("Enter in smurlrepo func FindURL()")
	smurl, err := usecase.repository.FindURL(ctx, smallUrl)
	if err != nil {
		usecase.logger.Error("",
			zap.Error(err))
		return nil, fmt.Errorf("find url error: %w", err)
	}
	return smurl, nil
}

func (usecase *SmurlUsecase) ReadStat(ctx context.Context, adminUrl string) (*models.Smurl, error) {
	usecase.logger.Debug("Enter in smurlrepo func ReadStat()")
	smurl, err := usecase.repository.ReadStat(ctx, adminUrl)
	if err != nil {
		usecase.logger.Error("",
			zap.Error(err))
		return nil, fmt.Errorf("error on read statistics:%w", err)
	}
	return smurl, nil
}
