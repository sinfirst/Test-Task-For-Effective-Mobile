package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"github.com/sinfirst/Test-Task-For-Effective-Mobile/config"
	"github.com/sinfirst/Test-Task-For-Effective-Mobile/internal/handlers"
	"github.com/sinfirst/Test-Task-For-Effective-Mobile/internal/models"
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

func (p *PGDB) CreateInDB(ctx context.Context, sub models.SubJSON) (int, error) {
	var id int
	query := `
		INSERT INTO subs (name_service, cost_per_month, user_uuid, date_start, date_end)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := p.db.QueryRow(ctx, query, sub.ServiceName, sub.Price, sub.UserUUID, sub.StartDate, sub.EndDate).Scan(&id)
	if err != nil {
		p.logger.Errorw("Problem with create in db: ", err)
		return 0, err
	}

	return id, nil
}

func (p *PGDB) ReadFromDB(ctx context.Context, id string) (models.SubJSON, error) {
	var sub models.SubJSON

	query := `SELECT (name_service, cost_per_month, user_uuid, date_start, date_end) FROM subs WHERE id = $1`
	row := p.db.QueryRow(ctx, query, id)
	err := row.Scan(&sub.ServiceName, &sub.Price, &sub.UserUUID, &sub.StartDate, &sub.EndDate)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.SubJSON{}, fmt.Errorf("not found")
		}
		p.logger.Errorw("Problem with read from db: ", err)
		return models.SubJSON{}, err
	}
	return sub, nil
}

func (p *PGDB) UpdateInDB(ctx context.Context, sub models.SubJSON) error {
	exist, err := p.checkSubExistByID(ctx, sub.ID)
	if err != nil {
		p.logger.Errorw("Problem with check sub exist in db: ", err)
		return err
	}

	if !exist {
		return fmt.Errorf("not found")
	}

	query := "UPDATE subs SET "
	args := []interface{}{}
	paramCount := 0

	if sub.ServiceName != "" {
		paramCount++
		query += fmt.Sprintf("name_service = $%d, ", paramCount)
		args = append(args, sub.ServiceName)
	}

	if sub.Price != 0 {
		paramCount++
		query += fmt.Sprintf("cost_per_month = $%d, ", paramCount)
		args = append(args, sub.Price)
	}

	if sub.UserUUID != "" {
		paramCount++
		query += fmt.Sprintf("user_uuid = $%d, ", paramCount)
		args = append(args, sub.UserUUID)
	}

	if sub.StartDate != "" {
		paramCount++
		query += fmt.Sprintf("date_start = $%d, ", paramCount)
		args = append(args, sub.StartDate)
	}

	if sub.EndDate != "" {
		paramCount++
		query += fmt.Sprintf("date_end = $%d, ", paramCount)
		args = append(args, sub.EndDate)
	}

	query = strings.TrimSuffix(query, ", ")
	paramCount++
	query += fmt.Sprintf(" WHERE id = $%d", paramCount)
	args = append(args, sub.ID)

	result, err := p.db.Exec(ctx, query, args...)
	if err != nil {
		p.logger.Errorw("Problem with update in db: ", err)
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("not found")
	}

	return nil
}

func (p *PGDB) DeleteFromDB(ctx context.Context, id string) error {
	exist, err := p.checkSubExistByID(ctx, id)
	if err != nil {
		p.logger.Errorw("Problem with check sub exist in db: ", err)
		return err
	}

	if !exist {
		return fmt.Errorf("not found")
	}

	query := `DELETE FROM subs
				WHERE id = $1`

	_, err = p.db.Exec(ctx, query, id)

	if err != nil {
		p.logger.Errorw("Problem with deleting from db: ", err)
		return err
	}
	return nil
}

func (p *PGDB) ListFromDB(ctx context.Context, id string) ([]models.SubJSON, error) {
	var subs []models.SubJSON

	query := `SELECT (id, name_service, cost_per_month, date_start, date_end) FROM subs WHERE user_uuid = $1`
	rows, err := p.db.Query(ctx, query, id)

	if err != nil {
		p.logger.Errorw("Problem with create list from db: ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var sub models.SubJSON

		err := rows.Scan(&sub.ServiceName, &sub.Price, &sub.UserUUID, &sub.StartDate, &sub.EndDate)
		if err != nil {
			p.logger.Errorw("Problem with create list from db: ", err)
			return nil, err
		}

		subs = append(subs, sub)
	}

	return subs, nil
}

func (p *PGDB) CostSumSubFromDB(ctx context.Context, req models.SubJSON) (int, error) {
	var totalCost int
	query := `
        SELECT COALESCE(SUM(price), 0) as total_cost
        FROM subs
        WHERE date_start <= $1 
        AND (end_date IS NULL OR end_date >= $2)
    `
	endTime, startTime, err := handlers.DateParse(req.StartDate, req.EndDate)
	if err != nil {
		p.logger.Errorw("Problem with parse date: ", err)
		return 0, err
	}

	args := []interface{}{endTime, startTime}
	argCounter := 3

	if req.UserUUID != "" {
		query += fmt.Sprintf(" AND user_id = $%d", argCounter)
		args = append(args, req.UserUUID)
		argCounter++
	}

	if req.ServiceName != "" {
		query += fmt.Sprintf(" AND service_name = $%d", argCounter)
		args = append(args, req.ServiceName)
		argCounter++
	}

	err = p.db.QueryRow(ctx, query, args...).Scan(&totalCost)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate total cost: %w", err)
	}

	return totalCost, nil
}

func (p *PGDB) checkSubExistByID(ctx context.Context, id string) (bool, error) {
	var exists bool

	err := p.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM subs WHERE id = $1
		)
	`, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking user existence: %w", err)
	}
	return exists, nil
}

func InitMigrations(conf config.Config, logger zap.SugaredLogger) error {
	if conf.DatabaseDsn == "" {
		return fmt.Errorf("DBDsn isn't set")
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
