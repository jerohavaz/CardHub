// Package service owns application rules without any gRPC or SQL imports.
package service

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

var ErrNotFound = errors.New("example not found")
var ErrInvalidName = errors.New("name must be 1-120 characters")
var ErrInvalidID = errors.New("invalid example UUID")
var uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

type Example struct {
	ID   string
	Name string
}

type Repository interface {
	Create(context.Context, string) (Example, error)
	Get(context.Context, string) (Example, error)
}

type Service struct{ repo Repository }

func New(repo Repository) *Service { return &Service{repo: repo} }
func (s *Service) Create(ctx context.Context, name string) (Example, error) {
	name = strings.TrimSpace(name)
	if utf8.RuneCountInString(name) == 0 || utf8.RuneCountInString(name) > 120 {
		return Example{}, ErrInvalidName
	}
	return s.repo.Create(ctx, name)
}
func (s *Service) Get(ctx context.Context, id string) (Example, error) {
	if !uuidPattern.MatchString(id) {
		return Example{}, ErrInvalidID
	}
	return s.repo.Get(ctx, id)
}
