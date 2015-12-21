package main

import (
	"database/sql"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type product struct {
	Id              int
	Name            string
	CreateName      string
	ReviewTime      string
	PublishTime     string
	InterestStart   string
	InterestEnd     string
	Period          uint16
	FundsAmount     uint
	MinPurchase     uint
	InterestYear    float32
	RepaymentType   uint8
	Category        uint8
	Status          uint8
	IntegerRequired uint8
	RequiredLv      uint8
	BorrowId        uint
}

func syncProduct(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	db, err := sql.Open("mysql", "root:6:u2WkG=B:@tcp(139.196.13.186:3306)/sqd?autocommit=true")
	if err != nil {
		log.Fatalf("Open database error: %s\n", err)
	}
	defer db.Close()

	// Prepare statement for reading data
	conn := pool.Get()
	id, _ := redis.Int(conn.Do("INCR", "incr_product_id"))
	rows, err := db.Query("SELECT * FROM `product`")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan()
	}

}
