package repository

import (
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/db"
)

// SharedRepo that can be embedded.
type SharedRepo struct {
	DBs db.ReadWriteMyDBs
}

func NewSharedRepo(dbs db.ReadWriteMyDBs) SharedRepo {
	return SharedRepo{
		DBs: dbs,
	}
}

func (r SharedRepo) RetrieveAssignee(id string) (licence.Assignee, error) {
	var a licence.Assignee
	err := r.DBs.Read.Get(&a, licence.StmtAssigneeByID, id)
	if err != nil {
		return licence.Assignee{}, err
	}

	return a, nil
}

// ArchiveMembership save membership prior to modification.
func (r SharedRepo) ArchiveMembership(m reader.MembershipVersioned) error {
	_, err := r.DBs.Write.NamedExec(reader.StmtVersionMembership, m)
	if err != nil {
		return err
	}

	return nil
}

// SaveVersionedLicence whenever a licence is changed.
// This could happen when a licence is created, renewed,
// granted or revoked.
func (r SharedRepo) SaveVersionedLicence(s licence.Versioned) error {
	_, err := r.DBs.Write.NamedExec(licence.StmtVersioned, s)
	if err != nil {
		return err
	}

	return nil
}
