package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"github.com/sinfirst/Test-Task-For-Effective-Mobile/internal/config"
	"go.uber.org/zap"
)

type PGDB struct {
	logger zap.SugaredLogger
	db     *pgxpool.Pool
}

func NewPGDB(conf config.Config, logger zap.SugaredLogger) *PGDB {
	db, err := pgxpool.New(context.Background(), conf.DatabaseDsn)

	if err != nil {
		logger.Errorw("Problem with connecting to db: ", err)
		return nil
	}

	err = db.Ping(context.Background())

	if err != nil {
		logger.Errorw("Problem with ping to db: ", err)
		return nil
	}

	return &PGDB{logger: logger, db: db}
}

func (p *PGDB) CreateInDB(ctx context.Context) {

}
func (p *PGDB) ReadFromDB(ctx context.Context) {

}
func (p *PGDB) UpdateInDB(ctx context.Context) {

}
func (p *PGDB) DeleteFromDB(ctx context.Context) {

}
func (p *PGDB) ListFromDB(ctx context.Context) {

}
func (p *PGDB) CostSumSubFromDB(ctx context.Context) {

}

func InitMigrations(conf config.Config, logger zap.SugaredLogger) error {
	if conf.DatabaseDsn == "" {
		return fmt.Errorf("DB url is not set")
	}

	db, err := sql.Open("pgx", conf.DatabaseDsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.Up(db, "internal/storage/migrations"); err != nil {
		return err
	}

	logger.Infow("Migrations applied successfully")
	return nil
}
