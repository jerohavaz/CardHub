// Package logging shares the Zap logger between services.
// Adapted from TimBRA's foundation/logging implementation.
package logging

import "go.uber.org/zap"

func New(development bool) (*zap.Logger, error) {
	if development {
		return zap.NewDevelopment()
	}
	return zap.NewProduction()
}
