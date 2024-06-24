package redis

import (
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"
)

type server struct {
	Addr string

	mu    sync.Mutex
	Store map[any]any
}

func NewServer(options Options) *server {
	addr := ":8080"
	if len(options.Addr) > 2 {
		addr = options.Addr
	}

	return &server{
		Addr:  addr,
		Store: make(map[any]any),
	}
}

func (s *server) Run() error {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	slog.Info("tcp server is running", "port", s.Addr)
	for {
		conn, err := l.Accept()
		if err != nil {
			slog.Error("error acception client conn", "err", err)
			continue
		}
		go func() {
			err = s.HandleClientConn(conn)
			if err != nil {
				slog.Error("handling client conn", "err", err)
			}
		}()
	}
}

func (s *server) HandleClientConn(conn net.Conn) error {
	defer func() {
		slog.Info("closeing client conn", "addr", conn.RemoteAddr())
		conn.Close()
	}()
	slog.Info("handling client conn", "addr", conn.RemoteAddr())

	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		return err
	}

	//serialize the string
	sm := string(buf[:n]) //serialize message
	slog.Info("attempting to Unmarshal", "serialized", sm)
	resp, err := Unmarshal(sm)
	if err != nil {
		return fmt.Errorf("Unmarshal err %v", err)
	}

	//set write timeout
	conn.SetWriteDeadline(time.Now().Add(1 * time.Second))

	switch t := dataType(sm[:1]); t {
	case "array":
		cmd := resp.([]any)[0].(string)
		key, val := resp.([]any)[1].(string), ""
		if s, ok := resp.([]any); ok && len(s) == 3 {
			val = resp.([]any)[2].(string)
		}
		slog.Info("handleCommand "+cmd, "key", key, "val", val)
		if err := s.HandleArrayCommand(cmd, key, val); err != nil {
			// Send error response
			if _, err := conn.Write([]byte(fmt.Sprintf("-error handling Command %v%v", cmd, newLine))); err != nil {
				return err
			}
			slog.Info("send back", "resp", "-error handling Command")
			return err
		}

		// Send OK response
		if _, err := conn.Write([]byte("+OK" + newLine)); err != nil {
			return err
		}
		slog.Info("send back", "resp", "+OK")
		return nil
	default:
		// Send error response
		if _, err := conn.Write([]byte(fmt.Sprintf("-unknow type %v", newLine))); err != nil {
			return err
		}
		return fmt.Errorf("unknow type")
	}
}

func (s *server) HandleArrayCommand(cmd, k, v string) error {
	slog.Info("handle array command", "cmd", cmd, "key", k, "value", v)
	switch cmd {
	case "set":
		return s.set(k, v)
	case "get": //todo
		s, err := s.get(k)
		if err != nil {
			return err
		}
		fmt.Println(s)
		return err
	default:
		return fmt.Errorf("unknow cmd")
	}
}

func (s *server) set(k, v string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Store[k] = v
	return nil
}

func (s *server) get(k string) (any, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if value, exist := s.Store[k]; exist {
		return value, nil
	}
	return "", fmt.Errorf("key not found")
}