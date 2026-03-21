package agentcore

import (
	"errors"
	"sync"
	"time"
)

type Session struct {
	SessionID   string
	UserID      string
	Type        string
	Status      string
	GroupInfo   *GroupInfo
	LastAgentID string
	CreatedAt   int64
	UpdatedAt   int64
	ExpiresAt   int64
}

type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
	}
}

func (sm *SessionManager) GetOrCreateSession(userID string, sessionID string, sessionType string) (*Session, error) {
	now := time.Now().Unix()

	if sessionID != "" {
		sm.mu.RLock()
		s, ok := sm.sessions[sessionID]
		sm.mu.RUnlock()
		if ok {
			if s.Type == "group" || sessionType == "group" {
				sm.mu.Lock()
				s.Type = "group"
				if s.GroupInfo == nil {
					s.GroupInfo = &GroupInfo{
						GroupID:   s.SessionID,
						GroupName: s.SessionID,
						CreatorID: s.UserID,
					}
				}
				s.UpdatedAt = now
				sm.mu.Unlock()
				return s, nil
			}

			if s.UserID != userID {
				return nil, errors.New("session user mismatch")
			}
			sm.mu.Lock()
			s.UpdatedAt = now
			sm.mu.Unlock()
			return s, nil
		}
	}

	if sessionType == "" {
		sessionType = "single"
	}
	if sessionType != "single" && sessionType != "group" {
		sessionType = "single"
	}

	id := sessionID
	if id == "" {
		id = generateID()
	}

	s := &Session{
		SessionID: id,
		UserID:    userID,
		Type:      sessionType,
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now + 86400,
	}
	if sessionType == "group" {
		s.GroupInfo = &GroupInfo{
			GroupID:   id,
			GroupName: id,
			CreatorID: userID,
		}
	}

	sm.mu.Lock()
	sm.sessions[id] = s
	sm.mu.Unlock()
	return s, nil
}

func (sm *SessionManager) UpdateSession(session *Session) error {
	if session == nil || session.SessionID == "" {
		return errors.New("invalid session")
	}
	now := time.Now().Unix()
	sm.mu.Lock()
	defer sm.mu.Unlock()
	session.UpdatedAt = now
	if session.ExpiresAt == 0 {
		session.ExpiresAt = now + 86400
	}
	sm.sessions[session.SessionID] = session
	return nil
}
