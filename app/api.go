package main

import (
	"database/sql"
	"github.com/ant0ine/go-json-rest/rest"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(2)
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/", ApiHandler),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}

type Product struct {
	Amount string `json:"amount"`
	Name   string `json:"name"`
	Des    string `json:"des"`
}

func ApiHandler(w rest.ResponseWriter, r *rest.Request) {
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
	for rows.Next() {
		err := rows.Scan(&p.Name, &p.Amount)
		p.Des = "hell哦"
		s = append(s, *p)
		if err != nil {
			log.Fatal(err)
		}
	}

	//	map[string]string{"Body": "Hello World!"}
	w.WriteJson(s)
}
