-- name: InsertSurchargeRuleCondition :exec
INSERT INTO system_surcharge_rule_conditions (surcharge_id, condition_id)
VALUES ($1, $2);

-- name: DeleteSurchargeRuleConditionsBySurchargeID :exec
DELETE FROM system_surcharge_rule_conditions
WHERE surcharge_id = $1;

-- name: GetConditionIDsBySurchargeID :many
SELECT condition_id
FROM system_surcharge_rule_conditions
WHERE surcharge_id = $1;

