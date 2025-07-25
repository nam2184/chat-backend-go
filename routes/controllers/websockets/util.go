package websockets

import (
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/nam2184/mymy/models/db"
)

func (hub *ChatHub) closeConn(conn *websocket.Conn) {
	err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		hub.log.Debug().Msgf("Error sending close message: %s", err.Error())
	}
	err = conn.Close()
	if err != nil {
		hub.log.Debug().Msgf("Error closing websocket connection: %s", err.Error())
	}
}

func writeMessageDB(tx *sqlx.Tx, message db.Message) error {
	query := `INSERT INTO messages (chat_id, sender_id, sender_name, receiver_id, content, image, type, is_typing, timestamp)
        VALUES (:chat_id, :sender_id, :sender_name, :receiver_id, :content, :image, :type, :is_typing, :timestamp)
			  RETURNING id;`

	_, err := tx.NamedExec(query, message)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func writeEncryptedMessageDB(tx *sqlx.Tx, message db.EncryptedMessage) error {
	query := `INSERT INTO encrypted_messages (chat_id, sender_id, sender_name, receiver_id, content, image, type, is_typing, timestamp)
        VALUES (:chat_id, :sender_id, :sender_name, :receiver_id, :content, :image, :type, :is_typing, :timestamp)
			  RETURNING id;`

	_, err := tx.NamedExec(query, message)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
