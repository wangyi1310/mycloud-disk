package session

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// NewMemoryStore 返回一个新的基于内存的 session store
func NewMemoryStore(keyPairs ...[]byte) *MemoryStore {
	return &MemoryStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &sessions.Options{
			Path:   "/",
			MaxAge: 86400 * 30, // 默认30天
		},
		sessions: make(map[string]map[string]*sessionData),
	}
}

type sessionData struct {
	values  map[interface{}]interface{}
	options *sessions.Options
	expires time.Time
}

type MemoryStore struct {
	Codecs  []securecookie.Codec
	Options *sessions.Options

	mu       sync.RWMutex
	sessions map[string]map[string]*sessionData // map[sessionID]map[name]*sessionData
}

// Get 实现 Store 接口的 Get 方法
func (s *MemoryStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

// New 实现 Store 接口的 New 方法
func (s *MemoryStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(s, name)
	session.Options = &(*s.Options) // 复制默认选项

	// 尝试从 cookie 中获取现有 session
	if cookie, err := r.Cookie(name); err == nil {
		if err = securecookie.DecodeMulti(name, cookie.Value, &session.ID, s.Codecs...); err == nil {
			s.mu.RLock()
			defer s.mu.RUnlock()

			if data, exists := s.sessions[session.ID]; exists {
				if sd, ok := data[name]; ok && sd.expires.After(time.Now()) {
					session.Values = sd.values
					session.Options = sd.options
					session.IsNew = false
					return session, nil
				}
			}
		}
	}

	// 创建新 session
	session.ID = string(securecookie.GenerateRandomKey(32))
	session.IsNew = true

	return session, nil
}

// Save 实现 Store 接口的 Save 方法
func (s *MemoryStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 如果 MaxAge < 0，则删除 session
	if session.Options.MaxAge < 0 {
		if data, exists := s.sessions[session.ID]; exists {
			delete(data, session.Name())
			if len(data) == 0 {
				delete(s.sessions, session.ID)
			}
		}

		// 设置删除 cookie
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
		return nil
	}

	// 编码 session ID
	encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, s.Codecs...)
	if err != nil {
		return err
	}

	// 存储 session 数据
	if _, exists := s.sessions[session.ID]; !exists {
		s.sessions[session.ID] = make(map[string]*sessionData)
	}

	expires := time.Now().Add(time.Second * time.Duration(session.Options.MaxAge))
	s.sessions[session.ID][session.Name()] = &sessionData{
		values:  session.Values,
		options: session.Options,
		expires: expires,
	}

	// 设置 cookie
	http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))

	// 定期清理过期 session
	go s.cleanupExpiredSessions()

	return nil
}

// cleanupExpiredSessions 清理过期的 session
func (s *MemoryStore) cleanupExpiredSessions() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, data := range s.sessions {
		for name, sd := range data {
			if sd.expires.Before(time.Now()) {
				delete(data, name)
			}
		}
		if len(data) == 0 {
			delete(s.sessions, id)
		}
	}
}
