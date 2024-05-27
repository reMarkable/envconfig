package types

import (
	"errors"
	"regexp"
)

// -----------------------------------------------------------------------------
// PUBSUB TOPIC
// -----------------------------------------------------------------------------

var (
	// ErrInvalidGoogleTopicID means the configured topic has the wrong format.
	ErrInvalidGoogleTopicID = errors.New("topic is not valid format")

	googleTopicRegexp = regexp.MustCompile(`projects\/([\w-]+)\/topics\/([\w-]+)`)
)

type GooglePubSubTopic struct {
	ProjectID string
	TopicID   string
}

func (pst *GooglePubSubTopic) Set(value string) error {
	m := googleTopicRegexp.FindStringSubmatch(value)
	if len(m) != 3 {
		return ErrInvalidGoogleTopicID
	}

	pst.ProjectID = m[1]
	pst.TopicID = m[2]

	return nil
}

// -----------------------------------------------------------------------------
// FIRESTORE DATABASE
// -----------------------------------------------------------------------------

var (
	// ErrInvalidGoogleFirestoreID means the configured database id has the wrong format.
	ErrInvalidGoogleFirestoreID = errors.New("firestore id is not valid format")

	googleFirestoreRegexp = regexp.MustCompile(`projects\/([\w-]+)\/databases\/([\w-]+)`)
)

type GoogleFirestoreDatabase struct {
	ProjectID string
	Database  string
}

func (pst *GoogleFirestoreDatabase) Set(value string) error {
	m := googleFirestoreRegexp.FindStringSubmatch(value)
	if len(m) != 3 {
		return ErrInvalidGoogleFirestoreID
	}

	pst.ProjectID = m[1]
	pst.Database = m[2]

	return nil
}
