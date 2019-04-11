package junipero

// Channel 呈現了一個頻道。
type Channel struct {
	// name 是這個頻道的名稱。
	name string
	// sessions 是所有訂閱此頻道的階段客戶端連線。
	Sessions map[int]*Session
	// isClosed 表示此頻道是否已經關閉了。
	isClosed bool
	// config 是頻道設置。
	config *ChannelConfig
}

// ChannelConfig 是頻道設置。
type ChannelConfig struct {
}

// NewChannel 會建立一個新的可訂閱頻道。
func (e *Engine) NewChannel(name string, conf *ChannelConfig) *Channel {
	ch := &Channel{
		name:   name,
		config: conf,
	}
	e.channels[name] = ch
	return ch
}

// Broadcast 能夠將文字訊息廣播給頻道中的所有客戶端。
func (c *Channel) Broadcast(msg string) error {
	if c.isClosed {
		return ErrChannelClosed
	}
	for _, v := range c.Sessions {
		v.Write(msg)
	}
	return nil
}

// BroadcastFilter 能夠將文字訊息廣播給頻道中被篩選客戶端。
func (c *Channel) BroadcastFilter(msg string, fn func(*Session) bool) error {
	if c.isClosed {
		return ErrChannelClosed
	}
	for _, v := range c.Sessions {
		if fn(v) {
			v.Write(msg)
		}
	}
	return nil
}

// BroadcastOthers 能夠將文字訊息廣播給頻道中指定以外的所有客戶端。
func (c *Channel) BroadcastOthers(msg string, s *Session) error {
	if c.isClosed {
		return ErrChannelClosed
	}
	for _, v := range c.Sessions {
		if v != s {
			v.Write(msg)
		}
	}
	return nil
}

// BroadcastBinary 能夠將二進制訊息廣播給頻道中的所有客戶端。
func (c *Channel) BroadcastBinary(msg []byte) error {
	if c.isClosed {
		return ErrChannelClosed
	}
	for _, v := range c.Sessions {
		v.WriteBinary(msg)
	}
	return nil
}

// BroadcastBinaryFilter 能夠將二進制訊息廣播給頻道中被篩選客戶端。
func (c *Channel) BroadcastBinaryFilter(msg []byte, fn func(*Session) bool) error {
	if c.isClosed {
		return ErrChannelClosed
	}
	for _, v := range c.Sessions {
		if fn(v) {
			v.WriteBinary(msg)
		}
	}
	return nil
}

// BroadcastBinaryOthers 能夠將二進制訊息廣播給頻道中指定以外的所有客戶端。
func (c *Channel) BroadcastBinaryOthers(msg []byte, s *Session) error {
	if c.isClosed {
		return ErrChannelClosed
	}
	for _, v := range c.Sessions {
		if v != s {
			v.WriteBinary(msg)
		}
	}
	return nil
}

// Close 會關閉頻道並將其訂閱的客戶端全部取消訂閱。
func (c *Channel) Close() error {
	if c.isClosed {
		return ErrChannelClosed
	}
	c.isClosed = true
	for _, v := range c.Sessions {
		v.Unsubscribe(c.name)
	}
	return nil
}

// CloseWithMsg 會關閉頻道並取消所有客戶端訂閱，但在那之前會先發送一則文字訊息。
func (c *Channel) CloseWithMsg(msg string) error {
	if c.isClosed {
		return ErrChannelClosed
	}
	c.isClosed = true
	for _, v := range c.Sessions {
		v.Unsubscribe(c.name)
		v.Write(msg)
	}
	return nil
}

// CloseWithBinary 會關閉頻道並取消所有客戶端訂閱，但在那之前會先發送一則二進制訊息。
func (c *Channel) CloseWithBinary(msg []byte) error {
	if c.isClosed {
		return ErrChannelClosed
	}
	c.isClosed = true
	for _, v := range c.Sessions {
		v.Unsubscribe(c.name)
		v.WriteBinary(msg)
	}
	return nil
}

// IsClosed 會回傳表示這個頻道是否已經關閉。
func (c *Channel) IsClosed() bool {
	return c.isClosed
}

// Contains 會表示指定的客戶端是否有訂閱此頻道。
func (c *Channel) Contains(s *Session) bool {
	_, ok := c.Sessions[s.id]
	return ok
}

// Len 會表示頻道的總訂閱客戶端數量。
func (c *Channel) Len() int {
	return len(c.Sessions)
}
