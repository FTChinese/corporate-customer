package plan

import (
	"github.com/FTChinese/go-rest/rand"
	"testing"
)

func TestGeneratePlanID(t *testing.T) {

	t.Logf("Plan id: plan_%s", rand.String(12))
}
