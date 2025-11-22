-- name: CreateUserPreferences :one
INSERT INTO user_preferences (
    user_id,
    analytics_enabled,
    error_tracking_enabled,
    marketing_enabled
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetUserPreferencesByID :one
SELECT * FROM user_preferences
WHERE id = $1;

-- name: GetUserPreferencesByUserID :one
SELECT * FROM user_preferences
WHERE user_id = $1;

-- name: UpdateUserPreferences :one
UPDATE user_preferences
SET
    analytics_enabled = $2,
    error_tracking_enabled = $3,
    marketing_enabled = $4,
    updated_at = NOW()
WHERE user_id = $1
RETURNING *;

-- name: DeleteUserPreferences :exec
DELETE FROM user_preferences
WHERE user_id = $1;

-- name: CountUserPreferences :one
SELECT COUNT(*) FROM user_preferences;

-- name: ListUsersWithAnalyticsEnabled :many
SELECT user_id FROM user_preferences
WHERE analytics_enabled = true;

-- name: ListUsersWithMarketingEnabled :many
SELECT user_id FROM user_preferences
WHERE marketing_enabled = true;
