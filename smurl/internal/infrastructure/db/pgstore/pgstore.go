package pgstore

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sanyarise/smurl/internal/entities/smurlentity"
	"github.com/sanyarise/smurl/internal/infrastructure/api/helpers"
	"github.com/sanyarise/smurl/internal/usecases/repos/smurlrepo"

	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
)

//Слой для работы с базой данных Postgresql
var _ smurlrepo.SmurlStore = &SmurlStore{}

type PgSmurl struct {
	SmallURL  string
	CreatedAt time.Time
	LongURL   string
	AdminURL  string
	IPInfo    string
	Count     string
}

type SmurlStore struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewSmurlStore(dsn string, l *zap.Logger) (*SmurlStore, error) {
	l.Debug("pgstore new smurl store")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		l.Error("error on sql open",
			zap.Error(err))
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		l.Error("error on db ping",
			zap.Error(err))
		db.Close()
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS smurls (
		small_url varchar NOT NULL,
		created_at timestamptz NOT NULL,
		long_url varchar NOT NULL,
		admin_url varchar NOT NULL,
		count varchar,
		ip_info varchar,

		CONSTRAINT smurls_pk PRIMARY KEY (small_url)
		)`)
	if err != nil {
		l.Error("error on create table",
			zap.Error(err))
		db.Close()
		return nil, err
	}
	su := &SmurlStore{
		db:     db,
		logger: l,
	}
	return su, nil
}
func (su *SmurlStore) Close() {
	su.logger.Debug("db close")
	su.db.Close()
}

//Метод для сохранения длинного урл, короткого урл и админского урл в базу данных
func (ss *SmurlStore) CreateURL(ctx context.Context, ses smurlentity.Smurl) (*smurlentity.Smurl, error) {
	l := ss.logger
	l.Debug("pgstore create url")
	//Создание короткого и админского урл
	smallURL := helpers.RandString(l)
	adminURL := helpers.RandString(l)

	pgs := &PgSmurl{
		LongURL:   ses.LongURL,
		CreatedAt: time.Now(),
		SmallURL:  smallURL,
		AdminURL:  adminURL,
		Count:     "0",
		IPInfo:    "",
	}
	//Начало транзакции для записи данных в бд
	tx, err := ss.db.Begin()
	if err != nil {
		l.Error("error on begin transaction",
			zap.Error(err))
	}
	//Запись в бд
	_, err = tx.ExecContext(ctx, `INSERT INTO smurls
	(small_url, created_at, long_url, admin_url, count, ip_info)
	values ($1, $2, $3, $4, $5, $6)`,
		pgs.SmallURL,
		pgs.CreatedAt,
		pgs.LongURL,
		pgs.AdminURL,
		pgs.Count,
		pgs.IPInfo,
	)
	if err != nil {
		//Возврат к исходному значению
		//в случае неудачной записи
		tx.Rollback()
		l.Error("error on insert values into table",
			zap.Error(err))
		return nil, err
	}
	//Завершение транзакции
	tx.Commit()
	l.Debug("pg store create smurl successfull")
	//Возвращаем объект с коротким и админским урл
	return &smurlentity.Smurl{
		SmallURL: pgs.SmallURL,
		AdminURL: pgs.AdminURL,
	}, nil
}

//Метод для обновления данных статистики при переходе по уменьшенному урл
func (ss *SmurlStore) UpdateStat(ctx context.Context, ses smurlentity.Smurl) (*smurlentity.Smurl, error) {
	l := ss.logger
	l.Debug("pgstore update stat")
	pgs := &PgSmurl{
		LongURL:  ses.LongURL,
		Count:    ses.Count,
		IPInfo:   ses.IPInfo,
		SmallURL: ses.SmallURL,
	}
	//Начало транзакции для записи обновленных данных
	tx, err := ss.db.Begin()
	if err != nil {
		l.Error("error on begin transaction",
			zap.Error(err))
	}
	//Запись обновленных данных
	_, err = tx.ExecContext(ctx, `UPDATE smurls SET count = $1, ip_info = $2
	WHERE small_url = $3`, pgs.Count, pgs.IPInfo, pgs.SmallURL)
	if err != nil {
		l.Error("error on update values into table",
			zap.Error(err))

		//Возврат к исходному значению
		//в случае неудачной записи
		tx.Rollback()
		return nil, err
	}
	//Завершение транзакции
	tx.Commit()
	l.Debug("pg store update stat successfull")

	return &smurlentity.Smurl{
		LongURL: pgs.LongURL,
		Count:   pgs.Count,
		IPInfo:  pgs.IPInfo,
	}, nil
}

//Метод для чтения данных статистики
func (ss *SmurlStore) ReadStat(ctx context.Context, ses smurlentity.Smurl) (*smurlentity.Smurl, error) {
	l := ss.logger
	l.Debug("pgstore read stat")
	pgs := &PgSmurl{}
	//Выполняем поиск по базе данных
	rows, err := ss.db.QueryContext(ctx,
		`SELECT small_url, long_url, admin_url, count, ip_info FROM smurls
	 WHERE admin_url = $1`, ses.AdminURL)
	if err != nil {
		l.Error("error on query in table",
			zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&pgs.SmallURL,
			&pgs.LongURL,
			&pgs.AdminURL,
			&pgs.Count,
			&pgs.IPInfo,
		); err != nil {
			l.Error("error on rows scan",
				zap.Error(err))
			return nil, err
		}
	}
	err = fmt.Errorf("admin_url: %s not found", ses.AdminURL)
	if pgs.AdminURL == "" {
		l.Error("",
			zap.Error(err))
		return nil, err
	}
	sms := &smurlentity.Smurl{
		SmallURL: pgs.SmallURL,
		LongURL:  pgs.LongURL,
		AdminURL: pgs.AdminURL,
		Count:    pgs.Count,
		IPInfo:   pgs.IPInfo,
	}
	l.Debug("pg store read stat successfull")

	return sms, nil
}

//Метод для поиска малого урл в базе данных
func (ss *SmurlStore) FindURL(ctx context.Context, ses smurlentity.Smurl) (*smurlentity.Smurl, error) {
	l := ss.logger
	l.Debug("pg store find url")

	pgs := &PgSmurl{}
	rows, err := ss.db.QueryContext(ctx,
		`SELECT small_url, long_url, count, ip_info FROM smurls WHERE small_url = $1`, ses.SmallURL)
	if err != nil {
		l.Error("error on search small url into table",
			zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&pgs.SmallURL,
			&pgs.LongURL,
			&pgs.Count,
			&pgs.IPInfo,
		); err != nil {
			l.Error("error find small url",
				zap.Error(err))
			return nil, err
		}
	}
	if pgs.SmallURL == "" {
		err = fmt.Errorf("small_url not found")
		l.Error("",
			zap.Error(err))
		return nil, err
	}
	l.Debug("pg store find url successfull")
	return &smurlentity.Smurl{
		SmallURL: ses.SmallURL,
		LongURL:  pgs.LongURL,
		Count:    pgs.Count,
		IPInfo:   pgs.IPInfo,
	}, nil
}
