package cleanup

import (
	"context"
	"testing"
)

func TestCleanup(t *testing.T) {

	m := PubSubMessage{
		Data: []byte("instance templates"),
	}
	Cleanup(context.Background(), m)
}
