package main

import (
	"context"
	"fmt"

	//импорт драйвера для миграций
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"

	//импорт драйвера pgx
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	repo := New()

	for i := range 5000 {
		t := tag{
			s:    fmt.Sprintf("opcua_r%d", i+1),
			lo:   10,
			lolo: 5,
			hi:   20,
			hihi: 25,
		}
		repo.SaveUser(context.Background(), t)
	}

}

type tag struct {
	s    string
	lo   int
	lolo int
	hi   int
	hihi int
}

type PostgresRepo struct {
	Db *pgxpool.Pool
}

func New() *PostgresRepo {
	dsn := "postgres://postgres:postgres@localhost:5432/edgex_db?sslmode=disable"

	dbpool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		panic(err)
	}

	return &PostgresRepo{
		Db: dbpool,
	}
}

func (r *PostgresRepo) Close() {
	r.Db.Close()
}

func (r *PostgresRepo) SaveUser(ctx context.Context, t tag) error {
	_, err := r.Db.Exec(ctx,
		`INSERT INTO core_data.threshold_table(project_id, source_name, lo, lolo, hi, hihi)
    VALUES ($1, $2, $3, $4, $5, $6)`,
		"1",
		t.s,
		t.lo,
		t.lolo,
		t.hi,
		t.hihi,
	)
	return err
}
