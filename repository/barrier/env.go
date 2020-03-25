package barrier

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var logger = logrus.WithField("package", "repository")

type Env struct {
	db *sqlx.DB
}

func (env Env) Login() (string, error) {

}
