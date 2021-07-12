package adminor

import "github.com/FTChinese/ftacademy/internal/pkg/admin"

// SavePwResetSession saves a PwResetSession used to generate
// a password reset email.
func (env Env) SavePwResetSession(s admin.PwResetSession) error {

	_, err := env.DBs.Write.NamedExec(
		admin.StmtInsertPwResetSession,
		s,
	)

	if err != nil {
		return err
	}

	return nil
}

// PwResetSession retrieves PwResetSession by token.
func (env Env) PwResetSession(token string) (admin.PwResetSession, error) {
	var session admin.PwResetSession
	err := env.DBs.Read.Get(&session, admin.StmtPwResetSessionByToken, token)
	if err != nil {
		return admin.PwResetSession{}, err
	}

	return session, nil
}

// DisablePasswordReset disables a token used.
func (env Env) DisablePasswordReset(t string) error {

	_, err := env.DBs.Write.Exec(admin.StmtDisablePwResetToken, t)

	if err != nil {
		return err
	}

	return nil
}
