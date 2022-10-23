package smurlrepo

import (
	"context"
	"fmt"
	"log"

	"github.com/sanyarise/smurl/internal/entities/smurlentity"
	"github.com/sanyarise/smurl/internal/infrastructure/api/helpers"

	"go.uber.org/zap"
)

// Interface for communication with the database
type SmurlStore interface {
	CreateURL(ctx context.Context, ses smurlentity.Smurl) (*smurlentity.Smurl, error)
	FindURL(ctx context.Context, url smurlentity.Smurl) (*smurlentity.Smurl, error)
	UpdateStat(ctx context.Context, url smurlentity.Smurl) (*smurlentity.Smurl, error)
	ReadStat(ctx context.Context, url smurlentity.Smurl) (*smurlentity.Smurl, error)
}

type SmurlStorage struct {
	smurlStore SmurlStore
	logger     *zap.Logger
}

func NewSmurlStorage(smurlStore SmurlStore, l *zap.Logger) *SmurlStorage {
	l.Debug("Enter in func smurlrepo NewSmurlStorage()")
	return &SmurlStorage{
		smurlStore: smurlStore,
		logger:     l,
	}
}

func (ss *SmurlStorage) CreateURL(ctx context.Context, ses smurlentity.Smurl) (*smurlentity.Smurl, error) {
	l := ss.logger
	l.Debug("Enter in smurlrepo func CreateURL()")

	newSmurl, err := ss.smurlStore.CreateURL(ctx, ses)
	if err != nil {
		l.Error("",
			zap.Error(err))
		return nil, fmt.Errorf("create url error: %w", err)
	}
	return newSmurl, nil
}

func (ss *SmurlStorage) CreateStat(ctx context.Context, ses smurlentity.Smurl) (*smurlentity.Smurl, error) {
	l := ss.logger
	l.Debug("Enter in smurlrepo func CreateStat()")
	// Search for a small url in the database
	ru, err := ss.smurlStore.FindURL(ctx, ses)
	if err != nil {
		l.Error("",
			zap.Error(err))
		return nil, fmt.Errorf("error find small url: %s", err)
	}
	// Update the hit counter field
	ru.Count, err = helpers.CountUses(ru.Count, l)
	if err != nil {
		l.Error("",
			zap.Error(err))
	}
	l.Debug("Smurlrepo after find and count ++",
		zap.Any(":", ru))
	// Update the field with IP information
	ru.IPInfo = ru.IPInfo + ses.IPInfo
	l.Debug("IP info",
		zap.String(":", ru.IPInfo))
	// Call the database method to update statistics
	sus, err := ss.smurlStore.UpdateStat(ctx, *ru)
	if err != nil {
		l.Error("",
			zap.Error(err))
		return nil, fmt.Errorf("update stat error: %w", err)
	}
	return sus, nil
}

func (ss *SmurlStorage) FindURL(ctx context.Context, url smurlentity.Smurl) (*smurlentity.Smurl, error) {
	l := ss.logger
	l.Debug("Enter in smurlrepo func FindURL()")
	ru, err := ss.smurlStore.FindURL(ctx, url)
	if err != nil {
		l.Error("",
			zap.Error(err))
		return nil, fmt.Errorf("find url error: %w", err)
	}
	return ru, nil
}

func (ss *SmurlStorage) ReadStat(ctx context.Context, ses smurlentity.Smurl) (*smurlentity.Smurl, error) {
	l := ss.logger
	l.Debug("Enter in smurlrepo func ReadStat()")
	rss, err := ss.smurlStore.ReadStat(ctx, ses)
	if err != nil {
		l.Error("",
			zap.Error(err))
		return nil, fmt.Errorf("error on read statistics:%w", err)
	}
	log.Println(rss)
	return rss, nil
}
