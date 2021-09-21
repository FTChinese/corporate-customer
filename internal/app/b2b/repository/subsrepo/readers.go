package subsrepo

import (
	"database/sql"
	"errors"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
)

// FindAssignee tries to find a user by email.
// If user it not found, it is not taken as error and
// a zero value of Assignee is returned.
func (env Env) FindAssignee(email string) (licence.Assignee, error) {
	var a licence.Assignee
	err := env.DBs.Read.Get(&a, licence.StmtAssigneeByEmail, email)
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

	err := env.DBs.Read.Get(
		&m,
		reader.StmtSelectMember,
		compoundID,
	)

	if err != nil && err != sql.ErrNoRows {
		if err == sql.ErrNoRows {
			return reader.Membership{}, nil
		}

		return m, err
	}

	// Treat a non-existing member as a valid value.
	return m.Sync(), nil
}
