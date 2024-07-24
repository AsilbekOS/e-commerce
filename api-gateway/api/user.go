package api

import (
	upb "api-gateway/protos/user-proto"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
)

type UserClient struct {
	Client upb.UserServiceClient
}

func NewUserClient(cl upb.UserServiceClient) *UserClient {
	return &UserClient{Client: cl}
}

func (u *UserClient) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userReq upb.UserRequest

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "POST: bodydan o'qib olishda xatolik...", http.StatusBadRequest)
		return
	}
	protojson.Unmarshal(bytes, &userReq)

	resp, err := u.Client.CreateUser(r.Context(), &userReq)
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

func (u *UserClient) GetUserID(w http.ResponseWriter, r *http.Request) {
	var userReq upb.GetUserId

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "GET: bodydan o'qib olishda xatolik...", http.StatusBadRequest)
		return
	}
	protojson.Unmarshal(bytes, &userReq)

	resp, err := u.Client.GetUserID(r.Context(), &userReq)
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
