package db

import (
	"database/sql"
	"time"

	"github.com/astaxie/beego/logs"

	_ "github.com/go-sql-driver/mysql"
)

var MySQL = newMySQL()

// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
func newMySQL() *sql.DB {
	// todo: 隐藏密码
	dsn := "root:mima@tcp(127.0.0.1:49160)/juejin?loc=Local&parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for {
			<-ticker.C

			err := db.Ping()
			if err != nil {
				logs.Error("db ping err: %s", err)
				logs.GetBeeLogger().Flush()
				panic(err)
			}
		}
	}()
	return db
}
