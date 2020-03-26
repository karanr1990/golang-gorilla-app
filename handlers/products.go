package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"github.ibm.com/Quest-CIO/go-micro-app/data"
	"log"
	"net/http"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request)  {
	p.l.Println("Handle GET product")
	lp :=data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil{
		http.Error(rw,"Unable to marshal json",http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request)  {
	p.l.Println("Handle post product")

	prod := &data.Product{}

    err:=prod.FromJSON(r.Body)

    if err !=nil {
    	http.Error(rw,"unable to unmarshal json", http.StatusBadRequest)
	}
    data.AddProduct(prod)
}

func (p Products)UpdateProducts(rw http.ResponseWriter, r *http.Request)  {

	vars :=mux.Vars(r)
	id,err := strconv.Atoi(vars["id"])

	if err !=nil {
		http.Error(rw,"unable to convert id ",http.StatusBadRequest)
	}
	p.l.Println("Handle PUT product",id)
	prod := r.Context().Value(KeyProduct{}).(data.Product)


	err = data.UpdateProduct(id, &prod)

    if err == data.ErrorProductNotFound{
    	http.Error(rw,"product not found",http.StatusNotFound)
		return
	}
	if err !=nil{
		http.Error(rw,"product not found",http.StatusInternalServerError)
		return
	}
	
}

type KeyProduct struct{}

func(p Products) MiddlewareProductValidation(next http.Handler) http.Handler  {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request){
		prod := data.Product{}
		err := prod.FromJSON(r.Body)

		if err != nil{
			http.Error(rw,"Unable to unmarshal the JSON",http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(),KeyProduct{},prod)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw,req)

	})
}
