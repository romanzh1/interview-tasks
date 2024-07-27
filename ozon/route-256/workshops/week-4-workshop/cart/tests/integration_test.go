package tests

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"week-4-workshop/cart/internal/domain"
	"week-4-workshop/cart/internal/repository/db_cart_repo"
)

func TestAddItemDB(t *testing.T) {
	ctx := context.Background()

	const dbConnEnv = "DB_CONN"

	dbConnStr := os.Getenv(dbConnEnv)

	conn, err := pgx.Connect(ctx, dbConnStr)
	require.NoError(t, err)

	storage := db_cart_repo.NewStorage(conn)

	item := domain.Item{
		SKU:   1000,
		Count: 10,
	}

	err = storage.AddItem(ctx, 1, item)
	require.NoError(t, err)
}
