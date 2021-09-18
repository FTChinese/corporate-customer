package licence

import (
	"errors"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

func (l Licence) IsRevocable() bool {
	return l.Status == LicStatusGranted && l.AssigneeID.Valid
}

// Revoked unlink a user from a licence.
func (l Licence) Revoked() Licence {
	l.HintGrantMismatch = false
	l.Status = LicStatusAvailable
	l.LatestInvitation = InvitationJSON{}
	l.AssigneeID = null.String{}
	l.UpdatedUTC = chrono.TimeUTCNow()

	return l
}

type RevokeResult struct {
	LicenceVersion      Versioned                  `json:"licenceVersion"`
	MembershipVersioned reader.MembershipVersioned `json:"membershipVersioned"`
}

func RevokeLicence(lic Licence, mmb reader.Membership) (RevokeResult, error) {
	if !lic.IsGrantedTo(mmb) {
		return RevokeResult{}, errors.New("error revoking licence: membership not generated from this licence")
	}

	mv := mmb.LicenceRevoked().
		LicenceRevoked().
		Version(reader.B2BArchiver(reader.ArchiveActionRevoke)).
		WithPriorVersion(mmb).
		WithB2BTxnID(lic.LatestTransactionID.String)

	return RevokeResult{
		LicenceVersion: lic.Revoked().
			Versioned(VersionActionRevoke).
			WithPriorVersion(lic).
			WithMembershipVersioned(mv.ID),
		MembershipVersioned: mv,
	}, nil
}
