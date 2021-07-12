package admin

import "github.com/FTChinese/ftacademy/internal/pkg/input"

// Profile is used to retrieve base account and team data in one shot.
type Profile struct {
	BaseAccount
	input.TeamParams
}
