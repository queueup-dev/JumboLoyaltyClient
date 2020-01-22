package JumboLoyaltyClient

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

func NewSessionClient(sessionTable string) *sessionClient {
	sessionClient := new(sessionClient)
	sessionClient.sessionTable = sessionTable

	return sessionClient
}
