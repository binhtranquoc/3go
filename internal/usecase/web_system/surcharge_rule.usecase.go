package web_system

import (
	"context"
	"errors"

	common "go-structure/internal/common"
	dto_common "go-structure/internal/dto/common"
	dto "go-structure/internal/dto/web_system"
	"go-structure/internal/helper/database"
	"go-structure/internal/helper/parse"
	pgdb "go-structure/orm/db/postgres"
	account_repo "go-structure/internal/repository"
	"go-structure/internal/repository/model"
	websystem_model "go-structure/internal/repository/model/web_system"
	websystem_repo "go-structure/internal/repository/web_system"
	serviceTransformer "go-structure/internal/transformer/web_system"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrSurchargeRuleNotFound = errors.New("không tìm thấy quy tắc phụ thu")

type (
	ISurchargeRuleUsecase interface {
		Create(ctx context.Context, adminID uuid.UUID, req *dto.CreateSurchargeRuleRequestDto) (*dto.SurchargeRuleItemDto, error)
		GetByID(ctx context.Context, id uuid.UUID) (*dto.SurchargeRuleItemDto, error)
		List(ctx context.Context, serviceID, zoneID *uuid.UUID) (*dto.ListSurchargeRulesResponseDto, error)
		Update(ctx context.Context, adminID uuid.UUID, id uuid.UUID, req *dto.UpdateSurchargeRuleRequestDto) (*dto.SurchargeRuleItemDto, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	surchargeRuleUsecase struct {
		repo                   websystem_repo.ISurchargeRuleRepository
		conditionRepo          websystem_repo.ISurchargeRuleConditionRepository
		surchargeConditionRepo websystem_repo.ISurchargeConditionRepository
		serviceRepo            websystem_repo.IServiceRepository
		zoneRepo               account_repo.IZoneRepository
		transactionManager     database.TransactionManager
	}
)

func NewSurchargeRuleUsecase(
	repo websystem_repo.ISurchargeRuleRepository,
	conditionRepo websystem_repo.ISurchargeRuleConditionRepository,
	surchargeConditionRepo websystem_repo.ISurchargeConditionRepository,
	serviceRepo websystem_repo.IServiceRepository,
	zoneRepo account_repo.IZoneRepository,
	transactionManager database.TransactionManager,
) ISurchargeRuleUsecase {
	return &surchargeRuleUsecase{
		repo:                   repo,
		conditionRepo:          conditionRepo,
		surchargeConditionRepo: surchargeConditionRepo,
		serviceRepo:            serviceRepo,
		zoneRepo:               zoneRepo,
		transactionManager:     transactionManager,
	}
}

func (u *surchargeRuleUsecase) Create(ctx context.Context, adminID uuid.UUID, req *dto.CreateSurchargeRuleRequestDto) (*dto.SurchargeRuleItemDto, error) {
	if u.repo == nil {
		return nil, nil
	}
	serviceID, err := uuid.Parse(req.ServiceID)
	if err != nil {
		return nil, err
	}
	zoneID, err := uuid.Parse(req.ZoneID)
	if err != nil {
		return nil, err
	}
	if err := u.validateConditionIDs(ctx, req.ConditionIDs); err != nil {
		return nil, err
	}
	ruleToValidate := &websystem_model.SurchargeRule{
		Amount:   req.Amount,
		Unit:     req.Unit,
		Priority: int32(req.Priority),
	}
	if err := ruleToValidate.ValidateSurchargeRule(); err != nil {
		return nil, err
	}
	params := pgdb.CreateSurchargeRuleParams{
		ServiceID: serviceID,
		ZoneID:    zoneID,
		Amount:    common.Float64ToNumeric(req.Amount),
		Unit:      req.Unit,
		IsActive:  req.IsActive,
		Priority:  int32(req.Priority),
		CreatedBy: adminID,
		UpdatedBy: adminID,
	}

	rule, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (*websystem_model.SurchargeRule, error) {
			if u.serviceRepo != nil {
				if _, err := u.serviceRepo.GetServiceByID(txCtx, serviceID); err != nil {
					if errors.Is(err, pgx.ErrNoRows) {
						return nil, ErrServiceNotFound
					}
					return nil, err
				}
			}
			if u.zoneRepo != nil {
				if _, err := u.zoneRepo.GetZoneByID(txCtx, zoneID); err != nil {
					if errors.Is(err, pgx.ErrNoRows) {
						return nil, ErrZoneNotFound
					}
					return nil, err
				}
			}

			created, err := u.repo.Create(txCtx, params)
			if err != nil {
				return nil, err
			}

			// handle conditions pivot
			if u.conditionRepo != nil && len(req.ConditionIDs) > 0 {
				condUUIDs, err := parse.ParseUUIDStrings(req.ConditionIDs)
				if err != nil {
					return nil, err
				}
				if err := u.conditionRepo.AddConditions(txCtx, created.ID, condUUIDs); err != nil {
					return nil, err
				}
			}

			return created, nil
		},
	)
	if err != nil {
		return nil, err
	}

	item, err := u.enrichRuleWithConditions(ctx, rule)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (u *surchargeRuleUsecase) GetByID(ctx context.Context, id uuid.UUID) (*dto.SurchargeRuleItemDto, error) {
	if u.repo == nil {
		return nil, ErrSurchargeRuleNotFound
	}
	rule, err := u.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSurchargeRuleNotFound
		}
		return nil, err
	}
	item, err := u.enrichRuleWithConditions(ctx, rule)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (u *surchargeRuleUsecase) List(ctx context.Context, serviceID, zoneID *uuid.UUID) (*dto.ListSurchargeRulesResponseDto, error) {
	if u.repo == nil {
		return &dto.ListSurchargeRulesResponseDto{
			Items:      nil,
			Pagination: dto_common.PaginationMeta{Page: 1, Limit: 0, Total: 0},
		}, nil
	}
	rules, err := u.repo.List(ctx, serviceID, zoneID)
	if err != nil {
		return nil, err
	}
	items := make([]dto.SurchargeRuleItemDto, 0, len(rules))
	for _, r := range rules {
		item, err := u.enrichRuleWithConditions(ctx, r)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}
	return &dto.ListSurchargeRulesResponseDto{
		Items: items,
		Pagination: dto_common.PaginationMeta{
			Page:  1,
			Limit: len(items),
			Total: int64(len(items)),
		},
	}, nil
}

func (u *surchargeRuleUsecase) Update(ctx context.Context, adminID uuid.UUID, id uuid.UUID, req *dto.UpdateSurchargeRuleRequestDto) (*dto.SurchargeRuleItemDto, error) {
	if u.repo == nil {
		return nil, ErrSurchargeRuleNotFound
	}
	serviceID, err := uuid.Parse(req.ServiceID)
	if err != nil {
		return nil, err
	}
	zoneID, err := uuid.Parse(req.ZoneID)
	if err != nil {
		return nil, err
	}
	if err := u.validateConditionIDs(ctx, req.ConditionIDs); err != nil {
		return nil, err
	}
	// Entity business rule: validate trước khi persist (Clean Architecture)
	ruleToValidate := &websystem_model.SurchargeRule{
		Amount:   req.Amount,
		Unit:     req.Unit,
		Priority: int32(req.Priority),
	}
	if err := ruleToValidate.ValidateSurchargeRule(); err != nil {
		return nil, err
	}
	params := pgdb.UpdateSurchargeRuleParams{
		ID:        id,
		ServiceID: serviceID,
		ZoneID:    zoneID,
		Amount:    common.Float64ToNumeric(req.Amount),
		Unit:      req.Unit,
		IsActive:  req.IsActive,
		Priority:  int32(req.Priority),
		UpdatedBy: adminID,
	}

	rule, err := database.WithTransaction(
		u.transactionManager,
		ctx,
		func(txCtx context.Context) (*websystem_model.SurchargeRule, error) {
			if u.serviceRepo != nil {
				if _, err := u.serviceRepo.GetServiceByID(txCtx, serviceID); err != nil {
					if errors.Is(err, pgx.ErrNoRows) {
						return nil, ErrServiceNotFound
					}
					return nil, err
				}
			}
			if u.zoneRepo != nil {
				if _, err := u.zoneRepo.GetZoneByID(txCtx, zoneID); err != nil {
					if errors.Is(err, pgx.ErrNoRows) {
						return nil, ErrZoneNotFound
					}
					return nil, err
				}
			}

			updated, err := u.repo.Update(txCtx, params)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil, ErrSurchargeRuleNotFound
				}
				return nil, err
			}

			// reset & re-add conditions in pivot
			if u.conditionRepo != nil {
				if err := u.conditionRepo.DeleteBySurchargeID(txCtx, id); err != nil {
					return nil, err
				}
				if len(req.ConditionIDs) > 0 {
					condUUIDs, err := parse.ParseUUIDStrings(req.ConditionIDs)
					if err != nil {
						return nil, err
					}
					if err := u.conditionRepo.AddConditions(txCtx, id, condUUIDs); err != nil {
						return nil, err
					}
				}
			}

			return updated, nil
		},
	)
	if err != nil {
		return nil, err
	}
	item, err := u.enrichRuleWithConditions(ctx, rule)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (u *surchargeRuleUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if u.repo == nil {
		return ErrSurchargeRuleNotFound
	}
	return u.repo.Delete(ctx, id)
}

func (u *surchargeRuleUsecase) validateConditionIDs(ctx context.Context, conditionIDs []string) error {
	if len(conditionIDs) == 0 || u.surchargeConditionRepo == nil {
		return nil
	}
	condUUIDs, err := parse.ParseUUIDStrings(conditionIDs)
	if err != nil {
		return err
	}
	for _, id := range condUUIDs {
		_, err := u.surchargeConditionRepo.GetByID(ctx, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrSurchargeConditionNotFound
			}
			return err
		}
	}
	return nil
}

