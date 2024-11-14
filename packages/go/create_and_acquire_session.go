package redis_commands

import (
	"context"

	"github.com/my-ermes-labs/api-go/api"
)

// Creates a new session and acquires it. Returns the id of the session.
func (c *RedisCommands) CreateAndAcquireSession(
	ctx context.Context,
	options api.CreateAndAcquireSessionOptions,
) (string, error) {
	return "XXX", nil
}
