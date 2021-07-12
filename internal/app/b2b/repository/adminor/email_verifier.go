package adminor

import "github.com/FTChinese/ftacademy/internal/pkg/admin"

func (env Env) SaveEmailVerifier(v admin.EmailVerifier) error {
	_, err := env.DBs.Write.NamedExec(admin.StmtInsertEmailVerifier, v)

	if err != nil {
		return err
	}

	return nil
}

func (env Env) RetrieveEmailVerifier(token string) (admin.EmailVerifier, error) {
	var v admin.EmailVerifier
	err := env.DBs.Read.Get(&v, admin.StmtRetrieveEmailVerifier, token)
	if err != nil {
		return admin.EmailVerifier{}, err
	}

	return v, nil
}

func (env Env) EmailVerified(ID string) error {
	_, err := env.DBs.Write.Exec(admin.StmtEmailVerified, ID)

	if err != nil {
		return err
	}

	return nil
}
