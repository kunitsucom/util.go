package grpcgateway

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RegisterServices(ctx context.Context, grpcServer grpc.ServiceRegistrar, mux *runtime.ServeMux, conn *grpc.ClientConn, registrars ...*GRPCServiceRegistrar) error {
	for _, r := range registrars {
		grpcServer.RegisterService(r.grpcServiceDesc, r.grpcServer)
		if err := r.grpcGatewayHandler(ctx, mux, conn); err != nil {
			return fmt.Errorf("r.grpcGatewayHandler: %s: %w", r.grpcServiceDesc.ServiceName, err)
		}
	}

	return nil
}

type GRPCServiceRegistrar struct {
	grpcServiceDesc    *grpc.ServiceDesc
	grpcServer         interface{}
	grpcGatewayHandler func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
}

func NewGRPCServiceRegistrar[GRPCServiceServer interface{}](
	grpcServiceDesc *grpc.ServiceDesc,
	grpcServer GRPCServiceServer,
	grpcGatewayHandler func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error,
) *GRPCServiceRegistrar {
	return &GRPCServiceRegistrar{
		grpcServiceDesc:    grpcServiceDesc,
		grpcServer:         grpcServer,
		grpcGatewayHandler: grpcGatewayHandler,
	}
}
