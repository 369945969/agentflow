package agentcore

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

type LogEntry struct {
	Timestamp string
	Type      string
	Content   string
}

type LogPersistence struct {
	basePath    string
	fileHandles map[string]*os.File
	mu          sync.RWMutex
}

func NewLogPersistence() *LogPersistence {
	return &LogPersistence{
		basePath:    "logs",
		fileHandles: make(map[string]*os.File),
	}
}

func (lp *LogPersistence) SetBasePath(basePath string) {
	basePath = strings.TrimSpace(basePath)
	if basePath == "" {
		return
	}
	lp.mu.Lock()
	lp.basePath = basePath
	lp.mu.Unlock()
}

func (lp *LogPersistence) GetLogFilePath(userID string, sessionID string) string {
	safeUser := sanitizePathSegment(userID)
	safeSession := sanitizePathSegment(sessionID)
	return filepath.Join(lp.basePath, safeUser, safeSession+".log")
}

func (lp *LogPersistence) AppendLog(userID string, sessionID string, entryType string, content string) error {
	path := lp.GetLogFilePath(userID, sessionID)
	f, err := lp.getFileHandle(path)
	if err != nil {
		return err
	}
	ts := time.Now().Format("2006-01-02T15:04:05")
	_, err = f.WriteString(fmt.Sprintf("[%s] [%s] %s\n", ts, entryType, content))
	return err
}

func (lp *LogPersistence) FinalizeLog(userID string, sessionID string, totalLen int) error {
	path := lp.GetLogFilePath(userID, sessionID)
	f, err := lp.getFileHandle(path)
	if err != nil {
		return err
	}
	ts := time.Now().Format("2006-01-02T15:04:05")
	if _, err := f.WriteString(fmt.Sprintf("\n=== SESSION_END %s ===\nTotal Length: %d\n", ts, totalLen)); err != nil {
		return err
	}
	lp.mu.Lock()
	defer lp.mu.Unlock()
	if h, ok := lp.fileHandles[path]; ok {
		_ = h.Close()
		delete(lp.fileHandles, path)
	}
	return nil
}

func (lp *LogPersistence) ReadSessionHistory(userID string, sessionID string, limit int) ([]LogEntry, error) {
	path := lp.GetLogFilePath(userID, sessionID)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []LogEntry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "[") {
			continue
		}
		entry, ok := parseLogLine(line)
		if ok {
			lines = append(lines, entry)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if limit <= 0 || limit >= len(lines) {
		return lines, nil
	}
	return lines[len(lines)-limit:], nil
}

func (lp *LogPersistence) getFileHandle(path string) (*os.File, error) {
	lp.mu.RLock()
	if f, ok := lp.fileHandles[path]; ok {
		lp.mu.RUnlock()
		return f, nil
	}
	lp.mu.RUnlock()

	lp.mu.Lock()
	defer lp.mu.Unlock()
	if f, ok := lp.fileHandles[path]; ok {
		return f, nil
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	lp.fileHandles[path] = f
	return f, nil
}

var pathSegmentRE = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func sanitizePathSegment(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "unknown"
	}
	s = pathSegmentRE.ReplaceAllString(s, "_")
	s = strings.Trim(s, "._-")
	if s == "" {
		return "unknown"
	}
	return s
}

func parseLogLine(line string) (LogEntry, bool) {
	endTs := strings.Index(line, "]")
	if endTs <= 1 {
		return LogEntry{}, false
	}
	ts := line[1:endTs]
	rest := strings.TrimSpace(line[endTs+1:])
	if !strings.HasPrefix(rest, "[") {
		return LogEntry{}, false
	}
	endType := strings.Index(rest, "]")
	if endType <= 1 {
		return LogEntry{}, false
	}
	typ := rest[1:endType]
	content := strings.TrimSpace(rest[endType+1:])
	return LogEntry{Timestamp: ts, Type: typ, Content: content}, true
}
