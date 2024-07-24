package service

import (
	"context"
	"database/sql"
	"log"
	"userservice/proto"
)

type UserServer struct {
	proto.UnimplementedUserServiceServer
	db *sql.DB
	LG *log.Logger
}

func NewUserServer(DB *sql.DB, lg *log.Logger) *UserServer {
	return &UserServer{db: DB, LG: lg}
}

func (u *UserServer) CreateUser(ctx context.Context, req *proto.UserRequest) (*proto.UserResponse, error) {
	u.LG.Println("INFO: CreateUser servicega so'rov kelib tushdi...")
	query := `INSERT INTO users (firstname, lastname, age, email) VALUES ($1, $2, $3, $4) RETURNING id, registered_at`

	var id int32
	var time string
	err := u.db.QueryRow(query, req.Firstname, req.Lastname, req.Age, req.Email).Scan(&id, &time)
	if err != nil {
		u.LG.Println("ERROR: Tablega ma'lumot qo'shishda xatolik...")
		log.Fatalf("Tablega ma'lumot qo'shishda xatolik - %v", err)
	}
	resp := proto.UserResponse{Id: id, Firstname: req.Firstname, Lastname: req.Lastname, Age: req.Age, Email: req.Email, RegisteredAt: time}
	u.LG.Println("INFO: Create User service muvaffaqiyatli javob qaytardi...")
	return &resp, nil
}

func (u *UserServer) GetUserID(ctx context.Context, req *proto.GetUserId) (*proto.UserResponse, error) {
	query := `SELECT id, firstname, lastname, age, email, registered_at FROM users WHERE id = $1`
	var resp proto.UserResponse

	err := u.db.QueryRow(query, req.Id).Scan(&resp.Id, &resp.Firstname, &resp.Lastname, &resp.Age, &resp.Email, &resp.RegisteredAt)
	if err != nil {
		u.LG.Println("ERROR: Tabledan ma'lumot o'qishda xatolik...")
		log.Fatalf("Tabledan ma'lumotlarni o'qishda xatolik - %v", err)
	}

	return &resp, nil
}
