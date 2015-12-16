package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	r := httprouter.New()
	r.GET("/", HomeHandler)
	log.Fatal(http.ListenAndServe(":8080", r))
}

type Product struct {
	Amount int    `json:"amount"`
	Name   string `json:"name"`
	Des    string `json:"des"`
}
type Response struct {
	ResObj  []Product `json:"responseObject"`
	ResCode string    `json:"responseCode"`
	ResMsg  string    `json:"responseMessage"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	db, err := sql.Open("mysql", "root:sqdShengQianDai@tcp(121.43.110.32:3306)/sqd?autocommit=true")
	if err != nil {
		log.Fatalf("Open database error: %s\n", err)
	}
	defer db.Close()

	// Prepare statement for reading data
	rows, err := db.Query("SELECT `name`, `funds_amount` FROM `product`")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()

	var s []Product
	p := &Product{}
	p.Des = "hell哦"

	for rows.Next() {
		err := rows.Scan(&p.Name, &p.Amount)
		s = append(s, *p)
		if err != nil {
			log.Fatal(err)
		}
	}
	ret := Response{ResObj: s, ResCode: "success", ResMsg: "提交成功"}
	Reply(w, r, ret)
}

func Reply(w http.ResponseWriter, r *http.Request, ret interface{}) {
	//ret := Response{ResObj: s, ResCode: "success", ResMsg: "提交成功"}
	js, _ := json.Marshal(&ret)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