func (u *surchargeRuleUsecase) enrichRuleWithConditions(ctx context.Context, rule *websystem_model.SurchargeRule) (*dto.SurchargeRuleItemDto, error) {
	if rule == nil {
		return nil, nil
	}

	var (
		svc  *websystem_model.Service
		zone *model.Zone
	)

	if u.serviceRepo != nil {
		if s, err := u.serviceRepo.GetServiceByID(ctx, rule.ServiceID); err == nil {
			svc = s
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	if u.zoneRepo != nil {
		if z, err := u.zoneRepo.GetZoneByID(ctx, rule.ZoneID); err == nil {
			zone = z
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	if u.conditionRepo == nil || u.surchargeConditionRepo == nil {
		item := serviceTransformer.ToSurchargeRuleItemDtoWithConditions(rule, nil, nil, svc, zone)
		return &item, nil
	}

	ids, err := u.conditionRepo.GetConditionIDsBySurchargeID(ctx, rule.ID)
	if err != nil {
		return nil, err
	}

	conditions := make([]*websystem_model.SurchargeCondition, 0, len(ids))
	for _, id := range ids {
		cond, err := u.surchargeConditionRepo.GetByID(ctx, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				continue
			}
			return nil, err
		}
		conditions = append(conditions, cond)
	}

	item := serviceTransformer.ToSurchargeRuleItemDtoWithConditions(rule, ids, conditions, svc, zone)
	return &item, nil
}
