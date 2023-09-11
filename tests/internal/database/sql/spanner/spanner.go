package spanner

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	databaseadmin "cloud.google.com/go/spanner/admin/database/apiv1"
	databaseadminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	instanceadmin "cloud.google.com/go/spanner/admin/instance/apiv1"
	instanceadminpb "cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
	spannerdriver "github.com/googleapis/go-sql-spanner"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kunitsucom/util.go/env"
	errorz "github.com/kunitsucom/util.go/errors"
	syncz "github.com/kunitsucom/util.go/sync"
)

var _ spannerdriver.Driver

//nolint:gosec,revive,stylecheck
const (
	SPANNER_EMULATOR_HOST = "SPANNER_EMULATOR_HOST"
	GOOGLE_CLOUD_PROJECT  = "GOOGLE_CLOUD_PROJECT"
	SPANNER_INSTANCE_ID   = "SPANNER_INSTANCE_ID"
	SPANNER_DATABASE_ID   = "SPANNER_DATABASE_ID"
)

//nolint:gochecknoglobals
var (
	_DSN      string
	_Shutdown func(ctx context.Context) error
	once      syncz.Once
)

func NewTestDB(ctx context.Context) (spannerEmulatorHost string, cleanup func(ctx context.Context) error, err error) {
	return newTestDB(ctx, "latest")
}

//nolint:funlen,gocognit,cyclop
func newTestDB(ctx context.Context, imageTag string, extraStatements ...string) (spannerEmulatorHost string, cleanup func(ctx context.Context) error, err error) {
	if v := os.Getenv(SPANNER_EMULATOR_HOST); v != "" {
		_, _, _, dsn := newDSN()
		return dsn, func(_ context.Context) error { return nil /* noop */ }, nil
	}

	if err := once.Do(func() error {
		dockertestPool, err := dockertest.NewPool("")
		if err != nil {
			return errorz.Errorf("dockertest.NewPool: %w", err)
		}
		dockertestPool.MaxWait = 30 * time.Second

		const databaseDriver = "spanner"
		googleCloudProject, spannerInstanceID, spannerDatabaseID, dsn := newDSN()

		dockertestRunOptions := &dockertest.RunOptions{
			Repository: "gcr.io/cloud-spanner-emulator/emulator",
			Tag:        imageTag,
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

		log.Printf("✅: spannerDSN: %s", dsn)
		spannerEmulatorHost := dockertestResource.GetHostPort("9010/tcp")
		os.Setenv(SPANNER_EMULATOR_HOST, spannerEmulatorHost)
		log.Printf("✅: %s: %s", SPANNER_EMULATOR_HOST, spannerEmulatorHost)

		if err := dockertestPool.Retry(func() error {
			if v := os.Getenv(SPANNER_EMULATOR_HOST); v == "" {
				panic(fmt.Errorf("%s is empty", SPANNER_EMULATOR_HOST)) //nolint:goerr113
			}

			// NOTE: https://gihyo.jp/article/2023/06/tukinami-go-08
			instanceAdminClient, err := instanceadmin.NewInstanceAdminClient(ctx)
			if err != nil {
				err = errorz.Errorf("dockertestPool.Retry: instanceadmin.NewInstanceAdminClient: %w", err)
				log.Print(err)
				return err
			}
			defer instanceAdminClient.Close()
			createInstanceOp, err := instanceAdminClient.CreateInstance(ctx, &instanceadminpb.CreateInstanceRequest{
				Parent:     fmt.Sprintf("projects/%s", googleCloudProject),
				InstanceId: spannerInstanceID,
			})
			if err != nil && status.Code(err) != codes.AlreadyExists {
				err = errorz.Errorf("dockertestPool.Retry: instanceAdminClient.CreateInstance: %w", err)
				log.Print(err)
				return err
			}
			if createInstanceOp != nil {
				if _, err := createInstanceOp.Wait(ctx); err != nil {
					err = errorz.Errorf("dockertestPool.Retry: createInstanceOp.Wait: %w", err)
					log.Print(err)
					return err
				}
			}

			databaseAdminClient, err := databaseadmin.NewDatabaseAdminClient(ctx)
			if err != nil {
				err = errorz.Errorf("dockertestPool.Retry: databaseadmin.NewDatabaseAdminClient: %w", err)
				log.Print(err)
				return err
			}
			defer databaseAdminClient.Close()
			createDatabaseOp, err := databaseAdminClient.CreateDatabase(ctx, &databaseadminpb.CreateDatabaseRequest{
				Parent:          fmt.Sprintf("projects/%s/instances/%s", googleCloudProject, spannerInstanceID),
				CreateStatement: fmt.Sprintf("CREATE DATABASE `%s`", spannerDatabaseID),
				ExtraStatements: extraStatements,
			})
			if err != nil && status.Code(err) != codes.AlreadyExists {
				err = errorz.Errorf("dockertestPool.Retry: databaseAdminClient.CreateDatabase: %w", err)
				log.Print(err)
				return err
			}
			if createDatabaseOp != nil {
				if _, err := createDatabaseOp.Wait(ctx); err != nil && status.Code(err) != codes.AlreadyExists {
					err = errorz.Errorf("dockertestPool.Retry: createDatabaseOp.Wait: %w", err)
					log.Print(err)
					return err
				}
			}

			db, err := sql.Open(databaseDriver, dsn)
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

		_DSN = dsn

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

func newDSN() (googleCloudProject, spannerInstanceID, spannerDatabaseID, dsn string) {
	googleCloudProject = env.StringOrDefault(GOOGLE_CLOUD_PROJECT, "test-project")
	spannerInstanceID = env.StringOrDefault(SPANNER_INSTANCE_ID, "test-instance")
	spannerDatabaseID = env.StringOrDefault(SPANNER_DATABASE_ID, "testdb")
	dsn = fmt.Sprintf("projects/%s/instances/%s/databases/%s", googleCloudProject, spannerInstanceID, spannerDatabaseID)
	return
}
