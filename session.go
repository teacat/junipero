package junipero

import (
	"time"

	"github.com/gorilla/websocket"
)

// Session 是單個客戶端階段。
type Session struct {
	// id 是此階段的獨立號碼。
	id int
	// store 是階段存儲資料。
	store map[string]interface{}
	// Subscriptions 是此階段訂閱的所有頻道。
	Subscriptions map[string]*Channel
	// isClosed 表示此階段是否已經關閉了。
	isClosed bool
	// conn 是該階段的 WebSocket 連線。
	conn *websocket.Conn

	// engine 是此階段所屬的引擎。
	engine *Engine
}

// NewSession 會在引擎中建立一個新的客戶端階段。
func (e *Engine) NewSession(conn *websocket.Conn) *Session {
	e.lastID++
	s := &Session{
		id:            e.lastID,
		store:         make(map[string]interface{}),
		Subscriptions: make(map[string]*Channel),
		conn:          conn,
		engine:        e,
	}
	e.sessions[e.lastID] = s
	return s
}

// Get 能夠從客戶端階段中取得暫存資料。
func (s *Session) Get(k string) (v interface{}, ok bool) {
	v, ok = s.store[k]
	return
}

// MustGet 能夠從客戶端階段中取得暫存資料，如果該資料不存在則呼叫 `panic`。
func (s *Session) MustGet(k string) interface{} {
	v, ok := s.Get(k)
	if !ok {
		panic(ErrKeyNotFound)
	}
	return v
}

// GetString 能夠從客戶端階段中取得 `string` 型態的暫存資料。
func (s *Session) GetString(k string) string {
	v, ok := s.Get(k)
	if !ok {
		return ""
	}
	return v.(string)
}

// GetBool 能夠從客戶端階段中取得 `bool` 型態的暫存資料。
func (s *Session) GetBool(k string) bool {
	v, ok := s.Get(k)
	if !ok {
		return false
	}
	return v.(bool)
}

// GetDuration 能夠從客戶端階段中取得 `time.Duration` 型態的暫存資料。
func (s *Session) GetDuration(k string) time.Duration {
	v, ok := s.Get(k)
	if !ok {
		return time.Duration(0)
	}
	return v.(time.Duration)
}

// GetFloat64 能夠從客戶端階段中取得 `float64` 型態的暫存資料。
func (s *Session) GetFloat64(k string) float64 {
	v, ok := s.Get(k)
	if !ok {
		return 0
	}
	return v.(float64)
}

// GetInt 能夠從客戶端階段中取得 `int` 型態的暫存資料。
func (s *Session) GetInt(k string) int {
	v, ok := s.Get(k)
	if !ok {
		return 0
	}
	return v.(int)
}

// GetInt64 能夠從客戶端階段中取得 `int64` 型態的暫存資料。
func (s *Session) GetInt64(k string) int64 {
	v, ok := s.Get(k)
	if !ok {
		return 0
	}
	return v.(int64)
}

// GetStringMap 能夠從客戶端階段中取得 `map[string]interface{}` 型態的暫存資料。
func (s *Session) GetStringMap(k string) map[string]interface{} {
	v, ok := s.Get(k)
	if !ok {
		return nil
	}
	return v.(map[string]interface{})
}

// GetStringMapString 能夠從客戶端階段中取得 `map[string]string` 型態的暫存資料。
func (s *Session) GetStringMapString(k string) map[string]string {
	v, ok := s.Get(k)
	if !ok {
		return nil
	}
	return v.(map[string]string)
}

// GetStringSlice 能夠從客戶端階段中取得 `[]string` 型態的暫存資料。
func (s *Session) GetStringSlice(k string) []string {
	v, ok := s.Get(k)
	if !ok {
		return nil
	}
	return v.([]string)
}

// GetTime 能夠從客戶端階段中取得 `time.Time` 型態的暫存資料。
func (s *Session) GetTime(k string) time.Time {
	v, ok := s.Get(k)
	if !ok {
		return time.Time{}
	}
	return v.(time.Time)
}

