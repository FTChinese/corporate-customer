package plan

import (
	"encoding/json"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/go-rest/rand"
	"reflect"
	"testing"
)

func TestGeneratePlanID(t *testing.T) {

	t.Logf("Plan id: plan_%s", rand.String(12))
}

var stdPlan = BasePlan{
	PlanID: "plan_ICMPPM0UXcpZ",
	Price:  258,
	Tier:   enum.TierStandard,
	Cycle:  enum.CycleYear,
}

var prmPlan = BasePlan{
	PlanID: "plan_5iIonqaehig4",
	Price:  1998,
	Tier:   enum.TierPremium,
	Cycle:  enum.CycleYear,
}

var stdDiscountA = Discount{
	Quantity: 10,
	PriceOff: 15,
}

var stdDiscountB = Discount{
	Quantity: 20,
	PriceOff: 25,
}

var prmDiscountA = Discount{
	Quantity: 10,
	PriceOff: 100,
}

var prmDiscountB = Discount{
	Quantity: 20,
	PriceOff: 200,
}

var activePlans = map[string]Plan{
	stdPlan.PlanID: {
		BasePlan: stdPlan,
		Discounts: []Discount{
			stdDiscountA,
			stdDiscountB,
		},
	},
	prmPlan.PlanID: {
		BasePlan: prmPlan,
		Discounts: []Discount{
			prmDiscountA,
			prmDiscountB,
		},
	},
}

func TestExamplePlans(t *testing.T) {
	stdA := Plan{
		BasePlan: stdPlan,
	}

	b, err := json.MarshalIndent(stdA, "", "\t")
	if err != nil {
		t.Error(err)
	}
	t.Logf("Standard plan without discount: %s", b)

	stdB := Plan{
		BasePlan: stdPlan,
		Discounts: []Discount{
			stdDiscountA,
			stdDiscountB,
		},
	}

	b, err = json.MarshalIndent(stdB, "", "\t")
	if err != nil {
		t.Error(err)
	}

	t.Logf("Standard plan with discount: %s", b)

	prm := Plan{
		BasePlan: prmPlan,
		Discounts: []Discount{
			prmDiscountA,
			prmDiscountB,
		},
	}

	b, err = json.MarshalIndent(prm, "", "\t")
	if err != nil {
		t.Error(err)
	}
	t.Logf("Premium plan: %s", b)
}

func TestNewPlan(t *testing.T) {
	type args struct {
		rows []DiscountPlan
	}
	tests := []struct {
		name string
		args args
		want Plan
	}{
		{
			name: "Plan without discount",
			args: args{
				rows: []DiscountPlan{
					{
						BasePlan: stdPlan,
					},
				},
			},
			want: Plan{
				BasePlan:  stdPlan,
				Discounts: nil,
			},
		},
		{
			name: "Plan with discounts",
			args: args{
				rows: []DiscountPlan{
					{
						BasePlan: stdPlan,
						Discount: stdDiscountA,
					},
					{
						BasePlan: stdPlan,
						Discount: stdDiscountB,
					},
				},
			},
			want: Plan{
				BasePlan: stdPlan,
				Discounts: []Discount{
					stdDiscountA,
					stdDiscountB,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPlan(tt.args.rows)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPlan() = %v, want %v", got, tt.want)
			}

			b, err := json.MarshalIndent(got, "", "\t")
			if err != nil {
				t.Error(err)
			}
			t.Logf("%s", b)
		})
	}
}

func TestNewGroupedPlans(t *testing.T) {
	type args struct {
		rows []DiscountPlan
	}
	tests := []struct {
		name string
		args args
		want GroupedPlans
	}{
		{
			name: "Group discounts with plans",
			args: args{
				rows: []DiscountPlan{
					{
						BasePlan: stdPlan,
						Discount: stdDiscountA,
					},
					{
						BasePlan: stdPlan,
						Discount: stdDiscountB,
					},
					{
						BasePlan: prmPlan,
						Discount: prmDiscountA,
					},
					{
						BasePlan: prmPlan,
						Discount: prmDiscountB,
					},
				},
			},
			want: GroupedPlans{
				stdPlan.PlanID: Plan{
					BasePlan: stdPlan,
					Discounts: []Discount{
						stdDiscountA,
						stdDiscountB,
					},
				},
				prmPlan.PlanID: Plan{
					BasePlan: prmPlan,
					Discounts: []Discount{
						prmDiscountA,
						prmDiscountB,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGroupedPlans(tt.args.rows)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGroupedPlans() = %v, want %v", got, tt.want)
			}

			b, err := json.MarshalIndent(got, "", "\t")
			if err != nil {
				t.Error(err)
			}
			t.Logf("%s", b)
		})
	}
}

func TestPlan_FindDiscount(t *testing.T) {
	type fields struct {
		BasePlan  BasePlan
		Discounts []Discount
	}

	f := fields{
		BasePlan: stdPlan,
		Discounts: []Discount{
			stdDiscountA,
			stdDiscountB,
		},
	}
	type args struct {
		q int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Discount
	}{
		{
			name:   "5 copies",
			fields: f,
			args:   args{q: 5},
			want:   Discount{},
		},
		{
			name:   "15 copies",
			fields: f,
			args:   args{q: 15},
			want:   stdDiscountA,
		},
		{
			name:   "25 copies",
			fields: f,
			args:   args{q: 25},
			want:   stdDiscountB,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Plan{
				BasePlan:  tt.fields.BasePlan,
				Discounts: tt.fields.Discounts,
			}
			if got := p.FindDiscount(tt.args.q); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindDiscount() = %v, want %v", got, tt.want)
			}
		})
	}
}
