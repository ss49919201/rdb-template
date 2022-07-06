package main

import (
	"context"
	"fmt"
	"os"
	"time"

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

	insertWithTx(session, serializable, "")
}

type isolationLevel string

const (
	readUncommitted isolationLevel = "READ UNCOMMITTED"
	readCommit      isolationLevel = "READ COMMITTED"
	repeatableRead  isolationLevel = "REPEATABLE READ"
	serializable    isolationLevel = "SERIALIZABLE"
)

func insertWithTx(session *xorm.Session, level isolationLevel, id string) {
	sql := fmt.Sprintf("SET SESSION TRANSACTION ISOLATION LEVEL %s", level)
	_, err := session.Exec(sql)
	if err != nil {
		panic(err)
	}

	type User struct {
		ID        string `xorm:"id"`
		Name      string
		Count     int
		UpdatedAt time.Time
	}

	// トランザクション開始
	if err := session.Begin(); err != nil {
		panic(err)
	}

	// 処理をロック
	var a string
	fmt.Scan(&a)

	// 存在チェックA
	// READ UNCOMMITTED だと別トランザクションのコミット前INSERTが見える(ファントムリード)のでここで落ちる
	// SERIALIZABLE だと別トランザクションと直列化されるのでここでブロック
	if exsit, err := session.Table("users").Exist(&User{ID: id}); err != nil {
		panic(err)
	} else if exsit {
		panic("already exists A")
	}

	// 処理をロック
	fmt.Scan(&a)

	// 存在チェックB
	// READ COMMITTED だと別トランザクションのコミット後INSERTが見える(ファントムリード)のでここで落ちる
	if exsit, err := session.Table("users").Exist(&User{ID: id}); err != nil {
		panic(err)
	} else if exsit {
		panic("already exists B")
	}

	// 挿入
	// REPEATABLE READ だと別トランザクションのコミット前後INSERTが見えないのでここで落ちる
	if _, err := session.Table("users").Insert(&User{ID: id, Name: "samber", Count: 1, UpdatedAt: time.Now()}); err != nil {
		panic(err)
	}

	// 処理をロック
	fmt.Scan(&a)

	// コミット
	if err := session.Commit(); err != nil {
		panic(err)
	}
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
