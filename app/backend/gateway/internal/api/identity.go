package api

import (
	"context"
	"net/http"
	"time"

	"cardhub/backend/gateway/internal/auth"
	identityv1 "cardhub/backend/protocol/gen/identity/v1"

	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type registerInput struct {
	Body struct {
		Email    string `json:"email" format:"email"`
		Username string `json:"username" minLength:"3" maxLength:"30"`
		Password string `json:"password" minLength:"12"`
	}
}

type registerOutput struct {
	Body struct {
		UserID string `json:"user_id"`
	}
}

type loginInput struct {
	Body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
}

type tokensBody struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type tokensOutput struct {
	Body tokensBody
}

type refreshInput struct {
	Body struct {
		RefreshToken string `json:"refresh_token"`
	}
}

type logoutInput struct {
	Body struct {
		RefreshToken string `json:"refresh_token"`
	}
}

type getMeInput struct{}

type getMeOutput struct {
	Body struct {
		UserID   string `json:"user_id"`
		Email    string `json:"email"`
		Username string `json:"username"`
	}
}

func registerIdentity(
	h huma.API,
	client identityv1.IdentityServiceClient,
	logger *zap.Logger,
) {
	huma.Register(h, huma.Operation{
		OperationID:   "register",
		Method:        http.MethodPost,
		Path:          "/v1/auth/register",
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, in *registerInput) (*registerOutput, error) {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		result, err := client.Register(ctx,
			&identityv1.RegisterRequest{
				Email:    in.Body.Email,
				Username: in.Body.Username,
				Password: in.Body.Password,
			},
		)
		if err != nil {
			return nil, publicError(err, logger)
		}

		out := &registerOutput{}
		out.Body.UserID = result.GetUserId()
		return out, nil
	})

	huma.Register(h, huma.Operation{
		OperationID: "login",
		Method:      http.MethodPost,
		Path:        "/v1/auth/login",
	}, func(ctx context.Context, in *loginInput) (*tokensOutput, error) {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		result, err := client.Login(ctx,
			&identityv1.LoginRequest{
				Email:    in.Body.Email,
				Password: in.Body.Password,
			},
		)
		if err != nil {
			return nil, publicError(err, logger)
		}

		return &tokensOutput{
			Body: tokensBody{
				AccessToken:  result.GetAccessToken(),
				RefreshToken: result.GetRefreshToken(),
			},
		}, nil
	})

	huma.Register(h, huma.Operation{
		OperationID: "refresh-token",
		Method:      http.MethodPost,
		Path:        "/v1/auth/refresh",
	}, func(ctx context.Context, in *refreshInput) (*tokensOutput, error) {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		result, err := client.RefreshToken(ctx,
			&identityv1.RefreshTokenRequest{
				RefreshToken: in.Body.RefreshToken,
			},
		)
		if err != nil {
			return nil, publicError(err, logger)
		}

		return &tokensOutput{
			Body: tokensBody{
				AccessToken:  result.GetAccessToken(),
				RefreshToken: result.GetRefreshToken(),
			},
		}, nil
	})

	huma.Register(h, huma.Operation{
		OperationID:   "logout",
		Method:        http.MethodPost,
		Path:          "/v1/auth/logout",
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, in *logoutInput) (*struct{}, error) {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_, err := client.Logout(ctx,
			&identityv1.LogoutRequest{
				RefreshToken: in.Body.RefreshToken,
			},
		)
		if err != nil {
			return nil, publicError(err, logger)
		}

		return &struct{}{}, nil
	})

	huma.Register(h, requireAuth(huma.Operation{
		OperationID: "get-me",
		Method:      http.MethodGet,
		Path:        "/v1/users/me",
		Summary:     "Get current user",
		Tags:        []string{"Users"},
	}), func(ctx context.Context, in *getMeInput) (*getMeOutput, error) {
		userID, ok := auth.UserID(ctx)
		if !ok {
			return nil, huma.Error401Unauthorized("unauthorized")
		}

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// Forward only the verified user identity.
		ctx = metadata.AppendToOutgoingContext(
			ctx, "x-user-id", userID,
		)

		result, err := client.GetMe(
			ctx, &identityv1.GetMeRequest{},
		)
		if err != nil {
			return nil, publicError(err, logger)
		}

		out := &getMeOutput{}
		out.Body.UserID = result.GetUserId()
		out.Body.Email = result.GetEmail()
		out.Body.Username = result.GetUsername()

		return out, nil
	})

}
