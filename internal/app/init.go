package app

import (
	"context"
	"log"

	gateway_runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/Shemistan/uzum_admin/docs"
	"github.com/Shemistan/uzum_admin/internal/api"
	pb_admin "github.com/Shemistan/uzum_admin/pkg/admin_v1"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/jmoiron/sqlx"
)

func (a *App) initDB() {
	sqlConnectionString := a.getSqlConnectionString()

	var err error
	a.db, err = sqlx.Open("postgres", sqlConnectionString)
	if err != nil {
		log.Fatal("failed to opening connection to db: ", err.Error())
	}

	if err = a.db.Ping(); err != nil {
		log.Fatal("failed to connect to the database: ", err.Error())
	}
}

func (a *App) initReDoc() {
	a.reDoc = docs.Initialize()
}

func (a *App) initGRPCServer() {
	a.grpcAdminServer = grpc.NewServer()

	pb_admin.RegisterAdminV1Server(
		a.grpcAdminServer,
		&api.Admin{
			AdminService: a.getService(),
		},
	)
}

func (a *App) initHTTPServer(ctx context.Context) error {
	a.muxAdmin = gateway_runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := pb_admin.RegisterAdminV1HandlerFromEndpoint(ctx, a.muxAdmin, a.appConfig.App.PortGRPC, opts)
	if err != nil {
		return err
	}

	return nil
}
