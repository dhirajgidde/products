package main 

import(
	"database/sql"
	"fmt"
	"log"
	"net"
	// "time"
	"context"
	// "net/http"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"

	messager "grpc_ex.com/v1/productsProto"
	
)


// products struct
type Products struct {
	Title     			string `json:"title"`
	SKU       			string `json:"sku"`
	AccountCode       	string `json:"accountCode"`
}	

//server struct 
type server struct {
	db      *sql.DB
	messager .UnimplementedMessageReceiverServer
}

func main() {
	fmt.Println("Product Server is started")
	var err error
	dbURL := fmt.Sprint("root:dev@tcp(127.0.0.1:3306)/dummy_product_service")
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()	
	grpcServer := grpc.NewServer()
	s := &server{
		db: db,
	}
	err = s.dbPrepare()
	if err != nil {
		log.Println(err)
		return
	}
	messager.RegisterMessageReceiverServer(grpcServer, s)
	
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalln("Could not listen on port 9000:", err)
		return
	}
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln(err)
		return
	}
	
}


func (s *server) dbPrepare() error {
	var err error
	_, err = s.db.Exec(`
    CREATE TABLE IF NOT EXISTS
    products (
        id int not null auto_increment,
        title tinytext not null,
        sku tinytext not null,
        accountCode tinytext not null,
        primary key (id)
    )
    `)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}


func (s *server) EnableProducts(ctx context.Context,product *messager.Product) (*messager.ProdResponse, error) {
	prod := &Products{product.Title, product.SKU, product.AccountCode}
	err := s.storageProductInsert(prod)
	if err != nil {
		log.Println(err)
		return &messager.ProdResponse{ResponseMessage : "Something went wrong"}, err
	}	
	return &messager.ProdResponse{ResponseMessage : "Product Created"}, nil
}

func (s *server) storageProductInsert(p *Products) error {
	tx, err := s.db.Begin()
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return err
	}
	_, err = tx.Exec(`
    INSERT INTO products (
        title,
        sku,
		accountCode
    )
    VALUES (?, ?, ?)
    `,
		p.Title, p.SKU, p.AccountCode)
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return err
	}
	tx.Commit()
	return nil
}