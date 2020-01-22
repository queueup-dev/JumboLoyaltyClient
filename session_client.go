package JumboLoyaltyClient

import "time"

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

func (s sessionClient) ListExpiredSessions() (*[]session, error) {

	condition := QueryCondition{
		"expires",
		string(time.Now().Unix()),
		"LT",
	}

	conditions := []QueryCondition{
		condition,
	}

	result, err := Dynamo.listItems(s.sessionTable, "expires-index", conditions, &[]session{}, 50)

	if err != nil {
		return nil, err
	}

	session := result.(*[]session)
	return session, nil
}

func NewSessionClient(sessionTable string) *sessionClient {
	sessionClient := new(sessionClient)
	sessionClient.sessionTable = sessionTable

	return sessionClient
}
