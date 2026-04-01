package db

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	Dbtx  *pgxpool.Pool
	Query *Queries
	// Q     *QueryService
)

type Config struct {
	Ctx      context.Context
	Host     string
	Port     string
	User     string
	Name     string
	Password string
}

type QueryService struct {
	*Queries
	// ISentry
}

func NewDb() {

}

func InitDb(ctx context.Context) error {

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	database := os.Getenv("POSTGRES_DATABASE")

	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, password, database)

	fmt.Printf("psqlInfo = %s\n", psqlInfo)

	Dbtx, err = pgxpool.New(ctx, psqlInfo)
	if err != nil {
		panic(err)
	}

	Query = New(Dbtx)

	if Query == nil {
		return fmt.Errorf("error: nil db.Query !")
	}

	fmt.Printf("database initialized successfully\n")
	return nil
}
