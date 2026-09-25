package live

import (
	"context"
	"net/http"
	"time"

	"github.com/coder/websocket"
)

const writeTimeout = 5 * time.Second

// Register adds the live routes to mux. The route pattern and PathValue
// live next to each other, so their parameter names cannot drift apart.
func Register(mux *http.ServeMux, hub *Hub, allowedOrigins []string) {
	mux.HandleFunc("GET /ws/matches/{matchId}/stream", streamHandler(hub, allowedOrigins))
}

func streamHandler(hub *Hub, allowedOrigins []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		matchID := r.PathValue("matchId")

		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{OriginPatterns: allowedOrigins})
		if err != nil {
			return // Accept has already written the HTTP error
		}
		defer conn.CloseNow()

		client := hub.Subscribe(matchID)
		defer hub.Unsubscribe(matchID, client)

		// We never read from the client, but the library must still process
		// its control frames. ctx is cancelled when the client disconnects.
		ctx := conn.CloseRead(r.Context())

		for {
			select {
			case msg, ok := <-client.Messages():
				if !ok {
					conn.Close(websocket.StatusTryAgainLater, "dropped")
					return
				}
				if err := write(ctx, conn, msg); err != nil {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}
}

func write(ctx context.Context, conn *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, writeTimeout)
	defer cancel()
	return conn.Write(ctx, websocket.MessageText, msg)
}
