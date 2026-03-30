package inmemoryrepo

import (
	"context"
	"fmt"
	"log"
	"parsaraki-go/internal/infrastructure/telegram"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	connection *pgxpool.Pool
}

func NewPostgresRepo(ctx context.Context) *PostgresRepo {
	dsn := "postgres://admin:admin@localhost:5432/parsarakidb"
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
        CREATE TABLE IF NOT EXISTS telegramrequests (
            id SERIAL PRIMARY KEY,
			message_id INT NOT NULL,
			date TEXT NOT NULL,
            from_username TEXT NOT NULL,
            chat_id INT NOT NULL,
			text TEXT NOT NULL
        );
    `
	_, err = pool.Exec(context.Background(), createTable)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		<-ctx.Done()
		pool.Close()
		log.Println("Postgres DB connection closed")
	}()
	return &PostgresRepo{connection: pool}
}

func (pr *PostgresRepo) SaveBotRequest(botMsg telegram.BotMessage) error {
	saveRequest := `
		INSERT INTO telegramrequests (message_id, date, from_username, chat_id, text)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`
	var newID int
	err := pr.connection.QueryRow(
		context.Background(),
		saveRequest,
		botMsg.Message.MessageID,     // message_id
		botMsg.Message.Date,          // date
		botMsg.Message.From.Username, // from_username
		botMsg.Message.Chat.ID,       // chat_id
		botMsg.Message.Text,          // text
	).Scan(&newID)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Message saved to telegramrequests table: id=%v", newID)
	return nil
}
