package stream

import (
	"github.com/nats-io/nats.go"
)

func JetStreamInit(StreamName string) (nats.JetStreamContext, error) {

	return nil, nil
}

func CreateStream(jetStream nats.JetStreamContext, StreamName string) error {
	stream, err := jetStream.StreamInfo(StreamName)

	// stream not found, create it
	if stream == nil {
		log.Printf("Creating stream: %s\n", StreamName)

		_, err = jetStream.AddStream(&nats.StreamConfig{
			Name:     StreamName,
			Subjects: []string{"main"},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
