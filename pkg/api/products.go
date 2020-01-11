package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Kamva/mgm"
	"github.com/globalsign/mgo/bson"
)

type GetProductsResponse struct {
	Products []Product `json:"products"`
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	products := []Product{}
	err := mgm.Coll(&Product{}).SimpleFind(&products, bson.M{})
	if err != nil {
		log.Fatalln(err)
	}

	response := GetProductsResponse{Products: products}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatalln("Error marshalling data", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type AddProductResponse struct {
	Product Product `json:"product"`
}

func AddProduct(w http.ResponseWriter, r *http.Request) {
	var product Product

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		RespondWithMessage(w, "Invalid JSON Payload")
		return
	}

	product.CreatedBy = r.Context().Value("userID").(string)

	err = mgm.Coll(&product).Create(&product)
	if err != nil {
		log.Fatalln("Error saving product to DB", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	response := AddProductResponse{Product: product}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatalln("Error marshalling data", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
