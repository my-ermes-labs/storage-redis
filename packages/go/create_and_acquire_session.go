package redis_commands

import (
	"context"
	"fmt"

	"github.com/my-ermes-labs/api-go/api"
	"github.com/my-ermes-labs/log"
)

// Creates a new session and acquires it. Returns the id of the session.
func (c *RedisCommands) CreateAndAcquireSession(
	ctx context.Context,
	options api.CreateAndAcquireSessionOptions,
) (string, error) {
	log.MyLog("CREATE AND ACQUIRE SESSION")

	sessionId, err := c.CreateSession(ctx, options.CreateSessionOptions)
	log.MyLog("SESSION ID = " + sessionId)
	if err != nil {
		log.MyLog(fmt.Sprintf("Error during creation = %v ", err))
		return "", err
	}

	log.MyLog("creation done. SessionId = " + sessionId)
	_, err = c.AcquireSession(ctx, sessionId, options.AcquireSessionOptions)
	if err != nil {
		log.MyLog(fmt.Sprintf("Error during acquisition = %v ", err))
		return "", err
	}

	log.MyLog("acquisition done. SessionId = " + sessionId)
	return sessionId, nil
}
