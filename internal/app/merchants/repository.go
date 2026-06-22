package merchants

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) InsertMerchant(ctx context.Context, merchant Merchant) (Merchant, error) {
	query := `
		INSERT INTO merchants 
			(name, active, created_at, updated_at)
		VALUES
			(@name, @active, NOW(), NOW())
		RETURNING id, name, active, created_at, updated_at;`

	args := pgx.NamedArgs{
		"name":   merchant.Name,
		"active": merchant.Active,
	}

	var newMerchant Merchant

	err := r.pool.QueryRow(ctx, query, args).Scan(
		&newMerchant.ID,
		&newMerchant.Name,
		&newMerchant.Active,
		&newMerchant.CreatedAt,
		&newMerchant.UpdatedAt,
	)
	if err != nil {
		return Merchant{}, fmt.Errorf("Unable to insert row: %w", err)
	}

	return newMerchant, nil
}

func (r *Repository) ListMerchants(ctx context.Context) ([]Merchant, error) {
	query := `
		SELECT id, name, active, created_at, updated_at
		FROM merchants
		ORDER BY created_at DESC;
	`

	rows, err := r.pool.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("Unable to query merchants: %w", err)
	}
	defer rows.Close()

	merchants := []Merchant{}

	for rows.Next() {
		var merchant Merchant

		err := rows.Scan(
			&merchant.ID,
			&merchant.Name,
			&merchant.Active,
			&merchant.CreatedAt,
			&merchant.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("Unable to scan merchant: %w", err)
		}
		merchants = append(merchants, merchant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("merchant rows error: %w", err)
	}
	return merchants, nil
}

func (r *Repository) Merchant(ctx context.Context, merchantId int64) (Merchant, error) {
	query := `
		SELECT id, name, active, created_at, updated_at
		FROM merchants
		WHERE id = @merchantId
	`

	args := pgx.NamedArgs{
		"merchantId": merchantId,
	}

	var merchant Merchant

	err := r.pool.QueryRow(ctx, query, args).Scan(
		&merchant.ID,
		&merchant.Name,
		&merchant.Active,
		&merchant.CreatedAt,
		&merchant.UpdatedAt,
	)

	if err != nil {
		return Merchant{}, fmt.Errorf("unable to find merchant %d: %w", merchantId, err)
	}

	return merchant, nil
}

func (r *Repository) PatchMerchant(ctx context.Context, merchantId int64, update UpdateMerchantRequest) (Merchant, error) {
	var setClauses []string
	var args []interface{}
	argPosition := 1

	if update.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argPosition))
		args = append(args, *update.Name)
		argPosition++
	}

	if update.Active != nil {
		setClauses = append(setClauses, fmt.Sprintf("active = $%d", argPosition))
		args = append(args, *update.Active)
		argPosition++
	}

	if len(setClauses) == 0 {
		return Merchant{}, fmt.Errorf("no fields to update")
	}
	
	setClauses = append(setClauses, "updated_at = NOW()")
	args = append(args, merchantId)

	query := fmt.Sprintf(
		"UPDATE merchants SET  %s where id = $%d RETURNING id, name, active, created_at, updated_at",
		strings.Join(setClauses, ", "),
		argPosition,
	)

	var merchant Merchant

	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&merchant.ID,
		&merchant.Name,
		&merchant.Active,
		&merchant.CreatedAt,
		&merchant.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Merchant{}, fmt.Errorf("merchant with id %d not found: %w", merchantId, err)
		}

		return Merchant{}, fmt.Errorf("failed to update merchant: %w", err)
	}

	return merchant, nil
}
