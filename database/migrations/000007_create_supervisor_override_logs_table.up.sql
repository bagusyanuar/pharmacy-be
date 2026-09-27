-- 000007_create_supervisor_override_logs_table.up.sql
CREATE TABLE IF NOT EXISTS supervisor_override_logs (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    branch_id           UUID NOT NULL REFERENCES branches(id) ON DELETE RESTRICT,
    cashier_user_id     UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    supervisor_user_id  UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    action_type         VARCHAR(50) NOT NULL, -- VOID_TRANSACTION, DELETE_ITEM, MANUAL_DISCOUNT, OPEN_DRAWER, PETTY_CASH
    reason_category     VARCHAR(50) NOT NULL, -- PATIENT_CANCEL, WRONG_INPUT, INSUFFICIENT_FUNDS, EMPLOYEE_DISCOUNT, OTHER
    reason_notes        TEXT NULL,
    target_entity_type  VARCHAR(50) NULL,     -- SALE_ORDER, SALE_ITEM, CASH_DRAWER
    target_entity_id    UUID NULL,

    -- Universal Audit Trail
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID NULL REFERENCES users(id) ON DELETE RESTRICT,
    deleted_at          TIMESTAMPTZ NULL
);

-- Performance Indexes for audit trail searches
CREATE INDEX IF NOT EXISTS idx_override_logs_branch_date ON supervisor_override_logs(branch_id, created_at DESC) 
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_override_logs_supervisor ON supervisor_override_logs(supervisor_user_id) 
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_override_logs_cashier ON supervisor_override_logs(cashier_user_id) 
WHERE deleted_at IS NULL;

COMMENT ON TABLE supervisor_override_logs IS 'Rekam jejak otorisasi supervisor di kasir POS untuk audit kepatuhan dan integritas kasir.';
