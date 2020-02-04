package triper

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

// GenerateUUID returns an ULID id
func GenerateUUID() string {
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	entropy := rand.New(source)

	t := time.Now()
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
