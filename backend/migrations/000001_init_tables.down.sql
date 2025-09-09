-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_swap_requests_updated_at ON swap_requests;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_users_user_id;
DROP INDEX IF EXISTS idx_users_roll_no;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_mess;
DROP INDEX IF EXISTS idx_swap_requests_requester_id;
DROP INDEX IF EXISTS idx_swap_requests_status;
DROP INDEX IF EXISTS idx_swap_requests_type;

-- Drop tables
DROP TABLE IF EXISTS swap_requests;
DROP TABLE IF EXISTS users;

-- Drop extension (optional, might be used by other applications)
-- DROP EXTENSION IF EXISTS "uuid-ossp";
