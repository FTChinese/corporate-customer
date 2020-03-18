package reader

// Profile contains the data of a reader belong to an organization.
type Profile struct {
	FtcID string `json:"ftcId" db:"ftc_id"`
	Email string `json:"email" db:"email"`
	Membership
}
