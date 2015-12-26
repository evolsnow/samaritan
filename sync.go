package main

import (
	"database/sql"
	//	"fmt"
	"github.com/evolsnow/gosqd/conn"
	"github.com/evolsnow/gosqd/model"
	"github.com/evolsnow/httprouter"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

type customer struct {
}

func syncProduct(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	db, err := sql.Open("mysql", "root:6:u2WkG=B:@tcp(139.196.13.186:3306)/sqd?autocommit=true")
	if err != nil {
		log.Fatalf("Open database error: %s\n", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM `product`")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	p := &model.RawProduct{}
	for rows.Next() {
		err := rows.Scan(&p.Id, &p.Name, &p.CreateTime, &p.ReviewTime, &p.PublishTime, &p.InterestStart, &p.InterestEnd,
			&p.Period, &p.FundsAmount, &p.MinPurchase, &p.SalseAmount, &p.InterestYear, &p.RepaymentType, &p.Category,
			&p.Status, &p.IntegerRequired, &p.RequiredLv, &p.BorrowId, &p.CreateStuffId, &p.ReviewStuffId)
		if err != nil {
			log.Println(err)
		}

		reply, err := redis.Int(conn.CreateProduct(p))
		if err != nil {
			log.Println(err)
		}
		log.Printf("已完成标:%d", reply)

	}
}

func syncCustomer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}
