package repository

import (
	"context"

	"go-structure/internal/helper/database"
	pgdb "go-structure/orm/db/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	ISurchargeRuleConditionRepository interface {
		AddConditions(ctx context.Context, surchargeID uuid.UUID, conditionIDs []uuid.UUID) error
		DeleteBySurchargeID(ctx context.Context, surchargeID uuid.UUID) error
		GetConditionIDsBySurchargeID(ctx context.Context, surchargeID uuid.UUID) ([]uuid.UUID, error)
	}

	surchargeRuleConditionRepository struct {
		pool *pgxpool.Pool
	}
)

func NewSurchargeRuleConditionRepository(pool *pgxpool.Pool) ISurchargeRuleConditionRepository {
	return &surchargeRuleConditionRepository{pool: pool}
}

func (r *surchargeRuleConditionRepository) getDB(ctx context.Context) *pgdb.Queries {
	return database.GetQueries(ctx, r.pool)
}

func (r *surchargeRuleConditionRepository) AddConditions(ctx context.Context, surchargeID uuid.UUID, conditionIDs []uuid.UUID) error {
	if len(conditionIDs) == 0 {
		return nil
	}
	db := r.getDB(ctx)
	for _, condID := range conditionIDs {
		if err := db.InsertSurchargeRuleCondition(ctx, pgdb.InsertSurchargeRuleConditionParams{
			SurchargeID: surchargeID,
			ConditionID: condID,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *surchargeRuleConditionRepository) DeleteBySurchargeID(ctx context.Context, surchargeID uuid.UUID) error {
	db := r.getDB(ctx)
	return db.DeleteSurchargeRuleConditionsBySurchargeID(ctx, surchargeID)
}

func (r *surchargeRuleConditionRepository) GetConditionIDsBySurchargeID(ctx context.Context, surchargeID uuid.UUID) ([]uuid.UUID, error) {
	db := r.getDB(ctx)
	return db.GetConditionIDsBySurchargeID(ctx, surchargeID)
}

