package redis_commands

import (
	"context"
	"fmt"

	"github.com/my-ermes-labs/api-go/api"
	"github.com/my-ermes-labs/log"
)

// Acquires a session. If the session has been offloaded and not acquired it
// returns the new session sessionLocation, otherwise nil. The options defines how
// the session is acquired.
// errors:
// - ErrSessionNotFound: If no session with the given id is found.
// - ErrSessionIsOffloading: If the session is offloading and the required permission is read-write.
func (c *RedisCommands) AcquireSession(ctx context.Context, sessionId string, opt api.AcquireSessionOptions) (*api.SessionLocation, error) {
	var allow_offloading string

	log.MyLog("\nSTART ACQUISITION OF " + sessionId)
	if opt.AllowOffloading() {
		allow_offloading = "1"
	} else {
		allow_offloading = "0"
	}

	log.MyLog("Allow offloading = " + allow_offloading)

	var allow_while_offloading string
	if opt.AllowWhileOffloading() {
		allow_while_offloading = "1"
	} else {
		allow_while_offloading = "0"
	}
	log.MyLog("Allow while offloading = " + allow_while_offloading)

	res, err := c.client.FCall(ctx, "acquire_session", []string{sessionId}, allow_offloading, allow_while_offloading).StringSlice()

	log.MyLog(fmt.Sprintf("Acquisition result =%v ", res))

	if err != nil {
		log.MyLog(fmt.Sprintf("Error while acquiring session = %v. Res = %x", err, res))
		return nil, err
	}

	if len(res) == 3 {
		offloaded_to_host, offloaded_to_session := res[1], res[2]
		log.MyLog("offloaded to host = " + offloaded_to_host + "offloaded to session = " + offloaded_to_session)

		location := api.NewSessionLocation(offloaded_to_host, offloaded_to_session)

		log.MyLog(fmt.Sprintf("location =%v ", location))

		return &location, nil
	}

	return nil, err
}

func (c *RedisCommands) ReleaseSession(
	ctx context.Context,
	sessionId string,
	opt api.AcquireSessionOptions,
) (*api.SessionLocation, error) {
	var allow_offloading string
	if opt.AllowOffloading() {
		allow_offloading = "1"
	} else {
		allow_offloading = "0"
	}

	res, err := c.client.FCall(ctx, "acquire_session", []string{sessionId}, allow_offloading).StringSlice()

	if err != nil {
		return nil, err
	}

	if len(res) == 3 {
		offloaded_to_host, offloaded_to_session := res[1], res[2]

		location := api.NewSessionLocation(offloaded_to_host, offloaded_to_session)

		return &location, nil
	}

	return nil, err
}

// Returns the offloadable sessions, the function returns the new cursor, the
// list of session ids and an error. The cursor is used to paginate the results.
// If the cursor is empty, the function returns the first page of results.
// errors:
// - ErrInvalidCursor: If the cursor is invalid.
// - ErrInvalidCount: If the count is invalid.
func (c *RedisCommands) ScanOffloadableSessions(
	ctx context.Context,
	cursor uint64,
	count int64,
) ([]string, uint64, error) {
	results, newCursor, err := c.client.ZScan(ctx, "c:offloadable_sessions_set", cursor, "*", count).Result()
	if err != nil {
		return nil, 0, err
	}

	keys := make([]string, 0, len(results)/2)
	for i := 0; i < len(results); i += 2 {
		keys = append(keys, results[i])
	}

	return keys, newCursor, nil
}
