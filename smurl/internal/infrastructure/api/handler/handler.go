package handler

import (
	"context"
	"fmt"

	"github.com/sanyarise/smurl/internal/entities/smurlentity"
	"github.com/sanyarise/smurl/internal/usecases/repos/smurlrepo"

	"go.uber.org/zap"
)

type Handlers struct {
	sm     *smurlrepo.SmurlStorage
	logger *zap.Logger
}

func NewHandlers(sm *smurlrepo.SmurlStorage, l *zap.Logger) *Handlers {
	h := &Handlers{
		sm:     sm,
		logger: l,
	}
	return h
}

type Smurl struct {
	SmallURL string `json:"small_url,omitempty"`
	LongURL  string `json:"long_url,omitempty"`
	AdminURL string `json:"admin_url,omitempty"`
	IPInfo   string `json:"ip_info,omitempty"`
	Count    string `json:"count,omitempty"`
}

//Обработчик эндпоинта для создания уменьшенного url
func (h *Handlers) CreateSmurlHandle(ctx context.Context, hss Smurl) (Smurl, error) {
	l := h.logger
	l.Debug("handlers create smurl handle")
	ses := smurlentity.Smurl{LongURL: hss.LongURL}
	//Вызов метода из слоя с интерфейсами
	nses, err := h.sm.CreateURL(ctx, ses)
	if err != nil {
		l.Error("",
			zap.Error(err))
		return Smurl{}, fmt.Errorf("error when creating: %w", err)
	}
	return Smurl{
		LongURL:  nses.LongURL,
		SmallURL: nses.SmallURL,
		AdminURL: nses.AdminURL,
	}, nil
}

//Обработчик эндпоинта для поиска в бд уменьшенного url, обновления
//статистики и переадресации по найденному длинному адресу
func (h *Handlers) RedirectHandle(ctx context.Context, smallURL string, ip string) (Smurl, error) {
	l := h.logger
	l.Debug("handlers redirect handle")

	es := smurlentity.Smurl{
		SmallURL: smallURL,
		IPInfo:   ip,
	}
	//Вызов метода из слоя интерфейсов для создания статистики перед переходом по длиной ссылке
	hs, err := h.sm.CreateStat(ctx, es)
	if err != nil {
		l.Error("",
			zap.Error(err))
		return Smurl{}, fmt.Errorf("redirect error: %w", err)
	}

	return Smurl{
		LongURL: hs.LongURL,
	}, nil
}

//Обработчик эндпоинта для поиска админского урл в бд и получения статистики
//переходов по уменьшенному урл
func (h *Handlers) GetStatHandle(ctx context.Context, sm Smurl) (Smurl, error) {
	l := h.logger
	l.Debug("handlers get stat handle")
	es := smurlentity.Smurl{
		AdminURL: sm.AdminURL,
	}
	//Вызов метода для чтения статистики из слоя интерфейсов
	cu, err := h.sm.ReadStat(ctx, es)
	if err != nil {
		l.Error("",
			zap.Error(err))
		return Smurl{}, fmt.Errorf("creating stat error:%w", err)
	}
	scu := Smurl(*cu)
	return scu, nil
}
