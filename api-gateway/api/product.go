package api

import (
	"api-gateway/models"
	ppb "api-gateway/protos/product-proto"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
)

type PorductClient struct {
	Product ppb.ProductServiceClient
}

func NewProductClient(pr ppb.ProductServiceClient) *PorductClient {
	return &PorductClient{Product: pr}
}

func (p *PorductClient) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var prodReq ppb.ProductRequest

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "POST: bodydan o'qib olishda xatolik...", http.StatusBadRequest)
		return
	}
	protojson.Unmarshal(bytes, &prodReq)

	resp, err := p.Product.CreateProduct(r.Context(), &prodReq)
	if err != nil {
		log.Println("POST: Serverdan ma'lumot olishda xatolik...", err)
		http.Error(w, "POST: Serverdan ma'lumot olishda xatolik...", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Ma'lumotni encode qilishda xatolik...", http.StatusInternalServerError)
		return
	}
}

func (p *PorductClient) GetProductID(w http.ResponseWriter, r *http.Request) {
	var prodReq ppb.GetProductRequest

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "GET: bodydan o'qib olishda xatolik...", http.StatusBadRequest)
		return
	}
	protojson.Unmarshal(bytes, &prodReq)

	resp, err := p.Product.GetProductID(r.Context(), &prodReq)
	if err != nil {
		log.Println("GET: Serverdan ma'lumot olishda xatolik...", err)
		http.Error(w, "GET: Serverdan ma'lumot olishda xatolik...", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Ma'lumotni encode qilishda xatolik...", http.StatusInternalServerError)
		return
	}
}

func (p *PorductClient) ListProducts(w http.ResponseWriter, r *http.Request) {
	stream, err := p.Product.ListProducts(r.Context(), &ppb.Empty{})
	if err != nil {
		log.Println("GETS: Serverdan ma'lumot olishda xatolik...", err)
		http.Error(w, "GETS: Serverdan ma'lumot olishda xatolik...", http.StatusInternalServerError)
		return
	}
	var listProducts []models.ProductResponseSt
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		product := models.ProductResponseSt{
			ID:        resp.Id,
			Name:      resp.Name,
			Quantity:  resp.Quantity,
			Price:     resp.Price,
			CreatedAt: resp.CreatedAt,
		}
		listProducts = append(listProducts, product)
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(listProducts); err != nil {
		w.Write([]byte("Ma'lumotni encode qilishda xatolik..."))
		log.Println("Ma'lumotni encode qilishda xatolik...")
	}
}
