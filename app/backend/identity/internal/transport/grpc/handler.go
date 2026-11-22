package grpc

import (
	"context"
	"errors"

	"cardhub/backend/identity/internal/service"
	identityv1 "cardhub/backend/protocol/gen/identity/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Handler struct {
	identityv1.UnimplementedIdentityServiceServer
	svc *service.Service
	log *zap.Logger
}

func New(svc *service.Service, log *zap.Logger) *Handler {
	return &Handler{
		svc: svc,
		log: log,
	}
}

func (h *Handler) Register(
	ctx context.Context,
	req *identityv1.RegisterRequest,
) (*identityv1.RegisterResponse, error) {
	user, err := h.svc.Register(
		ctx,
		req.GetEmail(),
		req.GetUsername(),
		req.GetPassword(),
	)
	if err != nil {
		return nil, h.mapError(err)
	}

	return &identityv1.RegisterResponse{
		UserId: user.ID,
	}, nil
}

func (h *Handler) Login(
	ctx context.Context,
	req *identityv1.LoginRequest,
) (*identityv1.LoginResponse, error) {
	tokens, err := h.svc.Login(
		ctx,
		req.GetEmail(),
		req.GetPassword(),
	)
	if err != nil {
		return nil, h.mapError(err)
	}

	return &identityv1.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (h *Handler) RefreshToken(
	ctx context.Context,
	req *identityv1.RefreshTokenRequest,
) (*identityv1.RefreshTokenResponse, error) {
	tokens, err := h.svc.RefreshToken(
		ctx,
		req.GetRefreshToken(),
	)
	if err != nil {
		return nil, h.mapError(err)
	}

	return &identityv1.RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (h *Handler) Logout(
	ctx context.Context,
	req *identityv1.LogoutRequest,
) (*identityv1.LogoutResponse, error) {
	if err := h.svc.Logout(
		ctx,
		req.GetRefreshToken(),
	); err != nil {
		return nil, h.mapError(err)
	}

	return &identityv1.LogoutResponse{}, nil
}

func (h *Handler) GetMe(
	ctx context.Context,
	req *identityv1.GetMeRequest,
) (*identityv1.GetMeResponse, error) {
	// The gateway forwards the verified user ID.
	// This metadata must come from an authenticated
	// and trusted internal gRPC caller.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(
			codes.Unauthenticated,
			"missing authentication",
		)
	}

	ids := md.Get("x-user-id")
	if len(ids) != 1 || ids[0] == "" {
		return nil, status.Error(
			codes.Unauthenticated,
			"missing user identity",
		)
	}

	user, err := h.svc.GetMe(ctx, ids[0])
	if err != nil {
		return nil, h.mapError(err)
	}

	return &identityv1.GetMeResponse{
		UserId:   user.ID,
		Email:    user.Email,
		Username: user.Username,
	}, nil
}

func (h *Handler) mapError(err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidEmail),
		errors.Is(err, service.ErrInvalidUsername),
		errors.Is(err, service.ErrInvalidPassword),
		errors.Is(err, service.ErrInvalidID):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, service.ErrUserExists):
		return status.Error(
			codes.AlreadyExists,
			"user already exists",
		)

	case errors.Is(err, service.ErrInvalidCredentials),
		errors.Is(err, service.ErrInvalidSession):
		return status.Error(
			codes.Unauthenticated,
			"invalid credentials or session",
		)

	case errors.Is(err, service.ErrAccountDisabled):
		return status.Error(
			codes.PermissionDenied,
			"account disabled",
		)

	case errors.Is(err, service.ErrNotFound):
		return status.Error(
			codes.NotFound,
			"user not found",
		)

	default:
		h.log.Error(
			"identity operation failed",
			zap.Error(err),
		)
		return status.Error(
			codes.Internal,
			"internal server error",
		)
	}
}
