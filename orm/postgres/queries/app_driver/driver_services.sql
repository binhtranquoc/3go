-- name: CreateDriverService :one
INSERT INTO driver_services (
    driver_id,
    service_id,
    status
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetDriverServicesByDriverID :many
SELECT *
FROM driver_services
WHERE driver_id = $1
  AND deleted_at IS NULL;

