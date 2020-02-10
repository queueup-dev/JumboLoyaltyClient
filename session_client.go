package JumboLoyaltyClient

import (
	"strconv"
	"time"
)

type sessionClient struct {
	sessionTable string
}

func (s sessionClient) GetSession(sessionId string) (*Session, error) {
	result, err := Dynamo.getItem(s.sessionTable, "session_id", sessionId, &Session{})

	if err != nil {
		return nil, err
	}

	session := result.(*Session)
	return session, nil
}

func (s sessionClient) DeleteSession(sessionId string) error {
	return Dynamo.deleteItem(s.sessionTable, "session_id", sessionId)
}

func (s sessionClient) ListSessionsOfExternalId(externalId string) (*[]Session, error) {

	sessions := &[]Session{}

	result, err := Dynamo.listItems(s.sessionTable, "card_number-index", []QueryCondition{
		{
			Key:       "card_number",
			Value:     externalId,
			Operation: "EQ",
		},
	}, sessions, 50)

	if err != nil {
		return nil, err
	}

	sessions = result.(*[]Session)
	return sessions, nil
}

func (s sessionClient) ListExpiredSessions() (*[]Session, error) {

	sessions := &[]Session{}

	result, err := Dynamo.findItems(s.sessionTable, "expires <= :now", map[string]string{
		":now": strconv.Itoa(int(time.Now().Unix())),
	}, sessions, 50)

	if err != nil {
		return nil, err
	}

	sessions = result.(*[]Session)
	return sessions, nil
}

func NewSessionClient(sessionTable string) *sessionClient {
	sessionClient := new(sessionClient)
	sessionClient.sessionTable = sessionTable

	return sessionClient
}
