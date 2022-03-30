package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/GeneralKenobi/mailman/internal/db"
	"github.com/GeneralKenobi/mailman/pkg/mdctx"
	"github.com/GeneralKenobi/mailman/pkg/util"
)

// rowScanSupplier defines a function that returns an object and its properties that are then used in sql.Rows Scan method to convert
// columns of a single row into an object.
type rowScanSupplier[T any] func() (scanTarget *T, scanTargetProperties []any)

// selectingAll executes the query and converts obtained rows into a slice of T. It's not considered an error if the query returned no rows.
//
// queryName is only used in error messages.
func selectingAll[T any](ctx context.Context, queryName string, sql SqlExecutor, rowScanSupplier rowScanSupplier[T],
	query string, args ...any) ([]T, error) {

	rows, err := sql.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: error running query: %w", queryName, err)
	}
	defer closeRows(ctx, rows)

	var items []T
	for rows.Next() {
		item, props := rowScanSupplier()
		err = rows.Scan(props...)
		if err != nil {
			return nil, fmt.Errorf("%v: error reading customer row: %w", queryName, err)
		}
		items = append(items, *item)
	}
	return items, err
}

// selectingOne executes the query and converts the obtained row into an instance of T. It returns wrapped db.ErrNoRows if the
// query returned no rows and db.ErrTooManyRows if the query returned more than 1 row.
//
// queryName is only used in error messages.
func selectingOne[T any](ctx context.Context, queryName string, sql SqlExecutor, rowScanSupplier rowScanSupplier[T],
	query string, args ...any) (T, error) {

	rows, err := sql.QueryContext(ctx, query, args...)
	if err != nil {
		return util.ZeroValue[T](), fmt.Errorf("%v: error running query: %w", queryName, err)
	}
	defer closeRows(ctx, rows)

	if !rows.Next() {
		return util.ZeroValue[T](), fmt.Errorf("%s: %w", queryName, db.ErrNoRows)
	}
	result, props := rowScanSupplier()
	err = rows.Scan(props...)
	if err != nil {
		return util.ZeroValue[T](), err
	}
	if rows.Next() {
		return util.ZeroValue[T](), fmt.Errorf("%s: %w", queryName, db.ErrTooManyRows)
	}
	return *result, err
}

// affectingMany executes the query and returns the number of affected rows. It's not considered an error if the query affected no rows.
//
// queryName is only used in error messages.
func affectingMany(ctx context.Context, queryName string, sql SqlExecutor, query string, args ...any) (affectedRowsCount int64, err error) {
	result, err := sql.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("%v: error running query: %w", queryName, err)
	}

	affectedRowsCount, err = result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("%v: error obtaining the number of affected rows: %w", queryName, err)
	}

	return affectedRowsCount, nil
}

// affectingOne executes the query. It returns wrapped db.ErrNoRows if the query affected no rows and db.ErrTooManyRows
// if the query affected more than 1 row.
//
// queryName is only used in error messages.
func affectingOne(ctx context.Context, queryName string, sql SqlExecutor, query string, args ...any) error {
	affectedRowsCount, err := affectingMany(ctx, queryName, sql, query, args...)
	if err != nil {
		return err
	}

	if affectedRowsCount == 0 {
		return fmt.Errorf("%s: %w", queryName, db.ErrNoRows)
	}
	if affectedRowsCount > 1 {
		return fmt.Errorf("%s: %w", queryName, db.ErrTooManyRows)
	}
	return nil
}

func closeRows(ctx context.Context, rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		mdctx.Errorf(ctx, "Error closing rows: %v", err)
	}
}
