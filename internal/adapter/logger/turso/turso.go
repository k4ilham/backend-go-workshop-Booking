package turso

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type Logger struct {
	url   string
	token string
	httpc *http.Client
}

func New(url, token string) *Logger {
	return &Logger{url: url, token: token, httpc: &http.Client{Timeout: 5 * time.Second}}
}

type execReq struct {
	Statements []statement `json:"statements"`
}

type statement struct {
	Sql  string        `json:"sql"`
	Args []interface{} `json:"args"`
}

func (l *Logger) Log(action string, detail string, at time.Time) error {
	if l.url == "" {
		return nil
	}
	body := execReq{
		Statements: []statement{
			{
				Sql:  "INSERT INTO activity_logs (action, detail, created_at) VALUES (?, ?, ?)",
				Args: []interface{}{action, detail, at.Format(time.RFC3339)},
			},
		},
	}
	b, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", l.url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if l.token != "" {
		req.Header.Set("Authorization", "Bearer "+l.token)
	}
	resp, err := l.httpc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
