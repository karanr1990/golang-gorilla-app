package data

import (
	"encoding/json"
	"fmt"
	"io"
)

type Product struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Price float32 `json:"price"`
	SKU string `json:"sku"`
	CreatedOn string `json:"-"`
	UpdatedOn string `json:"-"`
	DeletedOn string `json:"-"`
}

type Products []*Product

func(p *Product) FromJSON(r io.Reader) error {
        err := json.NewDecoder(r)
	return err.Decode(p)
}


func(p *Products) ToJSON(w io.Writer) error{
	err := json.NewEncoder(w)
	return err.Encode(p)
	
}
func GetProducts() Products {
	return productList
}

func AddProduct(p *Product)  {
	p.ID = getNextID()
	productList = append(productList,p)
	
}

func getNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

func UpdateProduct(id int,p*Product) error  {
	_,pos,err := findProduct(id)

	if err != nil{
		return nil
	}
	p.ID = id
	productList[pos] = p

	return nil

}

var ErrorProductNotFound = fmt.Errorf("product not found")

func findProduct(id int)(*Product,int,error)  {

	for i,p := range productList {
      if p.ID == id {
      	return p,i,nil
	  }
	}
	return nil,-1 , ErrorProductNotFound
}

var productList = []*Product{
	&Product{
		ID: 1,
		Name: "string",
		Description: "string",
		Price: 32,
		SKU: "string",
		CreatedOn: "string",
		UpdatedOn: "string",
	},
	&Product{
		ID: 2,
		Name: "string",
		Description: "string",
		Price: 32,
		SKU: "string",
		CreatedOn: "string",
		UpdatedOn: "string",
	},

}
