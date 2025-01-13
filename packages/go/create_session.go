package redis_commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/my-ermes-labs/api-go/api"
	"github.com/my-ermes-labs/log"
)

// Creates a new session and returns the id of the session.
func (c *RedisCommands) CreateSession(
	ctx context.Context,
	opt api.CreateSessionOptions,
) (string, error) {
	log.MyLog("CREATE SESSION")
	clientGeoCoordinates := opt.ClientGeoCoordinates()
	log.MyLog("coordinates = " + clientGeoCoordinates.String())
	var latitude, longitude = "", ""
	if clientGeoCoordinates != nil {
		log.MyLog("client coordinates")
		latitude = strconv.FormatFloat(clientGeoCoordinates.Latitude, 'f', 6, 64)
		longitude = strconv.FormatFloat(clientGeoCoordinates.Longitude, 'f', 6, 64)
	}

	expiresAt := ""
	if opt.ExpiresAt() != nil {
		log.MyLog(fmt.Sprintf("Expires AT = %v ", *opt.ExpiresAt()))
		expiresAt = strconv.FormatInt(*opt.ExpiresAt()+100, 10)
	}

	log.MyLog("expiresAt = " + expiresAt)

	acquire := ""

	for {
		var sessionId string
		if opt.SessionId() == nil {
			sessionId = uuid.NewString()
		} else {
			sessionId = *opt.SessionId()
		}
		log.MyLog("REDIS CREATE SessionID =  " + sessionId)
		log.MyLog("REDIS CREATE latitude =  " + latitude)
		log.MyLog("REDIS CREATE longitude =  " + longitude)
		log.MyLog("REDIS CREATE expired at =  " + expiresAt)
		log.MyLog("REDIS CREATE acquire =  " + acquire)

		res, err := c.client.FCall(ctx, "create_session", []string{sessionId},
			latitude,
			longitude,
			expiresAt,
			acquire).Bool()

		log.MyLog("result from redis = " + strconv.FormatBool(res))

		if err != nil {
			log.MyLog(fmt.Sprintf("err from redis create session call = %v ", err))
			return "nil", err
		}

		if res {
			log.MyLog("res = TRUE ; SessionId = " + sessionId)
			return sessionId, nil
		} else if opt.SessionId() != nil {
			return "", api.ErrSessionIdAlreadyExists
		}
	}
}

// Returns the ids of the sessions.
func (c *RedisCommands) ScanSessions(
	ctx context.Context,
	cursor uint64,
	count int64,
) ([]string, uint64, error) {
	results, newCursor, err := c.client.ZScan(ctx, "c:sessions_set", cursor, "*", count).Result()
	if err != nil {
		return nil, 0, err
	}

	keys := make([]string, 0, len(results)/2)
	for i := 0; i < len(results); i += 2 {
		keys = append(keys, results[i])
	}

	return keys, newCursor, nil
}
