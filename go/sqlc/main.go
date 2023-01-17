package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ss49919201/rdb-template/config"
	"github.com/ss49919201/rdb-template/sqlc/book"
)

func exit1(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func main() {
	ctx := context.Background()

	db, err := sql.Open("mysql", config.DSN())
	if err != nil {
		exit1(err)
	}

	queries := book.New(db)
	res, err := queries.CreateBook(ctx, book.CreateBookParams{
		Title:       "title",
		PublishedAt: time.Now(),
	})
	if err != nil {
		exit1(err)
	}

	insertedID, err := res.LastInsertId()
	if err != nil {
		exit1(err)
	}

	created, err := queries.GetBook(ctx, int32(insertedID))
	if err != nil {
		exit1(err)
	}

	fmt.Printf("#%v\n", created)

	// 楽観的ロック
	res, err = queries.UpdateBook(ctx, book.UpdateBookParams{
		Title: sql.NullString{
			String: "updated title",
			Valid:  true,
		},
		ID:      created.ID,
		Version: created.Version,
	})
	if err != nil {
		exit1(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		exit1(err)
	}

	if affected == 0 {
		exit1(err)
	}

	updated, err := queries.GetBook(ctx, int32(created.ID))
	if err != nil {
		exit1(err)
	}
	fmt.Printf("#%v\n", updated)
}
