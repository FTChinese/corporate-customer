package subsrepo

import "github.com/FTChinese/ftacademy/internal/pkg/admin"

func (env Env) AdminProfile(id string) (admin.Profile, error) {
	var p admin.Profile
	err := env.DBs.Read.Get(&p, admin.StmtProfile, id)
	if err != nil {
		return admin.Profile{}, err
	}

	return p, nil
}
