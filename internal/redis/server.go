package redis

import (
	"fmt"
	"log/slog"
	"net"
	"reflect"
	"sync"
	"time"
)

type storeValue struct {
	Value any
	//Exp   time.Duration //expire after
	Exp string //expire after
}

type server struct {
	Addr string

	mu    sync.Mutex
	Store map[any]storeValue
}

func NewServer(options Options) *server {
	addr := ":8080"
	if len(options.Addr) > 2 {
		addr = options.Addr
	}

	return &server{
		Addr:  addr,
		Store: make(map[any]storeValue),
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
		var val, expireAfter string
		cmd := resp.([]any)[0].(string)
		key := resp.([]any)[1].(string)

		if s, ok := resp.([]any); ok {
			if len(s) > 2 {
				val = resp.([]any)[2].(string)
			}

			if len(s) > 4 {
				expireAfterFlag := resp.([]any)[3].(string)
				if expireAfterFlag == "EXP" {
					expireAfter = resp.([]any)[4].(string)
				}
			}
		}

		slog.Info("handleCommand "+cmd, "key", key, "val", val)
		s, err := s.HandleArrayCommand(cmd, key, val, expireAfter)
		if err != nil {
			// Send error response
			if _, err := conn.Write([]byte(fmt.Sprintf("-error handling Command %v%v %v", cmd, newLine, err))); err != nil {
				return err
			}
			slog.Info("send back", "resp", "-error handling Command")
			return err
		}

		if s != "" && cmd == "get" {
			//get response
			encodedString, err := Marshal(s)
			if err != nil {
				return fmt.Errorf("error marshal the string %v", err)
			}
			if _, err := conn.Write(encodedString); err != nil {
				return err
			}
		} else {
			// Send OK response
			// set key value
			if _, err := conn.Write([]byte("+OK" + newLine)); err != nil {
				return err
			}
			slog.Info("send back", "resp", "+OK")
		}

		return nil
	default:
		// Send error response
		if _, err := conn.Write([]byte(fmt.Sprintf("-unknow type %v", newLine))); err != nil {
			return err
		}
		return fmt.Errorf("unknow type")
	}
}

func (s *server) HandleArrayCommand(cmd, k, v, expireAfter string) (string, error) {
	slog.Info("handle array command", "cmd", cmd, "key", k, "value", v)
	switch cmd {
	case "set":
		return "", s.set(k, v, expireAfter)
	case "get": //todo
		s, err := s.get(k)
		if err != nil {
			return "", err
		}
		if str, ok := s.(string); ok {
			return str, nil
		}
		return "", fmt.Errorf("marshal error %v type %v", s, reflect.TypeOf(s))
	default:
		return "", fmt.Errorf("unknow cmd")
	}
}

func (s *server) set(k, v, expireAfter string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Store[k] = storeValue{Value: v, Exp: expireAfter}
	return nil
}

func (s *server) get(k string) (any, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if value, exist := s.Store[k]; exist {
		return value.Value, nil
	}
	return "", fmt.Errorf("key not found")
}
