package seeders

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bagusyanuar/pharmacy-be/pkg/password"
)

// SeedAuthAndStaff populates default roles, branch, super admin user, staff profile, and branch assignment.
func SeedAuthAndStaff(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Roles
	rolesSQL := `
	INSERT INTO roles (id, role_code, role_name, description, can_supervisor_override) VALUES
	('a0000000-0000-4000-8000-000000000001', 'SUPER_ADMIN', 'Super Administrator IT', 'Akses penuh ke konfigurasi sistem dan seluruh cabang.', TRUE),
	('a0000000-0000-4000-8000-000000000002', 'APOTEKER', 'Apoteker Pengelola Apotek (APA)', 'Penanggung jawab teknis kefarmasian, SP PBF, verifikasi resep & override kasir.', TRUE),
	('a0000000-0000-4000-8000-000000000003', 'ASISTEN_APOTEKER', 'Tenaga Teknis Kefarmasian (TTK)', 'Pelayanan resep dokter, peracikan obat, dan etiket.', FALSE),
	('a0000000-0000-4000-8000-000000000004', 'KASIR', 'Kasir Front-Office', 'Pelayanan kasir transaksi obat bebas (OTC) dan penerimaan pembayaran.', FALSE),
	('a0000000-0000-4000-8000-000000000005', 'GUDANG', 'Staf Gudang & Logistik', 'Penerimaan faktur barang PBF, penataan rak obat, mutasi cabang dan stock opname.', FALSE),
	('a0000000-0000-4000-8000-000000000006', 'FINANCE_OWNER', 'Owner & Manajemen Keuangan', 'Pengawasan keuangan, arus kas shift kasir, pelunasan hutang PBF, dan laporan laba rugi.', TRUE)
	ON CONFLICT (role_code) DO NOTHING;`
	if _, err := tx.ExecContext(ctx, rolesSQL); err != nil {
		return fmt.Errorf("failed to seed roles: %w", err)
	}

	// 2. Default Branch
	branchSQL := `
	INSERT INTO branches (id, branch_code, branch_name, branch_type, sia_number, sia_expired_date, phone, email, address, city, province, postal_code, is_active) VALUES
	('b0000000-0000-4000-8000-000000000001', 'AP-MLW-01', 'Apotek Sehat - Cabang Melawai', 'RETAIL_PHARMACY', '503/SIA-001/DPMPTSP/2024', '2029-01-15', '021-7201234', 'melawai@apoteksehat.id', 'Jl. Melawai Raya No. 12, Kebayoran Baru', 'Jakarta Selatan', 'DKI Jakarta', '12160', TRUE)
	ON CONFLICT (branch_code) DO NOTHING;`
	if _, err := tx.ExecContext(ctx, branchSQL); err != nil {
		return fmt.Errorf("failed to seed default branch: %w", err)
	}

	// 3. Super Admin User with Bcrypt cost >= 12
	defaultPassword := "AdminSuper2026!"
	passwordHash, err := password.HashWithCost(defaultPassword, 12)
	if err != nil {
		return fmt.Errorf("failed to hash default admin password: %w", err)
	}

	defaultPIN := "123456"
	pinHash, err := password.HashWithCost(defaultPIN, 10)
	if err != nil {
		return fmt.Errorf("failed to hash default admin PIN: %w", err)
	}

	userSQL := `
	INSERT INTO users (id, email, password_hash, pin_hash, barcode_card, role_id, is_active) VALUES
	('c0000000-0000-4000-8000-000000000001', 'admin@apoteksehat.id', $1, $2, 'CARD-SA-01', 'a0000000-0000-4000-8000-000000000001', TRUE)
	ON CONFLICT (email) WHERE deleted_at IS NULL DO NOTHING;`
	if _, err := tx.ExecContext(ctx, userSQL, passwordHash, pinHash); err != nil {
		return fmt.Errorf("failed to seed super admin user: %w", err)
	}

	// 4. Staff Profile for Super Admin (BR-PHARM-AUTH-13)
	staffSQL := `
	INSERT INTO staff_profiles (id, user_id, employee_code, nik, full_name, gender, phone, email, address, hire_date, job_position, is_active) VALUES
	('d0000000-0000-4000-8000-000000000001', 'c0000000-0000-4000-8000-000000000001', 'STF-ADM-001', '3171010101900001', 'Super Administrator IT', 'LAKI_LAKI', '081234567890', 'admin@apoteksehat.id', 'Kantor Pusat Apotek Sehat', '2024-01-01', 'ADMIN', TRUE)
	ON CONFLICT (employee_code) WHERE deleted_at IS NULL DO NOTHING;`
	if _, err := tx.ExecContext(ctx, staffSQL); err != nil {
		return fmt.Errorf("failed to seed super admin staff profile: %w", err)
	}

	// 5. User Branch Assignment
	userBranchSQL := `
	INSERT INTO user_branches (id, user_id, branch_id, is_default, can_operate_pos) VALUES
	('e0000000-0000-4000-8000-000000000001', 'c0000000-0000-4000-8000-000000000001', 'b0000000-0000-4000-8000-000000000001', TRUE, TRUE)
	ON CONFLICT (user_id, branch_id) WHERE deleted_at IS NULL DO NOTHING;`
	if _, err := tx.ExecContext(ctx, userBranchSQL); err != nil {
		return fmt.Errorf("failed to seed super admin branch assignment: %w", err)
	}

	return tx.Commit()
}
