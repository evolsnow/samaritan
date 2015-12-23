package handler

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mholt/binding"
	"log"
	"net/http"
)

type requestData struct {
	Cardno     int
	Methodname string
}

func (rd *requestData) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&rd.Cardno:     "cardNo",
		&rd.Methodname: "methodName",
	}
}

func bind(r *http.Request) *requestData {
	rd := new(requestData)
	err := binding.Bind(r, rd)
	if err != nil {
		log.Println(err)
	}
	return rd
}

func Hhhh(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rd := bind(r)
	log.Println(rd.Cardno + 1)
}
