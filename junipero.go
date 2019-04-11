package junipero

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Engine 是 WebSocket 引擎。
type Engine struct {
	handler  Handler
	sessions map[int]*Session
	channels map[string]*Channel
	config   *EngineConfig
	isClosed bool
	lastID   int
}

// EngineConfig 是引擎選項設置。
type EngineConfig struct {
	// WriteWait 是到逾時之前的等待時間。
	WriteWait time.Duration
	// PongWait 是等待 Pong 回應的時間，在指定時間內客戶端如果沒有任何響應，該 WebSocket 連線則會被迫中止。
	// 設置為 `0` 來停用無響應自動斷線的功能。
	PongWait time.Duration
	// PingPeriod 是 Ping 的週期時間。
	PingPeriod time.Duration
	// MaxMessageSize 是最大可接收的訊息位元組大小，
	// 溢出此大小的訊息會被拋棄。
	MaxMessageSize int64
	// Upgrader 是 WebSocket 升級的相關設置。
	Upgrader *websocket.Upgrader
}

// Handler 是 WebSocket 訊息和相關功能的處理函式。
type Handler interface {
	Close(*Session, CloseStatus, string) error
	Connect(*Session)
	Disconnect(*Session)
	Error(*Session, error)
	Message(*Session, string)
	MessageBinary(*Session, []byte)
	SentMessage(*Session, string)
	SentMessageBinary(*Session, []byte)
	Ping(*Session)
	Pong(*Session)
	Request(http.ResponseWriter, *http.Request, *Session)
}

// NewServer 會建立一個新的 WebSocket 伺服器。
func NewServer(conf *EngineConfig, handler Handler) *Engine {
	return &Engine{
		handler:  handler,
		config:   conf,
		sessions: make(map[int]*Session),
		channels: make(map[string]*Channel),
	}
}

// DefaultConfig 會回傳一個新的預設引擎設置。
func DefaultConfig() *EngineConfig {
	return &EngineConfig{
		WriteWait:      30,
		PongWait:       10,
		PingPeriod:     20,
		MaxMessageSize: 10 * 1024 * 1024,
		Upgrader: &websocket.Upgrader{
			HandshakeTimeout: 30 * time.Second,
			ReadBufferSize:   1024,
			WriteBufferSize:  1024,
		},
	}
}

// HandlerFunc 是用以傳入 HTTP 伺服器協助升級與接收 WebSocket 相關資訊的最重要函式。
func (e *Engine) HandlerFunc() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if e.isClosed {
			panic(ErrEngineClosed)
		}
		c, err := e.config.Upgrader.Upgrade(w, r, nil)
		s := e.NewSession(c)
		if err != nil {
			e.handler.Error(s, err)
			return
		}
		e.handler.Request(w, r, s)

		c.SetPingHandler(func(m string) error {
			e.handler.Ping(s)
			s.Pong()
			return nil
		})
		c.SetPongHandler(func(m string) error {
			e.handler.Pong(s)
			return nil
		})
		c.SetCloseHandler(func(code int, msg string) error {
			s.Close()
			e.handler.Close(s, CloseStatus(code), msg)
			if CloseStatus(code) == CloseNormalClosure {
				e.handler.Disconnect(s)
			}
			return nil
		})

		e.handler.Connect(s)

		defer func() {
			s.Close()
		}()

		for {
			typ, msg, err := c.ReadMessage()
			if err != nil {
				if !s.isClosed {
					s.Close()
					e.handler.Error(s, err)
				}
				break
			}
			switch MessageType(typ) {
			case TextMessage:
				e.handler.Message(s, string(msg))
				break
			case BinaryMessage:
				e.handler.MessageBinary(s, msg)
				break
			}
		}
	}
}

// Broadcast 會將文字訊息傳送到所有連線的客戶端。
func (e *Engine) Broadcast(msg string) {
	for _, v := range e.sessions {
		v.Write(msg)
	}
}

// BroadcastFilter 會將文字訊息傳送到經篩選的客戶端。
func (e *Engine) BroadcastFilter(msg string, fn func(*Session) bool) {
	for _, v := range e.sessions {
		if fn(v) {
			v.Write(msg)
		}
	}
}

// BroadcastOthers 會將文字訊息傳送到指定客戶端以外的所有客戶端。
func (e *Engine) BroadcastOthers(msg string, s *Session) {
	for _, v := range e.sessions {
		if v != s {
			v.Write(msg)
		}
	}
}

// BroadcastMultiple 會將文字訊息傳送到指定客戶端的客戶端們。
func (e *Engine) BroadcastMultiple(msg string, sessions []*Session) {
	for _, v := range sessions {
		v.Write(msg)
	}
}

// BroadcastBinary 會將二進制訊息傳送到所有連線的客戶端。
func (e *Engine) BroadcastBinary(msg []byte) {
	for _, v := range e.sessions {
		v.WriteBinary(msg)
	}
}

// BroadcastBinaryFilter 會將二進制訊息傳送到經篩選的客戶端。
func (e *Engine) BroadcastBinaryFilter(msg []byte, fn func(*Session) bool) {
	for _, v := range e.sessions {
		if fn(v) {
			v.WriteBinary(msg)
		}
	}
}

// BroadcastBinaryOthers 會將二進制訊息傳送到指定客戶端以外的所有客戶端。
func (e *Engine) BroadcastBinaryOthers(msg []byte, s *Session) {
	for _, v := range e.sessions {
		if v != s {
			v.WriteBinary(msg)
		}
	}
}

// BroadcastBinaryMultiple 會將二進制訊息傳送到指定客戶端的客戶端們。
func (e *Engine) BroadcastBinaryMultiple(msg []byte, sessions []*Session) {
	for _, v := range sessions {
		v.WriteBinary(msg)
	}
}

// Close 會關閉整個引擎並中斷所有連線。
func (e *Engine) Close() {
	for _, v := range e.sessions {
		v.Close()
	}
	e.isClosed = true
}

// CloseWithMsg 會關閉引擎並在那之前傳送最後一則文字訊息。
func (e *Engine) CloseWithMsg(msg string) {
	for _, v := range e.sessions {
		v.CloseWithMsg(msg)
	}
	e.isClosed = true
}

// CloseWithBinary 會關閉引擎並在那之前傳送最後一則二進制訊息。
func (e *Engine) CloseWithBinary(msg []byte) {
	for _, v := range e.sessions {
		v.CloseWithBinary(msg)
	}
	e.isClosed = true
}

// IsClosed 會表示該引擎是否已經關閉了。
func (e *Engine) IsClosed() bool {
	return e.isClosed
}

// Len 會取得正在連線的客戶端總數。
func (e *Engine) Len() int {
	return len(e.sessions)
}
