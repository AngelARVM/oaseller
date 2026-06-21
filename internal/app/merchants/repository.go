package merchants

import (
	"context"
	"fmt"

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
