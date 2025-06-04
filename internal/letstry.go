package internal

import (
	"context"

	"github.com/letstrygo/letstry/internal/application"
)

func LetsTry() {
	// Initialize the Context
	ctx := context.Background()

	// Start the application
	application.NewApplication(ctx).Start()
}
