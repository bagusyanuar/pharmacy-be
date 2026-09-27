-- 000004_create_staff_profiles_table.up.sql
CREATE TABLE IF NOT EXISTS staff_profiles (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NULL REFERENCES users(id) ON DELETE SET NULL,
    employee_code       VARCHAR(30) NOT NULL,
    nik                 VARCHAR(16) NULL,
    full_name           VARCHAR(100) NOT NULL,
    gender              VARCHAR(10) NULL, -- LAKI_LAKI / PEREMPUAN
    phone               VARCHAR(30) NOT NULL,
    email               VARCHAR(100) NULL,
    address             TEXT NULL,
    hire_date           DATE NOT NULL DEFAULT CURRENT_DATE,
    job_position        VARCHAR(50) NOT NULL, -- APOTEKER, ASISTEN_APOTEKER, KASIR, GUDANG, KURIR, FINANCE, ADMIN
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,

    -- Universal Audit Trail
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by          UUID NULL REFERENCES users(id) ON DELETE RESTRICT,
    deleted_at          TIMESTAMPTZ NULL
);

-- Partial Unique Indexes
CREATE UNIQUE INDEX IF NOT EXISTS uq_staff_employee_code ON staff_profiles(employee_code) 
WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uq_staff_user_id ON staff_profiles(user_id) 
WHERE deleted_at IS NULL AND user_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uq_staff_nik ON staff_profiles(nik) 
WHERE deleted_at IS NULL AND nik IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_staff_job_position ON staff_profiles(job_position) 
WHERE deleted_at IS NULL;

COMMENT ON TABLE staff_profiles IS 'Master data profil kepegawaian fisik (HR Record) apotek sesuai BR-PHARM-AUTH-13.';
