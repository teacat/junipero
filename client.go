package junipero

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Client 呈現了一個 WebSocket 客戶端。
type Client struct {
	// Config 是客戶端設置。
	config *ClientConfig
	// conn 是底層的 WebSocket 連線。
	conn *websocket.Conn
	// isClosed 會表示此客戶端是否已經關閉連線了。
	isClosed bool
}

// ClientConfig 是客戶端設置。
type ClientConfig struct {
	// Address 是遠端伺服器位置（如：`ws://127.0.0.1:1234/echo`）。
	Address string
	// header 是 WebSocket 初次發送時順帶傳輸的 HTTP 標頭資訊。
	Header http.Header
	// WriteWait 是每次訊息寫入時的逾時時間。
	WriteWait time.Duration
}

// NewClient 會建立客戶端並連線到指定的 WebSocket 伺服端。
func NewClient(conf *ClientConfig) (*Client, *http.Response, error) {
	if conf.WriteWait == 0 {
		conf.WriteWait = time.Second * 30
	}
	conn, resp, err := websocket.DefaultDialer.Dial(conf.Address, conf.Header)
	if err != nil {
		return nil, resp, err
	}
	client := &Client{
		config: conf,
		conn:   conn,
	}
	return client, resp, nil
}

// ReadMessage 會阻塞程式直到有訊息為止，接收到的訊息會 `string` 字串標準訊息。
// 任何系統訊息像是 Ping-Pong 與 Close 都不會出現在這裡。
func (c *Client) Read() (string, error) {
	if c.isClosed {
		return "", ErrConnectionClosed
	}
	for {
		typ, msg, err := c.ReadAll()
		if err != nil {
			return "", err
		}
		if typ != TextMessage {
			continue
		}
		return string(msg), nil
	}
}

// ReadBinary 會阻塞程式直到有訊息為止，接收到的訊息會是 `[]byte` 二進制標準訊息。
// 任何系統訊息像是 Ping-Pong 與 Close 都不會出現在這裡。
func (c *Client) ReadBinary() ([]byte, error) {
	if c.isClosed {
		return []byte(``), ErrConnectionClosed
	}
	for {
		typ, msg, err := c.ReadAll()
		if err != nil {
			return []byte(``), err
		}
		if typ != BinaryMessage {
			continue
		}
		return msg, nil
	}
}

// ReadAll 會阻塞程式直到有訊息為止，
// 這會接收到所有訊息像是 Ping-Pong 與 Close 或標準的文字甚至二進制訊息。
func (c *Client) ReadAll() (MessageType, []byte, error) {
	if c.isClosed {
		return 0, []byte(``), ErrConnectionClosed
	}
	typ, msg, err := c.conn.ReadMessage()
	if err != nil {
		return MessageType(typ), msg, err
	}
	return MessageType(typ), msg, nil
}

// Disconnect 會依照正常手續告訴伺服器關閉並結束客戶端連線。
func (c *Client) Disconnect() error {
	if c.isClosed {
		return ErrConnectionClosed
	}
	c.isClosed = true
	return c.conn.WriteControl(int(CloseMessage), websocket.FormatCloseMessage(int(CloseNormalClosure), ""), time.Now().Add(c.config.WriteWait))
}

// DisconnectWithMsg 會依照正常手續且帶有文字訊息告訴伺服器關閉並結束客戶端連線。
func (c *Client) DisconnectWithMsg(msg string) error {
	if c.isClosed {
		return ErrConnectionClosed
	}
	c.isClosed = true
	return c.conn.WriteControl(int(CloseMessage), websocket.FormatCloseMessage(int(CloseNormalClosure), msg), time.Now().Add(c.config.WriteWait))
}

// Close 會關閉並結束客戶端連線。
func (c *Client) Close() error {
	if c.isClosed {
		return ErrConnectionClosed
	}
	c.isClosed = true
	return c.conn.Close()
}

// Write 能夠傳送文字訊息至伺服端。
func (c *Client) Write(msg string) error {
	if c.isClosed {
		return ErrConnectionClosed
	}
	return c.conn.WriteMessage(int(TextMessage), []byte(msg))
}

// WriteBinary 能夠傳送二進制訊息至伺服端。
func (c *Client) WriteBinary(msg []byte) error {
	if c.isClosed {
		return ErrConnectionClosed
	}
	return c.conn.WriteMessage(int(BinaryMessage), msg)
}

// Ping 能夠發送 Ping 至伺服端並且等待 Pong 回應，
func (c *Client) Ping() error {
	if c.isClosed {
		return ErrConnectionClosed
	}
	return c.conn.WriteControl(int(PingMessage), []byte(``), time.Now().Add(c.config.WriteWait))
}

// Pong 能夠主動不等待 Ping 的情況下直接回應伺服端。
func (c *Client) Pong() error {
	if c.isClosed {
		return ErrConnectionClosed
	}
	return c.conn.WriteControl(int(PongMessage), []byte(``), time.Now().Add(c.config.WriteWait))
}

// IsClosed 會表示該連線是否已經關閉並結束了。
func (c *Client) IsClosed() bool {
	return c.isClosed
}
