package repository

import (
	"cmp"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewDB() (*sqlx.DB, error) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	mysqlConf := mysql.Config{
		Net:                  "tcp",
		User:                 cmp.Or(os.Getenv("DB_USER"), os.Getenv("NS_MARIADB_USER"), "root"),
		Passwd:               cmp.Or(os.Getenv("DB_PASSWORD"), os.Getenv("NS_MARIADB_PASSWORD"), "passsword"),
		Addr:                 cmp.Or(os.Getenv("DB_HOST"), os.Getenv("NS_MARIADB_HOSTNAME"), "db") + ":" + cmp.Or(os.Getenv("DB_PORT"), os.Getenv("NS_MARIADB_PORT"), "3306"),
		DBName:               cmp.Or(os.Getenv("DB_NAME"), os.Getenv("NS_MARIADB_DATABASE"), "h24s_24"),
		Loc:                  jst,
		AllowNativePasswords: true,
		ParseTime:            true,
		Collation:            "utf8mb4_general_ci",
	}

	db, err := sqlx.Connect("mysql", mysqlConf.FormatDSN())
	if err != nil {
		return nil, err
	}
	return db, nil
}
