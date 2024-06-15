package repository

import (
	"cmp"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type PostRepository struct {
	db *sqlx.DB
}

func NewPostRepository() (*PostRepository, error) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	mysqlConf := mysql.Config{
		User:                 cmp.Or(os.Getenv("DB_USER"), "root"),
		Passwd:               cmp.Or(os.Getenv("DB_PASSWORD"), "passsword"),
		Addr:                 cmp.Or(os.Getenv("DB_HOST"), "db") + ":" + cmp.Or(os.Getenv("DB_PORT"), "3306"),
		DBName:               cmp.Or(os.Getenv("DB_NAME"), "h24s_24"),
		Loc:                  jst,
		AllowNativePasswords: true,
		ParseTime:            true,
		Collation:            "utf8mb4_general_ci",
	}
	db, err := sqlx.Connect("mysql", mysqlConf.FormatDSN())
	if err != nil {
		return nil, err
	}
	return &PostRepository{db: db}, nil
}
