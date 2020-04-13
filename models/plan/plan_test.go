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
	PlanID:    "plan_ICMPPM0UXcpZ",
	Price:     258,
	Tier:      enum.TierStandard,
	Cycle:     enum.CycleYear,
	TrialDays: 3,
}

var prmPlan = BasePlan{
	PlanID:    "plan_5iIonqaehig4",
	Price:     1998,
	Tier:      enum.TierPremium,
	Cycle:     enum.CycleYear,
	TrialDays: 3,
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

	b, err := json.Marshal(stdA)
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

	b, err = json.Marshal(stdB)
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

	b, err = json.Marshal(prm)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Premium plan: %s", b)
}

func TestReduceDiscountPlan(t *testing.T) {
	type args struct {
		rows []DiscountPlanSchema
	}
	tests := []struct {
		name string
		args args
		want Plan
	}{
		{
			name: "Plan without discount",
			args: args{
				rows: []DiscountPlanSchema{
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
				rows: []DiscountPlanSchema{
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
			got := ReduceDiscountPlan(tt.args.rows)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReduceDiscountPlan() = %v, want %v", got, tt.want)
			}

			b, err := json.Marshal(got)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%s", b)
		})
	}
}

func TestGroupDiscountPlans(t *testing.T) {
	type args struct {
		rows []DiscountPlanSchema
	}
	tests := []struct {
		name string
		args args
		want GroupedPlans
	}{
		{
			name: "Group discounts with plans",
			args: args{
				rows: []DiscountPlanSchema{
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
			got := GroupDiscountPlans(tt.args.rows)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupDiscountPlans() = %v, want %v", got, tt.want)
			}

			b, err := json.Marshal(got)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%s", b)
		})
	}
}
