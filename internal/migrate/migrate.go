package migrate

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"path"
	"regexp"
	"sort"
	"strconv"
)

//go:embed sql/*.up.sql
var migrationFS embed.FS

var fileNameRe = regexp.MustCompile(`^(\d+)_[^.]+\.up\.sql$`)

const ensureSchemaMigrations = `
CREATE TABLE IF NOT EXISTS schema_migrations (
	version INT PRIMARY KEY,
	applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
`

// Up применяет все ещё не применённые .up.sql миграции по возрастанию version.
func Up(ctx context.Context, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("migrate: db is nil")
	}

	if _, err := db.ExecContext(ctx, ensureSchemaMigrations); err != nil {
		return fmt.Errorf("migrate: ensure schema_migrations: %w", err)
	}

	applied, err := loadApplied(ctx, db)
	if err != nil {
		return err
	}

	files, err := listUpMigrations()
	if err != nil {
		return err
	}

	for _, m := range files {
		if applied[m.version] {
			continue
		}

		body, err := fs.ReadFile(migrationFS, m.path)
		if err != nil {
			return fmt.Errorf("migrate: read %s: %w", m.path, err)
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("migrate: begin tx v%d: %w", m.version, err)
		}

		if _, err := tx.ExecContext(ctx, string(body)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("migrate: apply v%d (%s): %w", m.version, m.path, err)
		}

		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, m.version); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("migrate: record v%d: %w", m.version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("migrate: commit v%d: %w", m.version, err)
		}
	}

	return nil
}

type migrationFile struct {
	version int
	path    string
}

func listUpMigrations() ([]migrationFile, error) {
	entries, err := migrationFS.ReadDir("sql")
	if err != nil {
		return nil, fmt.Errorf("migrate: read sql dir: %w", err)
	}

	var out []migrationFile
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		m := fileNameRe.FindStringSubmatch(name)
		if m == nil {
			continue
		}
		v, err := strconv.Atoi(m[1])
		if err != nil {
			return nil, fmt.Errorf("migrate: bad version in %q", name)
		}
		out = append(out, migrationFile{version: v, path: path.Join("sql", name)})
	}

	sort.Slice(out, func(i, j int) bool { return out[i].version < out[j].version })

	seen := make(map[int]struct{})
	for _, f := range out {
		if _, dup := seen[f.version]; dup {
			return nil, fmt.Errorf("migrate: duplicate migration version %d", f.version)
		}
		seen[f.version] = struct{}{}
	}

	return out, nil
}

func loadApplied(ctx context.Context, db *sql.DB) (map[int]bool, error) {
	rows, err := db.QueryContext(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, fmt.Errorf("migrate: list applied: %w", err)
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return nil, fmt.Errorf("migrate: scan applied: %w", err)
		}
		applied[v] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("migrate: rows applied: %w", err)
	}
	return applied, nil
}
