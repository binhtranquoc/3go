-- name: CreateSurchargeCondition :one
INSERT INTO system_surcharge_conditions (code, name, condition_type, config, is_active)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetSurchargeConditionByID :one
SELECT * FROM system_surcharge_conditions
WHERE id = $1;

-- name: ListSurchargeConditions :many
SELECT * FROM system_surcharge_conditions
ORDER BY code ASC;

-- name: UpdateSurchargeCondition :one
UPDATE system_surcharge_conditions
SET
  code = $2,
  name = $3,
  condition_type = $4,
  config = $5,
  is_active = $6
WHERE id = $1
RETURNING *;

-- name: DeleteSurchargeCondition :exec
DELETE FROM system_surcharge_conditions
WHERE id = $1;

-- name: GetSurchargeConditionByCode :one
SELECT * FROM system_surcharge_conditions
WHERE code = $1;

