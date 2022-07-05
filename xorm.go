package main

import (
	"context"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/samber/lo"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

func main() {
	db := newMySQLClient()

	ctx := context.Background()
	session := db.NewSession().Context(ctx)
	defer session.Close()
}

func newMySQLClient() *xorm.Engine {
	// generate client
	db, err := xorm.NewEngine("mysql", "user:password@tcp(localhost:3306)/rdb")
	if err != nil {
		panic(err)
	}

	// set logger
	logger := log.NewSimpleLogger(os.Stdout)
	logger.ShowSQL(true)
	db.SetLogger(logger)

	return db
}

func getLock(session *xorm.Session, key string) (bool, error) {
	sql := "select get_lock(?, 100)"
	m, err := session.Query(sql, key)
	if err != nil {
		return false, err
	}

	byteMapToStringMap := func(v map[string][]byte) map[string]string {
		return lo.MapValues(v, func(b []byte, _ string) string { return string(b) })
	}
	res := lo.Map(m, func(v map[string][]byte, _ int) map[string]string {
		return byteMapToStringMap(v)
	})

	return lo.ContainsBy(res, func(m map[string]string) bool {
		v, ok := m["get_lock"]
		return ok && v == "1"
	}), nil
}
