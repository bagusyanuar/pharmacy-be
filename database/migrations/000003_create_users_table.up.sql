-- 000003_create_users_table.up.sql
CREATE TABLE IF NOT EXISTS users (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email               VARCHAR(100) NOT NULL,
    password_hash       VARCHAR(255) NOT NULL,
    pin_hash            VARCHAR(255) NULL,
    barcode_card        VARCHAR(50) NULL,
    role_id             UUID NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,

    -- Universal Audit Trail
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID NULL REFERENCES users(id) ON DELETE RESTRICT,
    deleted_at          TIMESTAMPTZ NULL
);

-- Partial Unique Indexes
CREATE UNIQUE INDEX IF NOT EXISTS uq_users_email ON users(email) 
WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uq_users_barcode_card ON users(barcode_card) 
WHERE deleted_at IS NULL AND barcode_card IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id) 
WHERE deleted_at IS NULL;

COMMENT ON TABLE users IS 'Akun identitas digital IAM untuk akses login sistem apotek.';
