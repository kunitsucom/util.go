package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"

	"github.com/kunitsucom/util.go/env"
	errorz "github.com/kunitsucom/util.go/errors"
	syncz "github.com/kunitsucom/util.go/sync"
)

//nolint:gosec,revive,stylecheck
const (
	MYSQL_DSN           = "MYSQL_DSN"
	MYSQL_ROOT_PASSWORD = "MYSQL_ROOT_PASSWORD"
	MYSQL_USER          = "MYSQL_USER"
	MYSQL_PASSWORD      = "MYSQL_PASSWORD"
	MYSQL_DATABASE      = "MYSQL_DATABASE"
)

type mysqlLogger struct {
	enable bool
}

func (m *mysqlLogger) Print(v ...interface{}) {
	if !m.enable {
		return
	}
	log.Print("[mysql] ", fmt.Sprint(v...))
}

//nolint:gochecknoglobals
var (
	_DSN      string
	_Shutdown func(ctx context.Context) error
	once      syncz.Once
)

func NewTestDB(ctx context.Context) (dsn string, cleanup func(ctx context.Context) error, err error) {
	return newTestDB(ctx, "8.1")
}

//nolint:funlen
func newTestDB(ctx context.Context, imageTag string) (dsn string, cleanup func(ctx context.Context) error, err error) {
	dsn, err = env.String(MYSQL_DSN)
	if err == nil {
		return dsn, func(_ context.Context) error { return nil /* noop */ }, nil
	}

	if err := once.Do(func() error {
		m := &mysqlLogger{}
		if err := mysql.SetLogger(m); err != nil {
			return errorz.Errorf("mysql.SetLogger: %w", err)
		}
		defer func() { m.enable = true }()

		dockertestPool, err := dockertest.NewPool("")
		if err != nil {
			return errorz.Errorf("dockertest.NewPool: %w", err)
		}
		dockertestPool.MaxWait = 30 * time.Second

		const databaseDriver = "mysql"
		var (
			databaseRootPassword = env.StringOrDefault(MYSQL_ROOT_PASSWORD, "password")
			databaseUser         = env.StringOrDefault(MYSQL_USER, "testuser")
			databasePassword     = env.StringOrDefault(MYSQL_PASSWORD, "testpassword")
			databaseName         = env.StringOrDefault(MYSQL_DATABASE, "testdb")
		)

		dockertestRunOptions := &dockertest.RunOptions{
			Repository: "mysql",
			Tag:        imageTag,
			Env: []string{
				fmt.Sprintf("%s=%s", MYSQL_ROOT_PASSWORD, databaseRootPassword),
				fmt.Sprintf("%s=%s", MYSQL_USER, databaseUser),
				fmt.Sprintf("%s=%s", MYSQL_PASSWORD, databasePassword),
				fmt.Sprintf("%s=%s", MYSQL_DATABASE, databaseName),
			},
		}

		dockertestResource, err := dockertestPool.RunWithOptions(dockertestRunOptions,
			func(config *docker.HostConfig) {
				config.AutoRemove = true
				config.RestartPolicy = docker.RestartPolicy{
					Name: "no",
				}
			},
		)
		if err != nil {
			return errorz.Errorf("dockertestPool.RunWithOptions: %w", err)
		}

		databaseDSN := fmt.Sprintf("%s:%s@tcp(%s)/%s", databaseUser, databasePassword, dockertestResource.GetHostPort("3306/tcp"), databaseName)
		log.Printf("✅: MYSQL_DSN: %s", databaseDSN)

		if err := dockertestPool.Retry(func() error {
			db, err := sql.Open(databaseDriver, databaseDSN)
			if err != nil {
				err = errorz.Errorf("dockertestPool.Retry: sql.Open: %w", err)
				log.Print(err)
				return err
			}
			defer db.Close()

			if err := db.PingContext(ctx); err != nil {
				err = errorz.Errorf("dockertestPool.Retry: db.Ping: %w", err)
				log.Print(err)
				return err
			}

			return nil
		}); err != nil {
			return errorz.Errorf("pool.Retry: %w", err)
		}

		_DSN = databaseDSN

		_Shutdown = func(_ context.Context) error {
			if err := dockertestPool.Purge(dockertestResource); err != nil {
				return errorz.Errorf("dockertestPool.Purge: %w", err)
			}
			return nil
		}

		return nil
	}); err != nil {
		return "", nil, errorz.Errorf("_MysqlOnce.Do: %w", err)
	}

	return _DSN, _Shutdown, nil
}
