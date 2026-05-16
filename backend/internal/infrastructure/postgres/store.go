package postgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"autopark/internal/domain"
	sqlfiles "autopark/sql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool

	listSQL   string
	getSQL    string
	insertSQL string
	updateSQL string
	deleteSQL string
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		pool:      pool,
		listSQL:   sqlfiles.MustRead("queries/list.sql"),
		getSQL:    sqlfiles.MustRead("queries/get.sql"),
		insertSQL: sqlfiles.MustRead("queries/insert.sql"),
		updateSQL: sqlfiles.MustRead("queries/update.sql"),
		deleteSQL: sqlfiles.MustRead("queries/delete.sql"),
	}
}

func (s *Store) List(ctx context.Context, e domain.Entity) ([]map[string]any, error) {
	q := fmt.Sprintf(s.listSQL, strings.Join(e.Columns(), ", "), e.Table, e.Order)
	return s.query(ctx, q)
}

func (s *Store) Get(ctx context.Context, e domain.Entity, id int64) (map[string]any, error) {
	q := fmt.Sprintf(s.getSQL, strings.Join(e.Columns(), ", "), e.Table, e.ID)
	return s.queryOne(ctx, q, id)
}

func (s *Store) Create(ctx context.Context, e domain.Entity, values map[string]any) (map[string]any, error) {
	cols, args, err := writeValues(e, values, true)
	if err != nil {
		return nil, err
	}
	if len(cols) == 0 {
		return nil, fmt.Errorf("empty payload")
	}
	q := fmt.Sprintf(s.insertSQL, e.Table, strings.Join(cols, ", "), placeholders(1, len(cols)), strings.Join(e.Columns(), ", "))
	return s.queryOne(ctx, q, args...)
}

func (s *Store) Update(ctx context.Context, e domain.Entity, id int64, values map[string]any) (map[string]any, error) {
	cols, args, err := writeValues(e, values, false)
	if err != nil {
		return nil, err
	}
	if len(cols) == 0 {
		return nil, fmt.Errorf("empty payload")
	}
	sets := make([]string, 0, len(cols))
	for idx, col := range cols {
		sets = append(sets, fmt.Sprintf("%s = $%d", col, idx+1))
	}
	args = append(args, id)
	q := fmt.Sprintf(s.updateSQL, e.Table, strings.Join(sets, ", "), e.ID, len(args), strings.Join(e.Columns(), ", "))
	return s.queryOne(ctx, q, args...)
}

func (s *Store) Delete(ctx context.Context, e domain.Entity, id int64) error {
	q := fmt.Sprintf(s.deleteSQL, e.Table, e.ID)
	_, err := s.pool.Exec(ctx, q, id)
	return err
}

func writeValues(e domain.Entity, values map[string]any, insert bool) ([]string, []any, error) {
	var cols []string
	var args []any
	for _, f := range e.Fields {
		if !insert && f.Name == e.ID {
			continue
		}
		raw, ok := values[f.Name]
		if !ok {
			continue
		}
		v, err := coerce(raw, f)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %w", f.Name, err)
		}
		if insert && v == nil {
			continue
		}
		cols = append(cols, f.Name)
		args = append(args, v)
	}
	return cols, args, nil
}

func coerce(v any, f domain.Field) (any, error) {
	if v == nil {
		return nil, nil
	}
	if s, ok := v.(string); ok {
		s = strings.TrimSpace(s)
		if s == "" && f.Nullable {
			return nil, nil
		}
		switch f.Kind {
		case domain.Text:
			return s, nil
		case domain.Int:
			return strconv.ParseInt(s, 10, 64)
		case domain.Numeric:
			return strconv.ParseFloat(s, 64)
		case domain.Date:
			return time.Parse("2006-01-02", s)
		}
	}

	switch f.Kind {
	case domain.Text:
		return fmt.Sprint(v), nil
	case domain.Int:
		switch x := v.(type) {
		case float64:
			return int64(x), nil
		case int64:
			return x, nil
		case int:
			return int64(x), nil
		default:
			return nil, fmt.Errorf("expected integer")
		}
	case domain.Numeric:
		switch x := v.(type) {
		case float64:
			return x, nil
		case int:
			return float64(x), nil
		case int64:
			return float64(x), nil
		default:
			return nil, fmt.Errorf("expected number")
		}
	case domain.Date:
		return nil, fmt.Errorf("expected yyyy-mm-dd")
	default:
		return v, nil
	}
}

func placeholders(start, count int) string {
	out := make([]string, count)
	for i := range out {
		out[i] = fmt.Sprintf("$%d", start+i)
	}
	return strings.Join(out, ", ")
}

func (s *Store) queryOne(ctx context.Context, q string, args ...any) (map[string]any, error) {
	rows, err := s.query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, pgx.ErrNoRows
	}
	return rows[0], nil
}

func (s *Store) query(ctx context.Context, q string, args ...any) ([]map[string]any, error) {
	rows, err := s.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fields := rows.FieldDescriptions()
	result := make([]map[string]any, 0)
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		item := make(map[string]any, len(fields))
		for i, fd := range fields {
			item[string(fd.Name)] = normalize(values[i])
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func normalize(v any) any {
	switch x := v.(type) {
	case nil:
		return nil
	case time.Time:
		if x.Hour() == 0 && x.Minute() == 0 && x.Second() == 0 && x.Nanosecond() == 0 {
			return x.Format("2006-01-02")
		}
		return x.Format(time.RFC3339)
	case []byte:
		return string(x)
	case pgtype.Numeric:
		if !x.Valid {
			return nil
		}
		f, err := x.Float64Value()
		if err == nil && f.Valid {
			return f.Float64
		}
		return x.Int.String()
	case pgtype.Date:
		if !x.Valid {
			return nil
		}
		return x.Time.Format("2006-01-02")
	case pgtype.Timestamp:
		if !x.Valid {
			return nil
		}
		return x.Time.Format(time.RFC3339)
	default:
		return v
	}
}
