package main

import (
	//"os"
	"fmt"
	"time"
	"log/slog"
    "net/http"
    "encoding/json"
	"crypto/rand"
	"math/big"
    //"strings"
)

type Board struct {
	GameName string `json:"gamename"`
	Moves []string `json:"moves"`
	Current int `json:"current"`
	TimePop time.Time `json:"timepop"`
	Token string `json:"token"`
}

type Server struct {
	logger *slog.Logger
	boards map[string]Board
}

func newServer(logger *slog.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

func (s *Server) Boardautodelete(token string) {
	myboard, found := s.boards[token]
	if !found {
		s.logger.Warn("unknown token when autodeleting", token)
		return
	}
	if time.Now().Before(myboard.TimePop) {
		select {
		case <-time.After(time.Until(myboard.TimePop)):
			s.Boardautodelete(token)
		}
	} else {
		delete(s.boards, token)
	}
}

func Addwaittime(timepop time.Time) time.Time {
	return timepop.Add(time.Minute * time.Duration(15))
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	s.boards = make(map[string]Board)
	
	slog.Info("Starting server on port 8080")
	mux.HandleFunc("/api/board/ping", s.Ping)
	mux.HandleFunc("/api/board/createnew", s.Createnewboard)
	mux.HandleFunc("/api/board/get", s.Getboard)
	err := http.ListenAndServe(":8080", mux)
	if err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Ping(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write([]byte("pong"))
	if err != nil {
		s.logger.Warn("couldn't send result message")
	}
}

func (s *Server) Createnewboard(w http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	GameName, found := vars["gamename"]
	if !found {
		GameName = []string{"chess"};
	}
	Token_num, _ := rand.Int(rand.Reader, big.NewInt(1e16))
	Token := fmt.Sprint(Token_num)
	newboard := Board{
		GameName: GameName[0],
		Moves: []string{},
		Current: 0,
		TimePop: Addwaittime(time.Now()),
		Token: Token,
	}
	s.boards[string(Token)] = newboard
	go s.Boardautodelete(string(Token))
	res, _ := json.Marshal(newboard)
	_, err := w.Write([]byte(res))
	if err != nil {
		s.logger.Warn("couldn't send result message")
	}
}

func (s *Server) Getboard(w http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	token,   _     := vars["token"]
	myboard, found := s.boards[token[0]]
	if !found {
		http.Error(w, "unknown token", http.StatusNotFound)
		return
	}
	myboard.TimePop    = Addwaittime(time.Now())
	s.boards[token[0]] = myboard

	res, _ := json.Marshal(myboard)
	_, err := w.Write([]byte(res))
	if err != nil {
		s.logger.Warn("couldn't send result message")
		return
	}
}


func main() {
	logger     := slog.Default()
	mainServer := newServer(logger)
	err        := mainServer.Run()
	if err != nil {
		logger.Error("server has been stopped", "error", err)
	}
}

// api:
//    /api/ping -> pong
//    /api/board/createnew?gamename=<optional> -> {"gamename":"chess","moves":[],"current":0,"timepop":"<date>","token":"<token>"}
//    /api/board/get?token=<token> -> same as api/board/createnew
