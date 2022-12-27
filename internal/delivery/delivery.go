package delivery

import (
	"context"
	"fmt"

	"github.com/sanyarise/smurl/internal/models"
	"github.com/sanyarise/smurl/internal/usecase"

	"go.uber.org/zap"
)

type Delivery struct {
	usecase usecase.Usecase
	logger  *zap.Logger
}

type Smurl struct {
	SmallURL string `json:"small_url,omitempty"`
	LongURL  string `json:"long_url,omitempty"`
	AdminURL string `json:"admin_url,omitempty"`
	IPInfo   string `json:"ip_info,omitempty"`
	Count    string `json:"count,omitempty"`
}

func NewDelivery(usecase usecase.Usecase, logger *zap.Logger) *Delivery{
	logger.Debug("Enter in delivery NewDelivery()")
	handlers := &Handlers{
		repo:   sm,
		logger: l,
	}
	return handlers
}

// Endpoint handler for creating a minified url
func (delivery *Delivery) Create(ctx context.Context, createdSmurl Smurl) (Smurl, error) {
	delivery.logger.Debug("Enter in delivery Create()")
	ses := smurlentity.Smurl{LongURL: hss.LongURL}

	// Calling a method from a layer with interfaces
	newSmurl, err := h.repo.CreateURL(ctx, ses)
	if err != nil {
		l.Error("",
			zap.Error(err))
		return Smurl{}, fmt.Errorf("error when creating: %w", err)
	}
	return Smurl{
		LongURL:  newSmurl.LongURL,
		SmallURL: newSmurl.SmallURL,
		AdminURL: newSmurl.AdminURL,
	}, nil
}

// Endpoint handler for searching the reduced url in the database, updating
// statistics and redirects to the found long address
func (delivery *Delivery) Redirect(ctx context.Context, smallURL string, ip string) (Smurl, error) {
	l := h.logger
	l.Debug("Enter in handlers func RedirectHandle()")

	es := smurlentity.Smurl{
		SmallURL: smallURL,
		IPInfo:   ip,
	}
	// Calling a method from the interface layer to generate statistics before following a long link
	hs, err := h.repo.CreateStat(ctx, es)
	if err != nil {
		l.Error("",
			zap.Error(err))
		return Smurl{}, fmt.Errorf("redirect error: %w", err)
	}

	return Smurl{
		LongURL: hs.LongURL,
	}, nil
}

// Endpoint handler for searching the admin url in the database and getting statistics
// transitions on a reduced url
func (delivery *Delivery) GetStat(ctx context.Context, sm Smurl) (Smurl, error) {
	l := h.logger
	l.Debug("Enter in handlers func GetStatHandle()")
	es := smurlentity.Smurl{
		AdminURL: sm.AdminURL,
	}
	// Calling the method for reading statistics from the interface layer
	cu, err := h.repo.ReadStat(ctx, es)
	if err != nil {
		l.Error("",
			zap.Error(err))
		return Smurl{}, fmt.Errorf("creating stat error:%w", err)
	}
	scu := Smurl(*cu)
	return scu, nil
}
