package handler

import (
	"context"
	"fmt"

	"github.com/SeiFlow-3P2/api_gateway/pkg/config"
	"github.com/SeiFlow-3P2/api_gateway/pkg/env"
	authProto "github.com/SeiFlow-3P2/auth_service/pkg/proto/v1"
	boardProto "github.com/SeiFlow-3P2/board_service/pkg/proto/v1"
	calendarProto "github.com/SeiFlow-3P2/calendar_service/pkg/proto/v1"
	paymentProto "github.com/SeiFlow-3P2/payment_service/pkg/proto/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func SetupHandlers(
	ctx context.Context,
	conf *config.Config,
	gwmux *runtime.ServeMux,
	dialOpts []grpc.DialOption,
) error {
	if err := authProto.RegisterAuthServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		env.GetAuthServiceAddr(),
		dialOpts,
	); err != nil {
		return fmt.Errorf("failed to register auth service: %w", err)
	}

	if err := boardProto.RegisterBoardServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		env.GetBoardServiceAddr(),
		dialOpts,
	); err != nil {
		return fmt.Errorf("failed to register board service: %w", err)
	}

	if err := paymentProto.RegisterPaymentServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		env.GetPaymentServiceAddr(),
		dialOpts,
	); err != nil {
		return fmt.Errorf("failed to register payment service: %w", err)
	}

	if err := calendarProto.RegisterCalendarServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		env.GetCalendarServiceAddr(),
		dialOpts,
	); err != nil {
		return fmt.Errorf("failed to register calendar service: %w", err)
	}

	return nil
}
