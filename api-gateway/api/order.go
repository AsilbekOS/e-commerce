package api

import (
	"api-gateway/models"
	opb "api-gateway/protos/order-proto"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
)

type OrderClient struct {
	Order opb.OrderServiceClient
}

func NewOrderClient(or opb.OrderServiceClient) *OrderClient {
	return &OrderClient{Order: or}
}

func (o *OrderClient) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var ord opb.OrderRequest

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "POST: bodydan o'qib olishda xatolik...", http.StatusBadRequest)
		return
	}

	err = protojson.Unmarshal(bytes, &ord)
	if err != nil {
		http.Error(w, "POST: JSON ni parse qilishda xatolik...", http.StatusBadRequest)
		return
	}

	resp, err := o.Order.CreateOrder(r.Context(), &ord)
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

func (o *OrderClient) CreateOrders(w http.ResponseWriter, r *http.Request) {
	// var ord opb.OrderRequest
	var listOrder []models.OrderResponseJson

	err := json.NewDecoder(r.Body).Decode(&listOrder)
	if err != nil {
		log.Println("Malumotlarni decode qilishda xatolik...", err)
	}

	stream, err := o.Order.CreateOrders(r.Context())
	if err != nil {
		log.Println("POSTs: Serverdan ma'lumot olishda xatolik...", err)
		http.Error(w, "POSTS: Serverdan ma'lumot olishda xatolik...", http.StatusInternalServerError)
	}

	for _, i := range listOrder {
		order := opb.OrderRequest{
			ProductId: i.ProductID,
			Quantity:  i.Quantity,
		}
		err := stream.Send(&order)
		if err != nil {
			http.Error(w, "Ma'lumotni jo'natishda xatolik...", http.StatusBadRequest)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		http.Error(w, "Xatolik...", http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Ma'lumotni encode qilishda xatolik...", http.StatusBadRequest)
	}

}

func (o *OrderClient) GetOrder(w http.ResponseWriter, r *http.Request) {
	var orderReq opb.GetOrderRequset

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "GET: bodydan o'qib olishda xatolik...", http.StatusBadRequest)
		return
	}

	protojson.Unmarshal(bytes, &orderReq)

	resp, err := o.Order.GetOrder(r.Context(), &orderReq)
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
