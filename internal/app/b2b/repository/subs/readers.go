package subs

import (
	"database/sql"
	"errors"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/guregu/null"
)

func (env Env) RetrieveAssignee(id string) (licence.Assignee, error) {
	var a licence.Assignee
	err := env.dbs.Read.Get(&a, licence.StmtAssigneeByID, id)
	if err != nil {
		return licence.Assignee{}, err
	}

	return a, nil
}

func (env Env) FindAssignee(email string) (licence.Assignee, error) {
	var a licence.Assignee
	err := env.dbs.Read.Get(&a, licence.StmtAssigneeByEmail, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return licence.Assignee{
				FtcID:    null.String{},
				UnionID:  null.String{},
				Email:    null.StringFrom(email),
				UserName: null.String{},
			}, nil
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
