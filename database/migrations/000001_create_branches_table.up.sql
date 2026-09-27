-- 000001_create_branches_table.up.sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

DO $$ BEGIN
    CREATE TYPE branch_type_enum AS ENUM (
        'RETAIL_PHARMACY',
        'CLINIC_PHARMACY',
        'CENTRAL_WAREHOUSE'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS branches (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    branch_code         VARCHAR(30) NOT NULL,
    branch_name         VARCHAR(100) NOT NULL,
    branch_type         branch_type_enum NOT NULL DEFAULT 'RETAIL_PHARMACY',
    sia_number          VARCHAR(100) NULL,
    sia_expired_date    DATE NULL,
    phone               VARCHAR(30) NULL,
    email               VARCHAR(100) NULL,
    address             TEXT NOT NULL,
    city                VARCHAR(100) NOT NULL,
    province            VARCHAR(100) NOT NULL,
    postal_code         VARCHAR(10) NULL,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,

    -- Universal Audit Trail
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID NULL,
    deleted_at          TIMESTAMPTZ NULL,

    CONSTRAINT uq_branches_code UNIQUE (branch_code)
);

CREATE INDEX IF NOT EXISTS idx_branches_is_active ON branches(is_active) WHERE deleted_at IS NULL;

COMMENT ON TABLE branches IS 'Master data cabang apotek, gudang pusat, dan izin SIA resmi.';
