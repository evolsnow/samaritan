package handler

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Hhhh(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "hhhh")
}
