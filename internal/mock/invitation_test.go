package mock

import (
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/FTChinese/go-rest/enum"
	"testing"
)

func TestZeroMembership(t *testing.T) {
	p := NewPersona()
	a := p.Account()

	repo := NewRepo()
	repo.CreateReader(a)

	t.Logf("%s", faker.MustMarshalIndent(a))
}

func TestOneTimeMembership(t *testing.T) {
	p := NewPersona()
	a := p.Account()

	repo := NewRepo()
	repo.CreateReader(a)
	repo.InsertMembership(p.MemberBuilderFTC().Build())

	t.Logf("%s", faker.MustMarshalIndent(a))
}

func TestAppleMembership(t *testing.T) {
	p := NewPersona()
	a := p.Account()

	repo := NewRepo()
	repo.CreateReader(a)
	repo.InsertMembership(p.MemberBuilderFTC().WithPayMethod(enum.PayMethodApple).Build())

	t.Logf("%s", faker.MustMarshalIndent(a))
}

func TestB2BMembership(t *testing.T) {
	p := NewPersona()
	a := p.Account()

	repo := NewRepo()
	repo.CreateReader(a)
	repo.InsertMembership(p.MemberBuilderFTC().WithB2B("lic_nMjtCCNerLNO").Build())

	t.Logf("%s", faker.MustMarshalIndent(a))
}
