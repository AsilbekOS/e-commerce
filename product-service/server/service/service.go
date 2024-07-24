package service

import (
	"context"
	"database/sql"
	"log"
	"productservice/proto"
	"time"
)

type ProductServer struct {
	proto.UnimplementedProductServiceServer
	db *sql.DB
	lg *log.Logger
}

func NewProductServer(DB *sql.DB, LG *log.Logger) *ProductServer {
	return &ProductServer{db: DB, lg: LG}
}

func (p *ProductServer) CreateProduct(ctx context.Context, req *proto.ProductRequest) (*proto.ProductResponse, error) {
	p.lg.Println("INFO: CreateProduct servicega ma'lumot kelib tushdi...")
	query := `INSERT INTO products (name, quantity, price) VALUES ($1, $2, $3) RETURNING id, created_at`

	var id int32
	var cd_at string

	err := p.db.QueryRow(query, req.Name, req.Quantity, req.Price).Scan(&id, &cd_at)
	if err != nil {
		p.lg.Println("ERROR: Tablega ma'lumot qo'shishda xatolik...")
		log.Fatalf("Tablega ma'lumot qo'shishda xatolik - %v", err)
	}

	resp := proto.ProductResponse{Id: id, Name: req.Name, Quantity: req.Quantity, Price: req.Price, CreatedAt: cd_at}
	p.lg.Println("INFO: CreateUser service muvaffaqiyatli javob qaytardi...")
	log.Println("Successfully created table...")
	return &resp, nil
}

func (p *ProductServer) GetProductID(ctx context.Context, req *proto.GetProductRequest) (*proto.ProductResponse, error) {
	p.lg.Println("INFO: GetProductID servicega ma'lumot kelib tushdi...")	
	var cost bool
	query := `SELECT EXISTS(SELECT * FROM products WHERE id = $1)`
	var resp proto.ProductResponse

	err := p.db.QueryRow(query, req.Id).Scan(&cost)
	if err != nil {
		log.Println("Exacda xatolik...", err)
		return nil, err
	}
	if !cost {
		log.Fatalf("Bunday ID mavjud emas. - %v", err)
	} else {
		// if r, _ := n.RowsAffected(); r == 0 {
		// 	log.Println("Bunday IDli product mavjud emas...")
		// 	return nil, fmt.Errorf("bunday IDli product mavjud emas")
		// }

		err = p.db.QueryRow("SELECT * FROM products WHERE id = $1", req.Id).Scan(&resp.Id, &resp.Name, &resp.Quantity, &resp.Price, &resp.CreatedAt)
		if err != nil {
			p.lg.Println("ERROR: Tabledan ma'lumot o'qishda xatolik...")
			log.Fatalf("ERROR: Tabledan ma'lumot o'qishda xatolik: %v", err)
		}

		p.lg.Println("INFO: GetProductID service muvaffaqiyatli javob qaytardi...")
		return &resp, nil
	}
	return nil,nil
}

func (p *ProductServer) ListProducts(req *proto.Empty, stream proto.ProductService_ListProductsServer) error {
	query := `SELECT id, name, quantity, price, created_at FROM products`

	rows, err := p.db.Query(query)
	if err != nil {
		p.lg.Println("ERROR: Tabledan barcha ma'lumotlarni o'qishda xatolik...")
		log.Fatalf("Tabledan barcha ma'lumotlarni o'qishda xatolik: %v", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var product proto.ProductResponse
		err := rows.Scan(&product.Id, &product.Name, &product.Quantity, &product.Price, &product.CreatedAt)
		if err != nil {
			p.lg.Println("ERROR: Qatorlarni skanerlashda xatolik...")
			log.Fatalf("Qatorlarni skanerlashda xatolik: %v", err)
			return err
		}

		err = stream.Send(&product)
		if err != nil {
			p.lg.Println("ERROR: Ma'lumotlarni stream orqali yuborishda xatolik...")
			log.Fatalf("Ma'lumotlarni stream orqali yuborishda xatolik: %v", err)
			return err
		}

		time.Sleep(time.Second)
	}

	if err = rows.Err(); err != nil {
		p.lg.Println("ERROR: Qatorlar bilan ishlashda xatolik...")
		log.Fatalf("Qatorlar bilan ishlashda xatolik: %v", err)
		return err
	}

	return nil
}
