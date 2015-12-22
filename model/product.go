package model

type Product struct {
	Id              int
	Name            string
	CreateTime      string
	ReviewTime      string
	PublishTime     string
	InterestStart   string
	InterestEnd     string
	Period          uint16
	FundsAmount     uint
	MinPurchase     uint
	SalseAmount     uint
	InterestYear    float32
	RepaymentType   uint8
	Category        uint8
	Status          uint8
	IntegerRequired uint8
	RequiredLv      uint8
	BorrowId        uint
	CreateStuffId   uint
	ReviewStuffId   uint
}
