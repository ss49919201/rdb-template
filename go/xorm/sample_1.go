package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/s-beats/rdb-template/model"
	"github.com/samber/lo"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

type isolationLevel string

const (
	readUncommitted isolationLevel = "READ UNCOMMITTED"
	readCommit      isolationLevel = "READ COMMITTED"
	repeatableRead  isolationLevel = "REPEATABLE READ"
	serializable    isolationLevel = "SERIALIZABLE"
)

func main() {
	engine := newMySQLClient()

	ctx := context.Background()
	session := engine.NewSession().Context(ctx)
	defer session.Close()

	var u model.User
	session.Begin()
	has, err := session.Where("id = ?", 1).Get(&u)
	// has, err := session.Where("id = ?", 1).ForUpdate().Get(&u)
	if err != nil || !has {
		fmt.Println(err)
		panic("failed to get")
	}
	if u.Count > 0 {
		if _, err := session.
			Where("id = ?", "1").Decr("count").Update(&model.User{}); err != nil {
			fmt.Println(err)
			panic("failed to update")
		}
	}
	time.Sleep(30 * time.Second)
	session.Commit()
}

func badExample(session *xorm.Session) {
	var u model.User
	session.Begin()
	has, err := session.Where("id = ?", 1).Get(&u)
	if err != nil || !has {
		fmt.Println(err)
		panic("failed to get")
	}
	// Countが1の場合に複数セッションから同時に更新されると-1になる可能性がある
	if u.Count > 0 {
		if _, err := session.
			Where("id = ?", "1").Decr("count").Update(&model.User{}); err != nil {
			fmt.Println(err)
			panic("failed to update")
		}
	}
	session.Commit()
}

func goodExample(session *xorm.Session) {
	var u model.User
	session.Begin()
	has, err := session.Where("id = ?", 1).ForUpdate().Get(&u)
	if err != nil || !has {
		fmt.Println(err)
		panic("failed to get")
	}
	// コミット又はロールバックされるまでロックされるので、Countが1の場合に複数セッションから同時に更新されることはない
	if u.Count > 0 {
		if _, err := session.
			Where("id = ?", "1").Decr("count").Update(&model.User{}); err != nil {
			fmt.Println(err)
			panic("failed to update")
		}
	}
	session.Commit()
}

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
	if exsit, err := session.Table("users").Exist(&model.User{ID: id}); err != nil {
		panic(err)
	} else if exsit {
		panic("already exists A")
	}

	// 処理をロック
	fmt.Scan(&a)

	// 存在チェックB
	// READ COMMITTED だと別トランザクションのコミット後INSERTが見える(ファントムリード)のでここで落ちる
	if exsit, err := session.Table("users").Exist(&model.User{ID: id}); err != nil {
		panic(err)
	} else if exsit {
		panic("already exists B")
	}

	// 挿入
	// REPEATABLE READ だと別トランザクションのコミット前後INSERTが見えないのでここで落ちる
	if _, err := session.Table("users").Insert(&model.User{ID: id, Name: "samber", Count: 1, UpdatedAt: time.Now()}); err != nil {
		panic(err)
	}

	// 処理をロック
	fmt.Scan(&a)

	// コミット
	if err := session.Commit(); err != nil {
		panic(err)
	}
}

func serializableUpdate(engine *xorm.Engine, key string) {
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		fmt.Println(decrCount(engine.NewSession().Context(ctx), key))
	}
}

func concurrencyUpdate(engine *xorm.Engine, key string) {
	ctx := context.Background()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			fmt.Println(decrCount(engine.NewSession().Context(ctx), key))
			wg.Done()
		}()
	}
	wg.Wait()
}

func decrCount(session *xorm.Session, key string) error {
	if err := session.Begin(); err != nil {
		return err
	}

	if _, err := session.Table("users").
		Where("id = ?", key).
		Exist(); err != nil {
		return err
	}

	if _, err := session.Table("users").
		Where("id = ?", key).
		Where("name = ?", "samber").
		Decr("count").
		Update(struct{}{}); err != nil {
		return err
	}

	return session.Commit()
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
