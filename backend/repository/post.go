package repository

import (
	"cmp"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func NewDB() (*sqlx.DB, error) {
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
	return db, nil
}

type PostRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) *PostRepository {

	return &PostRepository{db: db}
}

func (pr *PostRepository) CreatePost(ctx context.Context, postID uuid.UUID, originalMessage string, convertedMessage string, parentID uuid.UUID) error {
	db := pr.db
	var rootID uuid.UUID

	if postID == parentID { // リプライじゃない
		rootID = postID
	} else {
		rootId := &uuid.UUID{}
		err := db.Get(&rootId, "SELECT root_id FROM posts WHERE id=?", parentID)

		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("not found: %w", err)
		}
		if err != nil {
			return fmt.Errorf("failed to get %w", err)
		}
	}

	_, err := db.Exec("INSERT INTO posts (id, original_message, converted_message, parent_id, root_id)", postID, originalMessage, convertedMessage, parentID, rootID)
	if err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}
