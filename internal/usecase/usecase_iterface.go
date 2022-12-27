package usecase

import (
	"context"

	"github.com/sanyarise/smurl/internal/models"
)

type Usecase interface {
	Create(ctx context.Context, longUrl string) (*models.Smurl, error)
	UpdateStat(ctx context.Context, updatedSmurl models.Smurl) error
	FindURL(ctx context.Context, smallUrl string) (*models.Smurl, error)
	ReadStat(ctx context.Context, adminUrl string) (*models.Smurl, error)
}
