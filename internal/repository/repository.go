package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sanyarise/smurl/internal/models"
	"github.com/sanyarise/smurl/internal/usecase"

	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
)

var _ usecase.SmurlStore = &SmurlRepository{}

type Smurl struct {
	SmallURL  string
	CreatedAt time.Time
	LongURL   string
	AdminURL  string
	IPInfo    []string
	Count     uint64
}

type SmurlRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewSmurlRepository(dsn string, logger *zap.Logger) (*SmurlRepository, error) {
	logger.Debug("Enter in pgstore func NewSmurlStore()")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Error("error on sql open",
			zap.Error(err))
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		logger.Error("error on db ping",
			zap.Error(err))
		db.Close()
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS smurls (
		small_url varchar NOT NULL,
		created_at timestamptz NOT NULL,
		long_url varchar NOT NULL,
		admin_url varchar NOT NULL,
		count integer,
		ip_info text[]
		)`)
	if err != nil {
		logger.Error("error on create table",
			zap.Error(err))
		db.Close()
		return nil, err
	}
	su := &SmurlRepository{
		db:     db,
		logger: logger,
	}
	return su, nil
}

func (repo *SmurlRepository) Close() {
	repo.logger.Debug("Enter in repository Close()")
	repo.db.Close()
}

// CreateURL saving long url, short url and admin url to database
func (repo *SmurlRepository) Create(ctx context.Context, smurl models.Smurl) (*models.Smurl, error) {
	repo.logger.Debug("Enter in repository CreateURL()")
	repositorySmurl := &Smurl{
		LongURL:   smurl.LongURL,
		CreatedAt: time.Now(),
		SmallURL:  smurl.SmallURL,
		AdminURL:  smurl.AdminURL,
		Count:     0,
		IPInfo:    []string{},
	}
	// Starting a transaction to write data to the database
	tx, err := repo.db.Begin()
	if err != nil {
		repo.logger.Error("error on begin transaction",
			zap.Error(err))
	}
	// Write to database
	_, err = tx.ExecContext(ctx, `INSERT INTO smurls
	(small_url, created_at, long_url, admin_url, count, ip_info)
	values ($1, $2, $3, $4, $5, $6)`,
		repositorySmurl.SmallURL,
		repositorySmurl.CreatedAt,
		repositorySmurl.LongURL,
		repositorySmurl.AdminURL,
		repositorySmurl.Count,
		repositorySmurl.IPInfo,
	)
	if err != nil {
		//Return to original value in case of unsuccessful write
		tx.Rollback()
		repo.logger.Error("error on insert values into table",
			zap.Error(err))
		return nil, err
	}
	// End of transaction
	tx.Commit()
	repo.logger.Debug("Pgstore create smurl successfull")
	// Return object with short and admin url
	return &smurl.Smurl{
		SmallURL: repositorySmurl.SmallURL,
		AdminURL: repositorySmurl.AdminURL,
	}, nil
}

// UpdateStat updating statistics data when clicking on a reduced url
func (repo *SmurlRepository) UpdateStat(ctx context.Context, smurl models.Smurl) error {
	repo.logger.Debug("Enter in repository UpdateStat()")
	repositorySmurl := &Smurl{
		Count:    smurl.Count,
		IPInfo:   smurl.IPInfo,
		SmallURL: smurl.SmallURL,
	}
	// Starting a transaction to write updated data
	tx, err := repo.db.Begin()
	if err != nil {
		repo.logger.Error("error on begin transaction",
			zap.Error(err))
	}
	// Write updated data
	_, err = tx.ExecContext(ctx, `UPDATE smurls SET count = $1, ip_info = $2
	WHERE small_url = $3`, repositorySmurl.Count, repositorySmurl.IPInfo, repositorySmurl.SmallURL)
	if err != nil {
		repo.logger.Error("error on update values into table",
			zap.Error(err))

		// Return to original value in case of unsuccessful write
		tx.Rollback()
		return err
	}
	// End of transaction
	tx.Commit()
	repo.logger.Debug("Pgstore update stat successfull")

	return nil
}

// ReadStat reads statistics data
func (repo *SmurlRepository) ReadStat(ctx context.Context, adminUrl string) (*models.Smurl, error) {
	repo.logger.Debug("Enter in pgstore func ReadStat()")
	repositorySmurl := &Smurl{}
	// Performing a database search
	rows, err := repo.db.QueryContext(ctx,
		`SELECT small_url, long_url, admin_url, count, ip_info FROM smurls
	 WHERE admin_url = $1`, adminUrl)
	if err != nil {
		repo.logger.Error("error on query in table",
			zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&repositorySmurl.SmallURL,
			&repositorySmurl.LongURL,
			&repositorySmurl.AdminURL,
			&repositorySmurl.Count,
			&repositorySmurl.IPInfo,
		); err != nil {
			repo.logger.Error("error on rows scan",
				zap.Error(err))
			return nil, err
		}
	}
	err = fmt.Errorf("admin_url: %s not found", adminUrl)
	if repositorySmurl.AdminURL == "" {
		repo.logger.Error("",
			zap.Error(err))
		return nil, err
	}
	result := &models.Smurl{
		SmallURL: repositorySmurl.SmallURL,
		LongURL:  repositorySmurl.LongURL,
		AdminURL: repositorySmurl.AdminURL,
		Count:    repositorySmurl.Count,
		IPInfo:   repositorySmurl.IPInfo,
	}
	repo.logger.Debug("Pgstore read stat successfull")

	return result, nil
}

// FindURL search small url in database
func (repo *SmurlRepository) FindURL(ctx context.Context, smallUrl string) (*models.Smurl, error) {
	repo.logger.Debug("Enter in pgstore func FindUrl()")

	repositorySmurl := &Smurl{}
	rows, err := repo.db.QueryContext(ctx,
		`SELECT small_url, long_url, count, ip_info FROM smurls WHERE small_url = $1`, smallUrl)
	if err != nil {
		repo.logger.Error("error on search small url into table",
			zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			repositorySmurl.SmallURL,
			repositorySmurl.LongURL,
			repositorySmurl.Count,
			repositorySmurl.IPInfo,
		); err != nil {
			repo.logger.Error("error find small url",
				zap.Error(err))
			return nil, err
		}
	}
	if repositorySmurl.SmallURL == "" {
		err = fmt.Errorf("small url not found")
		repo.logger.Error("",
			zap.Error(err))
		return nil, err
	}
	repo.logger.Debug("URL find successfull")
	return &models.Smurl{
		SmallURL: models.SmallURL,
		LongURL:  repositorySmurl.LongURL,
		Count:    repositorySmurl.Count,
		IPInfo:   repositorySmurl.IPInfo,
	}, nil
}
