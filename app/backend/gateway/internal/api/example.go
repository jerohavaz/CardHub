package api

import (
	"context"
	"net/http"
	"time"

	examplev1 "cardhub/backend/protocol/gen/example/v1"
	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/zap"
)

type exampleBody struct {
	ID   string `json:"id" doc:"Generated example UUID"`
	Name string `json:"name" doc:"Display name"`
}

type exampleOutput struct{ Body exampleBody }

type createInput struct {
	Body struct {
		Name string `json:"name" minLength:"1" maxLength:"120"`
	}
}

type getInput struct {
	ID string `path:"id" doc:"Example UUID"`
}

func registerExample(h huma.API, client examplev1.ExampleServiceClient, logger *zap.Logger) {
	huma.Register(h, huma.Operation{
		OperationID:   "create-example",
		Method:        http.MethodPost,
		Path:          "/v1/examples",
		DefaultStatus: http.StatusCreated,
		Summary:       "Create an example",
	}, func(ctx context.Context, in *createInput) (*exampleOutput, error) {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		out, err := client.CreateExample(ctx, &examplev1.CreateExampleRequest{Name: in.Body.Name})
		if err != nil {
			return nil, publicError(err, logger)
		}

		return &exampleOutput{Body: exampleBody{ID: out.GetId(), Name: out.GetName()}}, nil
	})

	huma.Register(h, huma.Operation{
		OperationID: "get-example",
		Method:      http.MethodGet,
		Path:        "/v1/examples/{id}",
		Summary:     "Get an example by UUID",
	}, func(ctx context.Context, in *getInput) (*exampleOutput, error) {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		out, err := client.GetExample(ctx, &examplev1.GetExampleRequest{Id: in.ID})
		if err != nil {
			return nil, publicError(err, logger)
		}

		return &exampleOutput{Body: exampleBody{ID: out.GetId(), Name: out.GetName()}}, nil
	})
}
