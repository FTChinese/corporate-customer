package reader

import "github.com/guregu/null"

// UserIDs is used to identify an FTC user.
// A user might have an ftc uuid, or a wechat union id,
// or both.
// This type structure is used to ensure unique constraint
// for SQL columns that cannot be both null since SQL do not
// have a mechanism to do UNIQUE INDEX on two columns while
// keeping either of them nullable.
// A user's compound id is taken from either ftc uuid or
// wechat id, with ftc id taking precedence.
type UserIDs struct {
	CompoundID string      `json:"compoundId" db:"compound_id"`
	FtcID      null.String `json:"ftcId" db:"ftc_id"`
	UnionID    null.String `json:"unionId" db:"union_id"`
}
