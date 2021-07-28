// +build !production

package faker

import (
	"fmt"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/go-rest/rand"
	"github.com/brianvoe/gofakeit/v5"
	"time"
)

func SeedGoFake() {
	gofakeit.Seed(time.Now().UnixNano())
}

// GenVersion creates a semantic version string.
func GenVersion() string {
	return fmt.Sprintf("%d.%d.%d",
		rand.IntRange(1, 10),
		rand.IntRange(1, 10),
		rand.IntRange(1, 10))
}

func RandNumericString() string {
	return rand.StringWithCharset(9, "0123456789")
}

func RandomTier() enum.Tier {
	return enum.Tier(rand.IntRange(1, 3))
}

func RandomGender() enum.Gender {
	return enum.Gender(rand.IntRange(0, 3))
}

func GenLicenceID() string {
	return "lic_" + rand.String(12)
}

func GenPhone() string {
	SeedGoFake()
	return "1" + gofakeit.Phone()
}

func GenEmail() string {
	SeedGoFake()
	return gofakeit.Email()
}

func SimplePassword() string {
	return gofakeit.Password(true, false, true, false, false, 8)
}

func GenCardSerial() string {
	now := time.Now()
	anni := now.Year() - 2005
	suffix := rand.IntRange(0, 9999)

	return fmt.Sprintf("%d%02d%04d", anni, now.Month(), suffix)
}
