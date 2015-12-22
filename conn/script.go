package conn

import (
	"github.com/evolsnow/gosqd/model"
)

func CreateProduct(p *model.Product) (interface{}, error) {
	conn := Pool.Get()
	defer conn.Close()
	return CreateProductScript.Do(conn, p.Name, p.CreateTime, p.ReviewTime, p.PublishTime, p.InterestStart, p.InterestEnd,
		p.Period, p.FundsAmount, p.MinPurchase, p.SalseAmount, p.InterestYear, p.RepaymentType, p.Category,
		p.Status, p.IntegerRequired, p.RequiredLv, p.BorrowId, p.CreateStuffId, p.ReviewStuffId)
}
