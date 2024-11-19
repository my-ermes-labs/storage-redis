package redis_commands

import (
	"context"
	"fmt"

	"github.com/my-ermes-labs/api-go/api"
)

// Creates a new session and acquires it. Returns the id of the session.
func (c *RedisCommands) CreateAndAcquireSession(
	ctx context.Context,
	options api.CreateAndAcquireSessionOptions,
) (string, error) {
	log("CREATE AND ACQUIRE SESSION")

	sessionId, err := c.CreateSession(ctx, options.CreateSessionOptions)
	log("SESSION ID = " + sessionId)
	if err != nil {
		log(fmt.Sprintf("Error during creation = %v ", err))
		return "", err
	}

	log("creation done. SessionId = " + sessionId)
	_, err = c.AcquireSession(ctx, sessionId, options.AcquireSessionOptions)
	if err != nil {
		log(fmt.Sprintf("Error during acquisition = %v ", err))
		return "", err
	}

	log("acquisition done. SessionId = " + sessionId)
	return sessionId, nil
}
