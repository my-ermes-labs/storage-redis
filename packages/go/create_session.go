package redis_commands

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/my-ermes-labs/api-go/api"
)

// Creates a new session and returns the id of the session.
func (c *RedisCommands) CreateSession(
	ctx context.Context,
	opt api.CreateSessionOptions,
) (string, error) {
	clientGeoCoordinates := opt.ClientGeoCoordinates()
	var latitude, longitude = "", ""
	if clientGeoCoordinates != nil {
		latitude = strconv.FormatFloat(clientGeoCoordinates.Latitude, 'f', 6, 64)
		longitude = strconv.FormatFloat(clientGeoCoordinates.Longitude, 'f', 6, 64)
	}

	expiresAt := ""
	if opt.ExpiresAt() != nil {
		expiresAt = strconv.FormatInt(*opt.ExpiresAt(), 10)
	}

	acquire := ""

	log("opt.SessionID =  " + *opt.SessionId())
	for {
		var sessionId string
		if opt.SessionId() == nil {
			sessionId = uuid.NewString()
		} else {
			sessionId = *opt.SessionId()
		}
		log("SessionID =  " + sessionId)

		res, err := c.client.FCall(ctx, "create_session", []string{sessionId},
			latitude,
			longitude,
			expiresAt,
			acquire).Bool()

		log("result from redis = " + strconv.FormatBool(res))

		if err != nil {
			return "nil", err
		}

		if res {
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

func log(bodyContent string) (string, error) {
	url := "http://192.168.64.1:3000/rediscreatesession"

	requestBody := bytes.NewBufferString(bodyContent)

	req, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		return "", fmt.Errorf("error while creating the request: %v", err)
	}

	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error while sending the request: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error while reading the response: %v", err)
	}

	return string(responseBody), nil
}
