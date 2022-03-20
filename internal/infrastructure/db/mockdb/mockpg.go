package mockdb

import (
	"context"
	"fmt"

	smurlentity "github.com/sanyarise/smurl/internal/entities/smurlentity"
	"github.com/sanyarise/smurl/internal/usecases/repos/smurlrepo"

	"go.uber.org/zap"
)

//Имитация ответов базы данных для тестов
var _ smurlrepo.SmurlStore = &Mockdb{}

type Mockdb struct {
	logger *zap.Logger
}

func NewMockdb(a string, l *zap.Logger) *Mockdb {
	return &Mockdb{
		logger: l,
	}
}

func (mc *Mockdb) CreateURL(ctx context.Context, ses smurlentity.Smurl) (*smurlentity.Smurl, error) {
	return &smurlentity.Smurl{
		SmallURL: "abcdefgh",
		AdminURL: "zxcvbnma",
	}, nil
}

func (mc *Mockdb) FindURL(ctx context.Context, ses smurlentity.Smurl) (*smurlentity.Smurl, error) {
	if ses.SmallURL != "abcdefgh" {
		return nil, fmt.Errorf("error on find url")
	}
	return &smurlentity.Smurl{
		SmallURL: ses.SmallURL,
		LongURL:  "http://vk.com",
		Count:    "1",
		IPInfo:   "0.0.0.0",
	}, nil
}

func (mc *Mockdb) UpdateStat(ctx context.Context, ses smurlentity.Smurl) (*smurlentity.Smurl, error) {
	return &smurlentity.Smurl{
		LongURL: ses.LongURL,
		Count:   "1",
		IPInfo:  "0.0.0.0",
	}, nil
}

func (mc *Mockdb) ReadStat(ctx context.Context, ses smurlentity.Smurl) (*smurlentity.Smurl, error) {
	if ses.AdminURL != "zxcvbnma" {
		return nil, fmt.Errorf("error on read stat")
	}
	return &smurlentity.Smurl{
		SmallURL: "abcdefgh",
		LongURL:  "http://vk.com",
		AdminURL: "zxcvbnma",
		Count:    "1",
		IPInfo:   "0.0.0.0",
	}, nil
}
