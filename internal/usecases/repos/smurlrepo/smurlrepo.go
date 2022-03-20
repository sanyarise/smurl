package smurlrepo

import (
	"context"
	"fmt"
	"log"

	"github.com/sanyarise/smurl/internal/entities/smurlentity"
	"github.com/sanyarise/smurl/internal/infrastructure/api/helpers"

	"go.uber.org/zap"
)

//Интерфейс-порт для связи с базой данных
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
	return &SmurlStorage{
		smurlStore: smurlStore,
		logger:     l,
	}
}

//Промежуточный метод для перенаправления запроса на слой с базой данных
func (ss *SmurlStorage) CreateURL(ctx context.Context, ses smurlentity.Smurl) (*smurlentity.Smurl, error) {
	l := ss.logger
	l.Debug("smurlrepo create url")

	newSmurl, err := ss.smurlStore.CreateURL(ctx, ses)
	if err != nil {
		l.Error("",
			zap.Error(err))
		return nil, fmt.Errorf("create url error: %w", err)
	}
	return newSmurl, nil
}

//Промежуточный метод для перенаправления запроса на слой с базой данных
func (ss *SmurlStorage) CreateStat(ctx context.Context, ses smurlentity.Smurl) (*smurlentity.Smurl, error) {
	l := ss.logger
	l.Debug("smurlrepo create stat")
	//Сначала выполняем поиск малого урл в бд
	ru, err := ss.smurlStore.FindURL(ctx, ses)
	if err != nil {
		l.Error("",
			zap.Error(err))
		return nil, fmt.Errorf("error find small url: %s", err)
	}
	//Обновляем поле с счечиком посещений
	ru.Count, err = helpers.CountUses(ru.Count, l)
	if err != nil {
		l.Error("",
			zap.Error(err))
	}
	l.Debug("smurl after find and count ++",
		zap.Any(":", ru))
	//Обновляем поле с информацией об IP
	ru.IPInfo = ru.IPInfo + ses.IPInfo
	l.Debug("ip info",
		zap.String(":", ru.IPInfo))
	//Вызываем метод базы данных для обновления
	//статистики
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
	l.Debug("smurlrepo find url")
	ru, err := ss.smurlStore.FindURL(ctx, url)
	if err != nil {
		l.Error("",
			zap.Error(err))
		return nil, fmt.Errorf("find url error: %w", err)
	}
	return ru, nil
}

//Промежуточный метод для перенаправления на слой с базой данных для получения статистики
func (ss *SmurlStorage) ReadStat(ctx context.Context, ses smurlentity.Smurl) (*smurlentity.Smurl, error) {
	l := ss.logger
	l.Debug("smurlstorage Read stat")
	rss, err := ss.smurlStore.ReadStat(ctx, ses)
	if err != nil {
		l.Error("",
			zap.Error(err))
		return nil, fmt.Errorf("error on read statistics:%w", err)
	}
	log.Println(rss)
	return rss, nil
}
