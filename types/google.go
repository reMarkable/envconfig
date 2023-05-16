package types

import (
	"errors"
	"regexp"
)

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
