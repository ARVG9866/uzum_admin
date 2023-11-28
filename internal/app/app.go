package app

import (
	"context"
	"fmt"
	"log"

	gateway_runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/mvrilo/go-redoc"
	"google.golang.org/grpc"

	"github.com/Shemistan/uzum_admin/dev"
	"github.com/Shemistan/uzum_admin/internal/models"
	"github.com/Shemistan/uzum_admin/internal/service/admin_v1"
	"github.com/Shemistan/uzum_admin/internal/storage"
)

type App struct {
	appConfig *models.Config
	muxAdmin  *gateway_runtime.ServeMux

	grpcAdminServer *grpc.Server
	adminService    admin_v1.IService
	db              *sqlx.DB
	reDoc           redoc.Redoc
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}
	a.setConfig()
	a.initDB()
	a.initReDoc()
	a.initGRPCServer()

	if err := a.initHTTPServer(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) setConfig() {
	if dev.DEBUG {
		err := dev.SetConfig()
		if err != nil {
			log.Fatal("failed to get config", err.Error())
		}
	}

	conf := models.Config{}

	envconfig.MustProcess("", &conf)

	a.appConfig = &conf
}

func (a *App) getSqlConnectionString() string {
	sqlConnectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%v",
		a.appConfig.DB.User,
		a.appConfig.DB.Password,
		a.appConfig.DB.Host,
		a.appConfig.DB.Port,
		a.appConfig.DB.Database,
		a.appConfig.DB.SSLMode,
	)

	return sqlConnectionString
}

func (a *App) getService() admin_v1.IService {
	storage := storage.NewStorage(a.db)

	if a.adminService == nil {
		a.adminService = admin_v1.NewService(storage)
	}

	return a.adminService
}
