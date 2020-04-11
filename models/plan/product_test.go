package plan

import (
	"github.com/FTChinese/go-rest/rand"
	"testing"
)

func TestProductID(t *testing.T) {
	t.Logf("Product ID: prod_%s", rand.String(12))
}
