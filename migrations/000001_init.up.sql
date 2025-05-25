BEGIN;

CREATE TABLE "user" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "login" TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE user_order (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES "user"(id),
    "number" TEXT UNIQUE NOT NULL,
    status TEXT NOT NULL,
    accrual BIGINT,
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE user_balance (
    user_id UUID PRIMARY KEY REFERENCES "user"(id),
    "current" BIGINT NOT NULL DEFAULT 0,
    withdrawn BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE user_withdrawal (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES "user"(id),
    order_number TEXT NOT NULL,
    "sum" BIGINT NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX idx_user_order_user_id ON user_order(user_id);
CREATE INDEX idx_user_withdrawal_user_id ON user_withdrawal(user_id);

COMMIT;