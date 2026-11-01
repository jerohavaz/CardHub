package service

import (
	"context"
	"errors"
	"strings"
	"testing"
)

type fakeRepo struct {
	createName string
	found      Example
	err        error
}

func (r *fakeRepo) Create(_ context.Context, name string) (Example, error) {
	r.createName = name
	return Example{ID: "id", Name: name}, r.err
}
func (r *fakeRepo) Get(_ context.Context, id string) (Example, error) { return r.found, r.err }
func TestCreateTrimsAndValidates(t *testing.T) {
	r := &fakeRepo{}
	s := New(r)
	got, err := s.Create(context.Background(), "  Pikachu  ")
	if err != nil || got.Name != "Pikachu" || r.createName != "Pikachu" {
		t.Fatalf("got=%+v err=%v", got, err)
	}
	if _, err := s.Create(context.Background(), "  "); !errors.Is(err, ErrInvalidName) {
		t.Fatalf("error=%v", err)
	}
}
func TestUnicodeLength(t *testing.T) {
	s := New(&fakeRepo{})
	if _, err := s.Create(context.Background(), strings.Repeat("é", 120)); err != nil {
		t.Fatalf("120 unicode runes should be accepted: %v", err)
	}
	if _, err := s.Create(context.Background(), strings.Repeat("é", 121)); !errors.Is(err, ErrInvalidName) {
		t.Fatalf("expected name error: %v", err)
	}
}
func TestGetValidatesID(t *testing.T) {
	s := New(&fakeRepo{found: Example{ID: "123"}})
	if _, err := s.Get(context.Background(), "not-a-uuid"); !errors.Is(err, ErrInvalidID) {
		t.Fatalf("error=%v", err)
	}
	if _, err := s.Get(context.Background(), "123e4567-e89b-12d3-a456-426614174000"); err != nil {
		t.Fatal(err)
	}
}
