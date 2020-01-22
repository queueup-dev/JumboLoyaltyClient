package JumboLoyaltyClient

type sessionClient struct {
	sessionTable string

	db DatabaseDriver
}

func (s sessionClient) init(sessionTable string) {
	s.db = NewDynamoDatabase()

	s.sessionTable = sessionTable
}

func (s sessionClient) GetSession(sessionId string) (*session, error) {
	result, err := s.db.getItem(s.sessionTable, "session_id", sessionId, &session{})

	if err != nil {
		return nil, err
	}

	session := result.(session)
	return &session, nil
}

func NewSessionClient(sessionTable string) *sessionClient {
	sessionClient := new(sessionClient)
	sessionClient.init(sessionTable)

	return sessionClient
}
