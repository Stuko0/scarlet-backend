package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct{
	Pool *pgxpool.Pool
}

func NewPostgresConnection(connString string) (*Postgres, error) {
	config, err:= pgxpool.ParseConfig(connString)
	if err != nil{return nil, fmt.Errorf("failed to parse connection string: %v",err)}

	config.MaxConns=25
	config.MinConns=5
	config.MaxConnLifetime=time.Hour
	config.MaxConnIdleTime=30*time.Minute
	config.ConnConfig.RuntimeParams["search_path"]="scarlet"
	config.AfterConnect = func (ctx context.Context, conn *pgx.Conn)error{
		_, err := conn.Exec(ctx, "SET search_path TO scarlet")
		return err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err!=nil{return nil, fmt.Errorf("failed to create connection pool: %v", err)}

	ctx, cancel:=context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err:=pool.Ping(ctx); err!=nil{
		return nil, fmt.Errorf("failed to ping database: %v",err)
	}

	return &Postgres{Pool: pool}, nil
}



func (db *Postgres) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}