-- 000002_create_roles_table.up.sql
CREATE TABLE IF NOT EXISTS roles (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_code               VARCHAR(30) NOT NULL,
    role_name               VARCHAR(100) NOT NULL,
    description             VARCHAR(255) NULL,
    can_supervisor_override BOOLEAN NOT NULL DEFAULT FALSE,

    -- Universal Audit Trail
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by              UUID NULL,
    deleted_at              TIMESTAMPTZ NULL,

    CONSTRAINT uq_roles_code UNIQUE (role_code)
);

CREATE INDEX IF NOT EXISTS idx_roles_can_override ON roles(can_supervisor_override) WHERE deleted_at IS NULL;

COMMENT ON TABLE roles IS 'Matriks hak akses 6 persona apotek terstandarisasi.';
