package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

type Product struct {
	Id              int    `json:"id"`
	FundsAmount     int    `json:"account"`
	SalesAmount     int    `json:"accountYes"`
	Available       int    `json:"navailable"`
	InterestYear    string `json:"expectApr"`
	MinPurchase     int    `json:"lowestAccount"`
	InterestStart   string `json:"qxDate"`
	Name            string `json:"name"`
	Period          string `json:"timeLimit"`
	CompletePercent string `json:"completePercent"`
	StatusStr       string `json:"statusStr"`
	StatusInt       int
	ZhiDing         int    `json:"isZhiDing"`
	Day             int    `json:"isday"`
	AllStatus       string `json:"allStatus"`
}

const (
	STATUS_CREATED        = 0
	STATUS_REVIEW_SUCCESS = 1
	STATUS_REVIEW_FAILED  = 2
	STATUS_PUBLISHED      = 3
	STATUS_SOLD_OUT       = 4
	STATUS_INTEREST_START = 5
	STATUS_REPAID         = 6
	STATUS_REMOVED        = 7
	STATUS_REFUND         = 8
)

var ProductStatus = map[int]string{
	STATUS_CREATED:        "已创建",
	STATUS_REVIEW_SUCCESS: "审核通过",
	STATUS_REVIEW_FAILED:  "审核未通过",
	STATUS_PUBLISHED:      "热销",
	STATUS_SOLD_OUT:       "已售罄",
	STATUS_INTEREST_START: "还款中",
	STATUS_REPAID:         "已还款",
	STATUS_REMOVED:        "已撤标",
	STATUS_REFUND:         "已退款",
}

type Response struct {
	ResObj  []Product `json:"responseObject"`
	ResCode string    `json:"responseCode"`
	ResMsg  string    `json:"responseMessage"`
}

func ProductList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	db, err := sql.Open("mysql", "root:sqdShengQianDai@tcp(121.43.110.32:3306)/sqd?autocommit=true")
	if err != nil {
		log.Fatalf("Open database error: %s\n", err)
	}
	defer db.Close()

	// Prepare statement for reading data
	rows, err := db.Query("SELECT `id`, `funds_amount`,`sales_amount`, `interest_year`,`min_purchase`, `interest_start`,`name`, `period`,`status` FROM `product`")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var s []Product
	p := &Product{ZhiDing: 1, Day: 1, AllStatus: "1"}

	for rows.Next() {
		err := rows.Scan(&p.Id, &p.FundsAmount, &p.SalesAmount, &p.InterestYear, &p.MinPurchase, &p.InterestStart, &p.Name, &p.Period, &p.StatusInt)
		p.Available = p.FundsAmount - p.SalesAmount
		p.CompletePercent = strconv.FormatFloat((float64(p.SalesAmount) / float64(p.FundsAmount) * 100), 'f', 2, 64)
		p.StatusStr = ProductStatus[p.StatusInt]
		s = append(s, *p)
		if err != nil {
			log.Fatal(err)
		}
	}
	ret := Response{ResObj: s, ResCode: "success", ResMsg: "提交成功"}
	reply(w, r, ret)
}

func reply(w http.ResponseWriter, r *http.Request, ret interface{}) {
	js, _ := json.Marshal(&ret)
	w.Write(js)
}
