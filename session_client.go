package JumboLoyaltyClient

import (
	"strconv"
	"time"
)

type sessionClient struct {
	sessionTable string
}

func (s sessionClient) GetSession(sessionId string) (*session, error) {
	result, err := Dynamo.getItem(s.sessionTable, "session_id", sessionId, &session{})

	if err != nil {
		return nil, err
	}

	session := result.(*session)
	return session, nil
}

func (s sessionClient) DeleteSession(sessionId string) error {
	return Dynamo.deleteItem(s.sessionTable, "session_id", sessionId)
}

func (s sessionClient) ListSessionsOfExternalId(externalId string) (*[]session, error) {

	sessions := &[]session{}

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

	sessions = result.(*[]session)
	return sessions, nil
}

func (s sessionClient) ListExpiredSessions() (*[]session, error) {

	sessions := &[]session{}

	result, err := Dynamo.findItems(s.sessionTable, "expires <= :now", map[string]string{
		":now": strconv.Itoa(int(time.Now().Unix())),
	}, sessions, 50)

	if err != nil {
		return nil, err
	}

	sessions = result.(*[]session)
	return sessions, nil
}

func NewSessionClient(sessionTable string) *sessionClient {
	sessionClient := new(sessionClient)
	sessionClient.sessionTable = sessionTable

	return sessionClient
}
