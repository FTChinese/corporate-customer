//go:build !production
// +build !production

package mock

import (
	"github.com/FTChinese/ftacademy/internal/pkg/ids"
	"github.com/FTChinese/ftacademy/internal/pkg/input"
	"github.com/FTChinese/ftacademy/internal/pkg/licence"
	"github.com/FTChinese/ftacademy/internal/pkg/reader"
	"github.com/FTChinese/ftacademy/pkg/addon"
	"github.com/FTChinese/ftacademy/pkg/faker"
	"github.com/FTChinese/ftacademy/pkg/price"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"time"
)

type ReaderAccount struct {
	licence.Assignee
	Password string `json:"password" db:"password"`
}

type Persona struct {
	kind     enum.AccountKind
	ftcID    string
	unionID  string
	email    string
	userName string
	password string
}

func NewPersona() Persona {
	faker.SeedGoFake()

	return Persona{
		kind:     enum.AccountKindFtc,
		ftcID:    uuid.New().String(),
		unionID:  faker.GenWxID(),
		email:    gofakeit.Email(),
		userName: gofakeit.Username(),
		password: faker.SimplePassword(),
	}
}

func (p Persona) IsEmpty() bool {
	return p.ftcID == ""
}

func (p Persona) Account() ReaderAccount {
	return ReaderAccount{
		Assignee: p.Assignee(),
		Password: p.password,
	}
}

func (p Persona) Assignee() licence.Assignee {
	return licence.Assignee{
		FtcID:    null.StringFrom(p.ftcID),
		UnionID:  null.String{},
		Email:    null.StringFrom(p.email),
		UserName: null.StringFrom(p.userName),
	}
}

func (p Persona) SignupParams() input.SignupParams {
	return input.SignupParams{
		Credentials: input.Credentials{
			Email:    p.email,
			Password: p.password,
		},
		SourceURL: "",
	}
}

// MemberBuilder creates a default builder that will
// create a standard AnteChange purchase via alipay that
// will expire a month later.
func (p Persona) MemberBuilder(k enum.AccountKind) MemberBuilder {
	return MemberBuilder{
		accountKind:  k,
		ftcID:        p.ftcID,
		unionID:      p.unionID,
		price:        price.MockPriceStdYear,
		payMethod:    enum.PayMethodAli,
		expiration:   time.Now().AddDate(0, 7, 0),
		subsStatus:   0,
		autoRenewal:  false,
		addOn:        addon.AddOn{},
		iapTxID:      "",
		stripeSubsID: "",
		b2bLicID:     "",
	}
}

// MemberBuilderFTC creates a ftc-only AnteChange.
func (p Persona) MemberBuilderFTC() MemberBuilder {
	return p.MemberBuilder(enum.AccountKindFtc)
}

type MemberBuilder struct {
	accountKind  enum.AccountKind
	ftcID        string
	unionID      string
	price        price.Price
	payMethod    enum.PayMethod
	expiration   time.Time
	subsStatus   enum.SubsStatus
	autoRenewal  bool
	addOn        addon.AddOn
	iapTxID      string
	stripeSubsID string
	b2bLicID     string
}

func (b MemberBuilder) WithAccountKind(k enum.AccountKind) MemberBuilder {
	b.accountKind = k
	return b
}

func (b MemberBuilder) WithIDs(ids reader.UserIDs) MemberBuilder {
	b.ftcID = ids.FtcID.String
	b.unionID = ids.UnionID.String
	return b
}

func (b MemberBuilder) WithFtcID(id string) MemberBuilder {
	b.ftcID = id
	return b
}

func (b MemberBuilder) WithWxID(id string) MemberBuilder {
	b.unionID = id
	return b
}

func (b MemberBuilder) WithPrice(p price.Price) MemberBuilder {
	b.price = p

	return b
}

