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
func (r SharedRepo) ArchiveMembership(m reader.MemberSnapshot) error {
	_, err := r.DBs.Write.NamedExec(reader.StmtArchiveMembership, m)
	if err != nil {
		return err
	}

	return nil
}
