package stream

import (
	"log"

	"github.com/nats-io/nats.go"
	"github.com/ragul28/grpc-event-stream/pkg/utils"
)

func JetStreamInit(StreamName string) (nats.JetStreamContext, error) {
	// Connect to NATS
	nc, err := nats.Connect(utils.GetEnv("NATS_URL", nats.DefaultURL))
	if err != nil {
		return nil, err
	}

	// Create JetStream Context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}

	// Create a stream if it does not exist
	err = CreateStream(js, StreamName)
	if err != nil {
		return nil, err
	}

	return js, nil
}

func CreateStream(jetStream nats.JetStreamContext, StreamName string) error {
	StreamSubjects := StreamName + ".*"

	stream, err := jetStream.StreamInfo(StreamName)
	// stream not found, create it
	if stream == nil {
		log.Printf("Creating stream: %s\n", StreamName)

		_, err = jetStream.AddStream(&nats.StreamConfig{
			Name:     StreamName,
			Subjects: []string{StreamSubjects},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