// Close 會良好地結束與此客戶端的連線。
func (s *Session) Close() error {
	if s.isClosed {
		return ErrSessionClosed
	}
	s.isClosed = true
	return s.conn.Close()
}

// CloseWithMsg 會關閉與此客戶端的連線，並在那之前傳送最後一則文字訊息。
func (s *Session) CloseWithMsg(msg string) error {
	err := s.Write(msg)
	if err != nil {
		return err
	}
	err = s.Close()
	if err != nil {
		return err
	}
	return nil
}

// CloseWithBinary 會關閉與此客戶端的連線，並在那之前傳送最後一則二進制訊息。
func (s *Session) CloseWithBinary(msg []byte) error {
	err := s.WriteBinary(msg)
	if err != nil {
		return err
	}
	err = s.Close()
	if err != nil {
		return err
	}
	return nil
}

// IsClosed 會表示此客戶端階段是否已經關閉連線了。
func (s *Session) IsClosed() bool {
	return s.isClosed
}

// Set 能夠將指定的資料存儲到此客戶端階段中作為暫存快取。
func (s *Session) Set(k string, v interface{}) {
	s.store[k] = v
}

// Delete 會將指定資料從暫存快取中移除。
func (s *Session) Delete(k string) error {
	_, ok := s.store[k]
	if !ok {
		return ErrKeyNotFound
	}
	delete(s.store, k)
	return nil
}

// Write 能透將文字訊息寫入到客戶端中。
func (s *Session) Write(msg string) error {
	err := s.conn.WriteMessage(int(TextMessage), []byte(msg))
	if err == nil {
		s.engine.handler.SentMessage(s, msg)
	}
	return err
}

// WriteBinary 能透將二進制訊息寫入到客戶端中。
func (s *Session) WriteBinary(msg []byte) error {
	err := s.conn.WriteMessage(int(BinaryMessage), []byte(msg))
	if err == nil {
		s.engine.handler.SentMessageBinary(s, msg)
	}
	return err
}

// Pong 能夠自主地回應客戶端一個 Pong 訊息，表示伺服器仍然有回應。
func (s *Session) Pong() error {
	return s.conn.WriteControl(int(PongMessage), []byte(``), time.Now().Add(s.engine.config.WriteWait))
}

// Ping 能夠詢問此客戶端的連線反應狀況，
// 如果在指定時間內沒有接收到 Pong 回應則會關閉並結束此連線。
func (s *Session) Ping() error {
	return s.conn.WriteControl(int(PingMessage), []byte(``), time.Now().Add(s.engine.config.WriteWait))
}

// Subscribe 會訂閱一個頻道。
func (s *Session) Subscribe(ch string) error {
	v, ok := s.engine.channels[ch]
	if !ok {
		return ErrChannelNotFound
	}
	_, ok = v.Sessions[s.id]
	if ok {
		return ErrChannelSubscribed
	}
	v.Sessions[s.id] = s
	s.Subscriptions[ch] = v
	return nil
}

// Unsubscribe 會取消訂閱一個頻道。
func (s *Session) Unsubscribe(ch string) error {
	v, ok := s.engine.channels[ch]
	if !ok {
		return ErrChannelNotFound
	}
	_, ok = v.Sessions[s.id]
	if !ok {
		return ErrChannelNotSubscribed
	}
	delete(v.Sessions, s.id)
	delete(s.Subscriptions, ch)
	return nil
}

// UnsubscribeAll 會取消訂閱此客戶端所有訂閱的頻道。
func (s *Session) UnsubscribeAll() {
	for k := range s.Subscriptions {
		s.Unsubscribe(k)
	}
}

// IsSubscribed 會表示客戶端是否有訂閱指定的頻道。
func (s *Session) IsSubscribed(ch string) bool {
	_, ok := s.Subscriptions[ch]
	return ok
}
