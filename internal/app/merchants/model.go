package merchants

import "time"

type Merchant struct {
	ID        int64
	Name      string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