// WithPayMethod sets payment method.
// Deprecated.
func (b MemberBuilder) WithPayMethod(m enum.PayMethod) MemberBuilder {
	b.payMethod = m
	if m == enum.PayMethodStripe || m == enum.PayMethodApple {
		b.autoRenewal = true
		b.subsStatus = enum.SubsStatusActive
	}
	return b
}

func (b MemberBuilder) WithStripe(subsID string) MemberBuilder {
	if subsID == "" {
		subsID = faker.GenStripeSubID()
	}

	b.stripeSubsID = subsID
	return b.WithPayMethod(enum.PayMethodStripe)
}

func (b MemberBuilder) WithApple(txID string) MemberBuilder {
	if txID == "" {
		txID = faker.GenAppleSubID()
	}

	b.iapTxID = txID
	b.autoRenewal = true
	return b.WithPayMethod(enum.PayMethodApple)
}

func (b MemberBuilder) WithB2B(licID string) MemberBuilder {
	if licID == "" {
		licID = ids.LicenceID()
	}

	b.b2bLicID = licID
	return b.WithPayMethod(enum.PayMethodB2B)
}

func (b MemberBuilder) WithExpiration(t time.Time) MemberBuilder {
	b.expiration = t
	return b
}

// WithAutoRenewal turn auto renew on/off
// Deprecated.
func (b MemberBuilder) WithAutoRenewal(t bool) MemberBuilder {
	b.autoRenewal = t
	return b
}

func (b MemberBuilder) WithAutoRenewOn() MemberBuilder {
	b.autoRenewal = true
	return b
}

func (b MemberBuilder) WithAutoRenewOff() MemberBuilder {
	b.autoRenewal = false
	return b
}

func (b MemberBuilder) WithSubsStatus(s enum.SubsStatus) MemberBuilder {
	b.subsStatus = s
	return b
}

func (b MemberBuilder) WithAddOn(r addon.AddOn) MemberBuilder {
	b.addOn = r
	return b
}

func (b MemberBuilder) WithIapID(id string) MemberBuilder {
	b.iapTxID = id
	return b
}

func (b MemberBuilder) Build() reader.Membership {
	var userIDs reader.UserIDs
	switch b.accountKind {
	case enum.AccountKindFtc:
		userIDs = reader.UserIDs{
			CompoundID: b.ftcID,
			FtcID:      null.StringFrom(b.ftcID),
			UnionID:    null.String{},
		}
	case enum.AccountKindWx:
		userIDs = reader.UserIDs{
			CompoundID: b.unionID,
			FtcID:      null.String{},
			UnionID:    null.StringFrom(b.unionID),
		}
	case enum.AccountKindLinked:
		userIDs = reader.UserIDs{
			CompoundID: b.ftcID,
			FtcID:      null.StringFrom(b.ftcID),
			UnionID:    null.StringFrom(b.unionID),
		}
	}

	m := reader.Membership{
		UserIDs:       userIDs,
		Edition:       b.price.Edition,
		LegacyTier:    null.Int{},
		LegacyExpire:  null.Int{},
		ExpireDate:    chrono.DateFrom(b.expiration),
		PaymentMethod: b.payMethod,
		FtcPlanID:     null.String{},
		StripeSubsID:  null.String{},
		StripePlanID:  null.String{},
		AutoRenewal:   b.autoRenewal,
		Status:        b.subsStatus,
		AppleSubsID:   null.String{},
		B2BLicenceID:  null.String{},
		AddOn:         b.addOn,
	}

	switch b.payMethod {
	case enum.PayMethodAli, enum.PayMethodWx:
		m.FtcPlanID = null.StringFrom(b.price.ID)

	case enum.PayMethodStripe:
		m.StripeSubsID = null.StringFrom(b.stripeSubsID)
		m.StripePlanID = null.StringFrom(faker.GenStripePlanID())

	case enum.PayMethodApple:
		m.AppleSubsID = null.StringFrom(b.iapTxID)

	case enum.PayMethodB2B:
		m.B2BLicenceID = null.StringFrom(b.b2bLicID)
	}

	return m.Sync()
}
