package licence

import (
	"github.com/FTChinese/go-rest/chrono"
	"time"
)

type Creator struct {
	CreatorID string `json:"creatorId" db:"creator_id"`
	TeamID    string `json:"teamId" db:"team_id"`
}

type RowMeta struct {
	CreatedUTC chrono.Time `json:"createdUtc" db:"created_utc"`
	UpdatedUTC chrono.Time `json:"updatedUtc" db:"updated_utc"`
}

func NewRowMeta() RowMeta {
	now := time.Now().Truncate(time.Second).In(time.UTC)

	return RowMeta{
		CreatedUTC: chrono.TimeUTCFrom(now),
		UpdatedUTC: chrono.TimeUTCFrom(now),
	}
}
