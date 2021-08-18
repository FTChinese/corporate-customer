package subsrepo

import (
	"database/sql"
	"errors"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
)

func (env Env) RetrieveAssignee(id string) (licence.Assignee, error) {
	var a licence.Assignee
	err := env.dbs.Read.Get(&a, licence.StmtAssigneeByID, id)
	if err != nil {
		return licence.Assignee{}, err
	}

	return a, nil
}

// FindAssignee tries to find a user by email.
// If user it not found, it is not taken as error and
// an zero value of Assignee is returned.
func (env Env) FindAssignee(email string) (licence.Assignee, error) {
	var a licence.Assignee
	err := env.dbs.Read.Get(&a, licence.StmtAssigneeByEmail, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return licence.Assignee{}, nil
		}
		return licence.Assignee{}, err
	}

	return a, nil
}

func (env Env) RetrieveMembership(compoundID string) (reader.Membership, error) {
	var m reader.Membership

	err := env.dbs.Read.Get(
		&m,
		reader.StmtLockMember,
		compoundID,
	)

	if err != nil && err != sql.ErrNoRows {
		return m, err
	}

	// Treat a non-existing member as a valid value.
	return m.Sync(), nil
}

func (env Env) ArchiveMembership(m reader.MemberSnapshot) error {
	_, err := env.dbs.Write.NamedExec(reader.StmtArchiveMembership, m)
	if err != nil {
		return err
	}

	return nil
}
