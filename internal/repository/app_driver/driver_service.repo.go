package app_driver

import (
	"context"

	"go-structure/internal/helper/database"
	pgdb "go-structure/orm/db/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	IDriverServiceRepository interface {
		SetDriverServices(ctx context.Context, driverID uuid.UUID, serviceIDs []uuid.UUID) error
	}

	driverServiceRepository struct {
		pool *pgxpool.Pool
	}
)

func NewDriverServiceRepository(pool *pgxpool.Pool) IDriverServiceRepository {
	return &driverServiceRepository{pool: pool}
}

func (r *driverServiceRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *driverServiceRepository) SetDriverServices(ctx context.Context, driverID uuid.UUID, serviceIDs []uuid.UUID) error {
	if len(serviceIDs) == 0 {
		return nil
	}
	db := r.getDB(ctx)
	for _, sid := range serviceIDs {
		arg := pgdb.CreateDriverServiceParams{
			DriverID:  driverID,
			ServiceID: sid,
			Status:    pgdb.DriverServiceStatusPENDINGDOCUMENT,
		}
		if _, err := db.CreateDriverService(ctx, arg); err != nil {
			return err
		}
	}
	return nil
}
