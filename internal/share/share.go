package share

import (
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
)

// this defines the core entity of the application
// which is the Share entity, defining the properties in which will be interacted with the database and the application
type Share struct {
	ID        uuid.UUID `json:"id"`
	Note      string    `json:"note"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
	IP        net.IP    `json:"ip"`
}

func NewShare(url, note string, ip net.IP) (Share, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return Share{}, fmt.Errorf("could not generate UUID: %w", err)
	}

	return Share{
		ID:        id,
		Note:      note,
		CreatedAt: time.Now(),
		UpdateAt:  time.Now(),
		URL:       url,
		IP:        ip,
	}, nil
}
