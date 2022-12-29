package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sanyarise/smurl/internal/models"
	"github.com/sanyarise/smurl/internal/usecase"
	"go.uber.org/zap"
)

var _ usecase.SmurlStore = &SmurlRepository{}

type Smurl struct {
	SmallURL   string
	CreatedAt  time.Time
	ModifiedAt time.Time
	LongURL    string
	AdminURL   string
	IPInfo     []string
	Count      uint64
}

type SmurlRepository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewSmurlRepository(dns string, logger *zap.Logger) (*SmurlRepository, error) {
	logger.Debug("Enter in pgstore func NewSmurlStore()")

	conf, err := pgxpool.ParseConfig(dns)
	if err != nil {
		logger.Sugar().Errorf("can't init storage: %s", err)
		return nil, fmt.Errorf("can't init storage: %w", err)
	}
	db, err := pgxpool.ConnectConfig(context.Background(), conf)
	if err != nil {
		logger.Sugar().Errorf("can't create pool %s", err)
		return nil, fmt.Errorf("can't create pool %w", err)
	}
	_, err = db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS smurls (
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
	repository := &SmurlRepository{
		db:     db,
		logger: logger,
	}
	return repository, nil
}

func (repo *SmurlRepository) Close() {
	repo.logger.Debug("Enter in repository Close()")
	repo.db.Close()
}

// CreateURL saving long url, short url and admin url to database
func (repo *SmurlRepository) Create(ctx context.Context, smurl models.Smurl) (*models.Smurl, error) {
	repo.logger.Debug("Enter in repository CreateURL()")
	repositorySmurl := &Smurl{
		LongURL:    smurl.LongURL,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
		SmallURL:   smurl.SmallURL,
		AdminURL:   smurl.AdminURL,
		Count:      0,
		IPInfo:     []string{},
	}
	// Starting a transaction to write data to the database
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		repo.logger.Error("error on begin transaction",
			zap.Error(err))
	}
	// Write to database
	_, err = tx.Exec(ctx, `INSERT INTO smurls
	(small_url, created_at, modified_at, long_url, admin_url, count, ip_info)
	values ($1, $2, $3, $4, $5, $6, $7)`,
		repositorySmurl.SmallURL,
		repositorySmurl.CreatedAt,
		repositorySmurl.ModifiedAt,
		repositorySmurl.LongURL,
		repositorySmurl.AdminURL,
		repositorySmurl.Count,
		repositorySmurl.IPInfo,
	)
	if err != nil {
		//Return to original value in case of unsuccessful write
		tx.Rollback(ctx)
		repo.logger.Error("error on insert values into table",
			zap.Error(err))
		return nil, err
	}
	// End of transaction
	tx.Commit(ctx)
	repo.logger.Debug("Pgstore create smurl successfull")
	// Return object with short and admin url
	return &models.Smurl{
		SmallURL: repositorySmurl.SmallURL,
		AdminURL: repositorySmurl.AdminURL,
	}, nil
}

// UpdateStat updating statistics data when clicking on a reduced url
func (repo *SmurlRepository) UpdateStat(ctx context.Context, smurl models.Smurl) error {
	repo.logger.Debug("Enter in repository UpdateStat()")
	repositorySmurl := &Smurl{
		ModifiedAt: time.Now(),
		Count:      smurl.Count,
		IPInfo:     smurl.IPInfo,
		SmallURL:   smurl.SmallURL,
	}
	// Starting a transaction to write updated data
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		repo.logger.Error("error on begin transaction",
			zap.Error(err))
	}
	// Write updated data
	_, err = tx.Exec(ctx, `UPDATE smurls SET modified_at = $1, count = $2, ip_info = $3
	WHERE small_url = $4`, repositorySmurl.ModifiedAt, repositorySmurl.Count, repositorySmurl.IPInfo, repositorySmurl.SmallURL)
	if err != nil {
		repo.logger.Error("error on update values into table",
			zap.Error(err))

		// Return to original value in case of unsuccessful write
		tx.Rollback(ctx)
		return err
	}
	// End of transaction
	tx.Commit(ctx)
	repo.logger.Debug("Pgstore update stat successfull")

	return nil
}

// ReadStat reads statistics data
func (repo *SmurlRepository) ReadStat(ctx context.Context, adminUrl string) (*models.Smurl, error) {
	repo.logger.Debug("Enter in pgstore func ReadStat()")
	repositorySmurl := &Smurl{}
	// Performing a database search
	rows, err := repo.db.Query(ctx,
		`SELECT small_url, created_at, modified_at, long_url, admin_url, count, ip_info FROM smurls
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
			&repositorySmurl.CreatedAt,
			&repositorySmurl.ModifiedAt,
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
	if repositorySmurl.AdminURL == "" {
		return nil, models.ErrNotFound
	}
	result := &models.Smurl{
		SmallURL:   repositorySmurl.SmallURL,
		CreatedAt:  repositorySmurl.CreatedAt,
		ModifiedAt: repositorySmurl.ModifiedAt,
		LongURL:    repositorySmurl.LongURL,
		AdminURL:   repositorySmurl.AdminURL,
		Count:      repositorySmurl.Count,
		IPInfo:     repositorySmurl.IPInfo,
	}
	repo.logger.Debug("Pgstore read stat successfull")

	return result, nil
}

// FindURL search small url in database
func (repo *SmurlRepository) FindURL(ctx context.Context, smallUrl string) (*models.Smurl, error) {
	repo.logger.Debug("Enter in pgstore func FindUrl()")

	repositorySmurl := Smurl{}
	row := repo.db.QueryRow(ctx,
		`SELECT small_url, created_at, modified_at, long_url, count, ip_info FROM smurls WHERE small_url = $1`, smallUrl)
	if err := row.Scan(
		&repositorySmurl.SmallURL,
		&repositorySmurl.CreatedAt,
		&repositorySmurl.ModifiedAt,
		&repositorySmurl.LongURL,
		&repositorySmurl.Count,
		&repositorySmurl.IPInfo,
	); err != nil {
		repo.logger.Error("error find small url",
			zap.Error(err))
		return nil, err
	}

	if repositorySmurl.SmallURL == "" {
		err := fmt.Errorf("small url not found")
		repo.logger.Error("",
			zap.Error(err))
		return nil, models.ErrNotFound
	}
	repo.logger.Debug("URL find successfull")
	return &models.Smurl{
		SmallURL:   repositorySmurl.SmallURL,
		CreatedAt:  repositorySmurl.CreatedAt,
		ModifiedAt: repositorySmurl.ModifiedAt,
		LongURL:    repositorySmurl.LongURL,
		Count:      repositorySmurl.Count,
		IPInfo:     repositorySmurl.IPInfo,
	}, nil
}
