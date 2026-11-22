package api

import "github.com/danielgtaylor/huma/v2"

func requireAuth(op huma.Operation) huma.Operation {
	op.Security = []map[string][]string{
		{"BearerAuth": {}},
	}
	return op
}
