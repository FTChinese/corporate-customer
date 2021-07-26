package licence

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/FTChinese/ftacademy/internal/pkg"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/rand"
	"github.com/guregu/null"
	"time"
)

// Invitation is an email sent to team member to accept a licence.
// An invitation could in 3 phases:
// Initially created: it indicates an email is sent to reader;
// Accepted: reader clicked the link in the invitation email,
// it should not be used any longer;
// Revoked: admin could revoke an invitation before it is accepted.
// An accepted invitation could not be revoked since that is meaningless.
// TODO: should we allow an invitation be re-sent if user failed to receive the email? Or just let admin to create a new invitation?
type Invitation struct {
	ID string `json:"id" db:"invite_id"`
	Creator
	Description    null.String      `json:"description" db:"invite_desc"`
	ExpirationDays int64            `json:"expirationDays" db:"invite_expiration_days"`
	Email          string           `json:"email" db:"invite_email"`
	LicenceID      string           `json:"licenceId" db:"licence_id"`
	Status         InvitationStatus `json:"status" db:"invite_current_status"`
	Token          string           `json:"-" db:"invite_token"` // This field is used only when inserting data. Retrieval does not include this field. However, it is included when saving to the JSON column in licence.
	RowMeta
}

func NewInvitation(params input.InvitationParams, adminID string) (Invitation, error) {
	token, err := rand.Hex(32)
	if err != nil {
		return Invitation{}, err
	}

	return Invitation{
		ID: pkg.InvitationID(),
		Creator: Creator{
			CreatorID: adminID,
			TeamID:    params.TeamID,
		},
		Description:    params.Description,
		ExpirationDays: 7,
		Email:          params.Email,
		LicenceID:      params.LicenceID,
		Status:         InvitationStatusCreated,
		Token:          token,
		RowMeta: RowMeta{
			CreatedUTC: chrono.TimeNow(),
			UpdatedUTC: chrono.Time{},
		},
	}, nil
}

// IsExpired tests whether the invitation is expired.
// An expired invitation is not allowed grant its related licence.
func (i Invitation) IsExpired() bool {
	now := time.Now().Unix()

	created := i.CreatedUTC.Time.Unix()

	// Default 7 days * 24 * 60 * 60
	return (created + i.ExpirationDays*86400) < now
}

// IsAcceptable determines whether an invitation is valid.
// A valid invitation must be not expires, not revoked by admin, not accepted by any one.
// A valid invitation can be accepted or revoked.
func (i Invitation) IsAcceptable() bool {
	return i.Status == InvitationStatusCreated && !i.IsExpired()
}

// Accepted invalidates an invitation after reader accepted the licence associated with it.
func (i Invitation) Accepted() Invitation {
	i.Status = InvitationStatusAccepted
	i.UpdatedUTC = chrono.TimeNow()

	return i
}

func (i Invitation) IsRevocable() bool {
	return i.Status == InvitationStatusCreated
}

// Revoked invalidates an invitation by admin.
func (i Invitation) Revoked() Invitation {
	i.Status = InvitationStatusRevoked
	i.UpdatedUTC = chrono.TimeNow()

	return i
}

// Value implements Valuer interface by serializing an Invitation into
// JSON data.
func (i Invitation) Value() (driver.Value, error) {
	if i.ID == "" {
		return nil, nil
	}

	b, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	return string(b), nil
}

// Scan implements Valuer interface by deserializing an invitation field.
func (i *Invitation) Scan(src interface{}) error {
	if src == nil {
		*i = Invitation{}
		return nil
	}

	switch s := src.(type) {
	case []byte:
		var tmp Invitation
		err := json.Unmarshal(s, &tmp)
		if err != nil {
			return err
		}
		*i = tmp
		return nil

	default:
		return errors.New("incompatible type to scan")
	}
}

// InvitationVerified is returned after an invitation link
// is clicked and the corresponding Licence is found.
type InvitationVerified struct {
	Licence    Licence // The licence being invited.
	Assignee   Assignee
	Membership reader.Membership
}

// InvitationList is used for restful output.
type InvitationList struct {
	pkg.PagedList
	Data []Invitation `json:"data"`
}
