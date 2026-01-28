package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/ArtoIi/BIBLIOTECA/internal/book"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func main() {
	// 1. Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}

	ctx := context.Background()

	// 2. Configurar o Firestore
	sa := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	client, err := firestore.NewClient(ctx, "vigia-verde", sa)
	if err != nil {
		log.Fatalf("Erro ao iniciar Firestore: %v", err)
	}
	defer client.Close()

	// 3. A "Fiação" (Injeção de Dependência)
	// Aqui conectamos as camadas que você criou
	repo := book.NewRepository(client) // Passa o cliente do banco para o repo
	_, err = client.Collections(ctx).Next()
	if err != nil {
		log.Fatalf("Falha real de conexão com Firestore: %v", err)
	}
	log.Println("Conexão com Firestore estabelecida com sucesso!")
	service := book.NewService(repo)    // Passa o repo para o service
	handler := book.NewHandler(service) // Passa o service para o handler

	// 4. Rotas
	mux := http.NewServeMux()
	handler.Routes(mux)

	// 5. Start Server
	log.Println("Servidor rodando na porta 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
