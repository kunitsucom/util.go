package sql_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/kunitsucom/ilog.go"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

func mysqlDSN(databaseUser, databasePassword, hostAndPort, databaseName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", databaseUser, databasePassword, hostAndPort, databaseName)
}

func TestMain(m *testing.M) {
	var code int = -1
	defer func() { os.Exit(code) }()

	ilog.L().Debugf("start")
	ilog.SetStdLogger(ilog.L().AddCallerSkip(1))
	if err := mysql.SetLogger(log.New(ilog.L().Copy().AddCallerSkip(2), "", 0)); err != nil {
		ilog.L().Errorf("mysql.SetLogger: %v", err)
		return
	}

	const (
		databaseDriver   = "mysql"
		dockerRepository = "mysql"
		dockerTag        = "8.1"
		databaseUser     = "testuser"
		databasePassword = "testpassword"
		databaseName     = "testdb"
		portID           = "3306/tcp"
	)

	Envs := []string{
		fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", "password"),
		fmt.Sprintf("MYSQL_USER=%s", databaseUser),
		fmt.Sprintf("MYSQL_PASSWORD=%s", databasePassword),
		fmt.Sprintf("MYSQL_DATABASE=%s", databaseName),
	}

	ilog.L().Debugf("dockertest.NewPool")
	pool, err := dockertest.NewPool("")
	if err != nil {
		ilog.L().Errorf("Could not connect to docker: %v", err)
		return
	}
	pool.MaxWait = 30 * time.Second

	// pwd, _ := os.Getwd()

	runOptions := &dockertest.RunOptions{
		Repository: dockerRepository,
		// latest だと本番とマッチしなくなる場合があるのでバージョン指定
		Tag: dockerTag,
		// ポート番号は固定せずに 0 で listen する
		Env: Envs,
		// ここでデータベースの初期化ファイルを渡す
		// 	Mounts: []string{
		// 		pwd + "/db/schema.sql:/docker-entrypoint-initdb.d/schema.sql",
		// 	},
	}

	ilog.L().Debugf("pool.RunWithOptions")
	resource, err := pool.RunWithOptions(runOptions,
		func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		},
	)
	defer func() {
		if err := pool.Purge(resource); err != nil {
			ilog.L().Errorf("Could not purge resource: %v", err)
			return
		}
	}()
	if err != nil {
		ilog.L().Errorf("Could not start resource: %v", err)
		return
	}

	hostAndPort := resource.GetHostPort(portID)
	databaseDSN := mysqlDSN(databaseUser, databasePassword, hostAndPort, databaseName)

	// docker が起動するまで少し時間がかかるのでリトライする
	ilog.L().Debugf("pool.Retry: %s", databaseDSN)
	if err := pool.Retry(func() error {
		db, err := sql.Open(databaseDriver, databaseDSN)
		if err != nil {
			ilog.L().Infof("sql.Open: %v", err)
			return err
		}
		defer db.Close()

		if err := db.Ping(); err != nil {
			ilog.L().Infof("db.Ping: %v", err)
			return err
		}

		return nil
	}); err != nil {
		ilog.L().Errorf("Could not connect to database: %s", err)
		return
	}

	code = m.Run()
}
