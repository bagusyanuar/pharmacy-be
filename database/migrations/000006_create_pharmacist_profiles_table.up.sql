-- 000006_create_pharmacist_profiles_table.up.sql
CREATE TABLE IF NOT EXISTS pharmacist_profiles (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    staff_id                UUID NOT NULL REFERENCES staff_profiles(id) ON DELETE RESTRICT,
    license_type            VARCHAR(30) NOT NULL, -- SIPA (Apoteker) atau STRTTK / SIPTTK (TTK)
    license_number          VARCHAR(100) NOT NULL,
    license_expired_date    DATE NOT NULL,
    is_apa                  BOOLEAN NOT NULL DEFAULT FALSE,
    assigned_branch_id      UUID NULL REFERENCES branches(id) ON DELETE RESTRICT,

    -- Universal Audit Trail
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by              UUID NULL REFERENCES users(id) ON DELETE RESTRICT,
    deleted_at              TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_pharmacist_staff ON pharmacist_profiles(staff_id) 
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_pharmacist_license_expiry ON pharmacist_profiles(license_expired_date) 
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_pharmacist_assigned_branch ON pharmacist_profiles(assigned_branch_id) 
WHERE deleted_at IS NULL;

COMMENT ON TABLE pharmacist_profiles IS 'Profil legalitas profesi kefarmasian untuk keabsahan Surat Pesanan (SP) dan penanggung jawab etiket obat.';
