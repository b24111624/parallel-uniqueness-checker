package checker

import "context"

type Store interface {
	// Run starts checker process
	Run(ctx context.Context) func() error
}
