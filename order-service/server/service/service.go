package service

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"orderservice/internal"
	"orderservice/proto"
	"time"
)

type OrderServer struct {
	proto.UnimplementedOrderServiceServer
	db *sql.DB
	lg *log.Logger
}

func NewOrderServer(DB *sql.DB, LG *log.Logger) *OrderServer {
	return &OrderServer{db: DB, lg: LG}
}

func (o *OrderServer) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderResponse, error) {
	o.lg.Println("CreateOrder servicega ma'lumot kelib tushdi...")

	// var prc proto.OrderResponse

	product := internal.Conn(req.ProductId)

	TotalAmount := product.Price * float64(req.Quantity)

	query := `INSERT INTO orders (product_id, product_name, quantity, price, total_amount) VALUES ($1, $2, $3, $4, $5) RETURNING product_id, created_at`

	var orderId int
	var createdAt time.Time

	err := o.db.QueryRow(query, req.ProductId, product.Name, req.Quantity, product.Price, TotalAmount).Scan(&orderId, &createdAt)
	if err != nil {
		o.lg.Println("ERROR: Buyurtmani yaratishda xatolik...")
		log.Fatalf("Buyurtmani yaratishda xatolik - %v", err)
	}

	resp := proto.OrderResponse{
		Id:          int32(orderId),
		ProductName: product.Name,
		Quantity:    req.Quantity,
		Price:       product.Price,
		TotalAmount: TotalAmount,
		CreatedAt:   createdAt.Format(time.RFC3339)}

	o.lg.Println("INFO: Buyurtma muvaffaqiyatli yaratildi...")
	log.Println("Successfully created order...")
	return &resp, nil
}

func (o *OrderServer) CreateOrders(stream proto.OrderService_CreateOrdersServer) error {
	o.lg.Println("CreateOrders servicega ma'lumot kelib tushdi...")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&proto.MessageOrderResponse{Message: "Successfully created orders"})
		}
		if err != nil {
			o.lg.Println("ERROR: Streamdan ma'lumot olishda xatolik...", err)
			return err
		}

		product := internal.Conn(req.ProductId)

		totalAmount := product.Price * float64(req.Quantity)

		// var orderId int32
		// var createdAt time.Time
		query := `INSERT INTO orders (product_id, product_name, price, total_amount, quantity) VALUES ($1, $2, $3, $4, $5)`
		_, err = o.db.Exec(query, req.ProductId, product.Name, product.Price, totalAmount, req.Quantity)
		if err != nil {
			o.lg.Println("ERROR: Buyurtmani yaratishda xatolik...", err)
			return fmt.Errorf("buyurtmani yaratishda xatolik - %v", err)
		}
	}
}

func (o *OrderServer) GetOrder(ctx context.Context, req *proto.GetOrderRequset) (*proto.OrderResponse, error) {
	o.lg.Println("GetOrder servicega ma'lumot kelib tushdi...")

	var resp proto.OrderResponse

	query := `SELECT product_id, product_name, quantity, price, total_amount, created_at FROM orders WHERE product_id = $1`

	err := o.db.QueryRow(query, req.Id).Scan(&resp.Id, &resp.ProductName, &resp.Quantity, &resp.Price, &resp.TotalAmount, &resp.CreatedAt)
	if err != nil {
		o.lg.Println("ERROR: Tabledan ma'lumot o'qishda xatolik...")
		log.Fatalf("Tabledan ma'lumotlarni o'qishda xatolik - %v", err)
	}
	return &resp, nil
}
