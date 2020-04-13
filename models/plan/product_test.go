package plan

import (
	"encoding/json"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/go-rest/rand"
	"github.com/guregu/null"
	"testing"
)

var productStd = ProductSchema{
	BaseProduct: BaseProduct{
		ID:           "prod_IxN4111S1TIP",
		Tier:         enum.TierStandard,
		Heading:      "标准会员",
		SmallPrint:   null.String{},
		YearlyPlanID: "plan_ICMPPM0UXcpZ",
	},
	Description: "专享订阅内容每日仅需0.7元(或按月订阅每日0.9元)\r\n精选深度分析\r\n中英双语内容\r\n金融英语速读训练\r\n英语原声电台\r\n无限浏览7日前所有历史文章（近8万篇）",
}

var productPrm = ProductSchema{
	BaseProduct: BaseProduct{
		ID:           "prod_dcHBCHaBTn3w",
		Tier:         enum.TierPremium,
		Heading:      "高端会员",
		SmallPrint:   null.StringFrom("注：所有活动门票不可折算现金、不能转让、不含差旅与食宿"),
		YearlyPlanID: "plan_5iIonqaehig4",
	},
	Description: "专享订阅内容每日仅需5.5元\r\n享受“标准会员”所有权益\r\n编辑精选，总编/各版块主编每周五为您推荐本周必读资讯，分享他们的思考与观点\r\nFT中文网2018年度论坛门票2张，价值3999元/张 （不含差旅与食宿）",
}

func TestProductID(t *testing.T) {
	t.Logf("Product ID: prod_%s", rand.String(12))
}

func TestZipProductWithPlan(t *testing.T) {
	products := ZipProductWithPlan(
		[]ProductSchema{
			productStd,
			productPrm,
		},
		activePlans,
	)

	b, err := json.MarshalIndent(products, "", "\t")
	if err != nil {
		t.Error(err)
	}

	t.Logf("Products: %s", b)
}
