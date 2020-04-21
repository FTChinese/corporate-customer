package subs

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var logger = logrus.WithField("package", "repository")

type Env struct {
	db *sqlx.DB
}

func NewEnv(db *sqlx.DB) Env {
	return Env{db: db}
}

func (env Env) beginInvTx() (InvitationTx, error) {
	tx, err := env.db.Beginx()

	if err != nil {
		return InvitationTx{}, err
	}

	return InvitationTx{
		tx,
	}, nil
}
