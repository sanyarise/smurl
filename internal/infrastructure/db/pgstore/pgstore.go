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

// Postgresql database layer
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
	l.Debug("Pgstore NewSmurlStore")
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
	su.logger.Debug("DB close")
	su.db.Close()
}

// CreateURL saving long url, short url and admin url to database
func (ss *SmurlStore) CreateURL(ctx context.Context, ses smurlentity.Smurl) (*smurlentity.Smurl, error) {
	l := ss.logger
	l.Debug("Pgstore CreateURL")
	// Create a short and admin url
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
	// Starting a transaction to write data to the database
	tx, err := ss.db.Begin()
	if err != nil {
		l.Error("error on begin transaction",
			zap.Error(err))
	}
	// Write to database
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
		//Return to original value in case of unsuccessful write
		tx.Rollback()
		l.Error("error on insert values into table",
			zap.Error(err))
		return nil, err
	}
	// End of transaction
	tx.Commit()
	l.Debug("Pgstore create smurl successfull")
	// Return object with short and admin url
	return &smurlentity.Smurl{
		SmallURL: pgs.SmallURL,
		AdminURL: pgs.AdminURL,
	}, nil
}

// UpdateStat updating statistics data when clicking on a reduced url
func (ss *SmurlStore) UpdateStat(ctx context.Context, ses smurlentity.Smurl) (*smurlentity.Smurl, error) {
	l := ss.logger
	l.Debug("Pgstore UpdateStat")
	pgs := &PgSmurl{
		LongURL:  ses.LongURL,
		Count:    ses.Count,
		IPInfo:   ses.IPInfo,
		SmallURL: ses.SmallURL,
	}
	// Starting a transaction to write updated data
	tx, err := ss.db.Begin()
	if err != nil {
		l.Error("error on begin transaction",
			zap.Error(err))
	}
	// Write updated data
	_, err = tx.ExecContext(ctx, `UPDATE smurls SET count = $1, ip_info = $2
	WHERE small_url = $3`, pgs.Count, pgs.IPInfo, pgs.SmallURL)
	if err != nil {
		l.Error("error on update values into table",
			zap.Error(err))

		// Return to original value in case of unsuccessful write
		tx.Rollback()
		return nil, err
	}
	// End of transaction
	tx.Commit()
	l.Debug("Pgstore update stat successfull")

	return &smurlentity.Smurl{
		LongURL: pgs.LongURL,
		Count:   pgs.Count,
		IPInfo:  pgs.IPInfo,
	}, nil
}

// ReadStat reads statistics data
func (ss *SmurlStore) ReadStat(ctx context.Context, ses smurlentity.Smurl) (*smurlentity.Smurl, error) {
	l := ss.logger
	l.Debug("Pgstore ReadStat")
	pgs := &PgSmurl{}
	// Performing a database search
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
	l.Debug("Pgstore read stat successfull")

	return sms, nil
}

// FindURL search small url in database
func (ss *SmurlStore) FindURL(ctx context.Context, ses smurlentity.Smurl) (*smurlentity.Smurl, error) {
	l := ss.logger
	l.Debug("Pgstore FindUrl")

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
		err = fmt.Errorf("small url not found")
		l.Error("",
			zap.Error(err))
		return nil, err
	}
	l.Debug("Pgstore find url successfull")
	return &smurlentity.Smurl{
		SmallURL: ses.SmallURL,
		LongURL:  pgs.LongURL,
		Count:    pgs.Count,
		IPInfo:   pgs.IPInfo,
	}, nil
}
