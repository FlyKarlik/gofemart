BEGIN;

DROP INDEX IF EXISTS idx_user_withdrawal_user_id;
DROP INDEX IF EXISTS idx_user_order_user_id;


DROP TABLE IF EXISTS user_withdrawal;
DROP TABLE IF EXISTS user_balance;
DROP TABLE IF EXISTS user_order;
DROP TABLE IF EXISTS "user";

COMMIT;
