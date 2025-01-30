package main

import (
	"os"
	"time"
	"log/slog"
    "net/http"
)

type Board struct {
	GameName string `json:"gamename"`
	Moves []string `json:"moves"`
	Current int `json:"current"`
	LastUpdate int `json:"lastupdate"`
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
	mux.HandleFunc("/", s.SendClientFile)
	mux.HandleFunc("/client/{filename}", s.SendClientFile)
	mux.HandleFunc("/images/{filename}", s.SendClientFile)
	mux.HandleFunc("/favicon.ico", s.SendClientFile)
	mux.HandleFunc("/api/ping", s.Ping)

	mux.HandleFunc("/api/board/createnew", s.Createnewboard)
	mux.HandleFunc("/api/board/get", s.Getboard)
	mux.HandleFunc("/api/board/set", s.Setboard)

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

func (s *Server) SendClientFile(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Path
	if filename == "/" {
		filename = "/index.html"
	}
	res, err := os.ReadFile(filename)
	if err != nil {
		s.logger.Warn("couldn't open file ", filename)
		return
	}
	_, err = w.Write(res)
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

//    /api/board/createnew?gamename=<optional> -> {"gamename":"chess","moves":[],"current":0,"lastupdate":"<date>","timepop":"<date>","token":"<token>"}
//    /api/board/get?token=<token> -> same as api/board/createnew
//    /api/board/set (body as json body) -> success

