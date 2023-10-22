package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func QueryWithPrepare(db *sqlx.DB) []map[string]any {
	stmt, err := db.Preparex("SELECT id FROM tbl WHERE id = 1")
	if err != nil {
		panic(err)
	}
	rows, err := stmt.Queryx()
	if err != nil {
		panic(err)
	}

	result := make([]map[string]any, 0)
	for rows.Next() {
		val := map[string]any{}
		err := rows.MapScan(val)
		if err != nil {
			panic(err)
		}
		result = append(result, val)
	}
	return result
}

func Query(db *sqlx.DB) []map[string]any {
	rows, err := db.Queryx("SELECT id FROM tbl WHERE id = 1")
	if err != nil {
		panic(err)
	}

	result := make([]map[string]any, 0)
	for rows.Next() {
		val := map[string]any{}
		err := rows.MapScan(val)
		if err != nil {
			panic(err)
		}
		result = append(result, val)
	}
	return result
}

func NamedQuery(db *sqlx.DB) []map[string]any {
	rows, err := db.NamedQuery("SELECT id FROM tbl WHERE id = :v", map[string]any{"v": 1})
	if err != nil {
		panic(err)
	}

	result := make([]map[string]any, 0)
	for rows.Next() {
		val := map[string]any{}
		err := rows.MapScan(val)
		if err != nil {
			panic(err)
		}
		result = append(result, val)
	}
	return result
}
