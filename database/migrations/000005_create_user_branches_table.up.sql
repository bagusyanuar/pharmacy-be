-- 000005_create_user_branches_table.up.sql
CREATE TABLE IF NOT EXISTS user_branches (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    branch_id           UUID NOT NULL REFERENCES branches(id) ON DELETE RESTRICT,
    is_default          BOOLEAN NOT NULL DEFAULT FALSE,
    can_operate_pos     BOOLEAN NOT NULL DEFAULT TRUE,

    -- Universal Audit Trail
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID NULL REFERENCES users(id) ON DELETE RESTRICT,
    deleted_at          TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_user_branches_assignment ON user_branches(user_id, branch_id) 
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_user_branches_branch ON user_branches(branch_id) 
WHERE deleted_at IS NULL;

COMMENT ON TABLE user_branches IS 'Pemetaan penugasan staf ke cabang tertentu untuk isolasi multi-cabang (X-Branch-Id).';
