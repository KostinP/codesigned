DROP TABLE IF EXISTS user_inactivity_timeout;
DROP INDEX IF EXISTS idx_user_sessions_user_id;
DROP INDEX IF EXISTS idx_user_sessions_last_active_at;
DROP INDEX IF EXISTS idx_user_sessions_token;
DROP TABLE IF EXISTS user_sessions;
DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS abac_policies;
DROP INDEX IF EXISTS idx_abac_policies_target;
DROP INDEX IF EXISTS idx_abac_policies_priority;