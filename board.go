package main

import (
	"fmt"
	"time"
    "net/http"
    "encoding/json"
	"crypto/rand"
	"math/big"
)

func (s *Server) Createnewboard(w http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	GameName, found := vars["gamename"]
	if !found {
		GameName = []string{"chess"};
	}
	Token_num, _ := rand.Int(rand.Reader, big.NewInt(1e18))
	Token := fmt.Sprint(Token_num)
	newboard := Board{
		GameName: GameName[0],
		Moves: []string{},
		Current: 0,
		LastUpdate: 0,
		TimePop: Addwaittime(time.Now()),
		Token: Token,
	}
	s.boards[string(Token)[:9]] = newboard
	go s.Boardautodelete(string(Token)[:9])
	res, _ := json.Marshal(newboard)
	_, err := w.Write([]byte(res))
	if err != nil {
		s.logger.Warn("couldn't send result message")
	}
}

func (s *Server) Getboard(w http.ResponseWriter, req *http.Request) {
	vars := req.URL.Query()
	token, _ := vars["token"]
	myboard, found := s.boards[token[0][:9]]
	if !found {
		http.Error(w, "unknown token", http.StatusNotFound)
		return
	}
	myboard.TimePop = Addwaittime(time.Now())
	s.boards[token[0][:9]] = myboard
	myboard.Token = myboard.Token[:9]

	res, _ := json.Marshal(myboard)
	_, err := w.Write([]byte(res))
	if err != nil {
		s.logger.Warn("couldn't send result message")
		return
	}
}

func (s *Server) Setboard(w http.ResponseWriter, req *http.Request) {
	var myboard Board
	err := json.NewDecoder(req.Body).Decode(&myboard)
	if err != nil {
		http.Error(w, "unresolved request", http.StatusNotFound)
		return
	}
	Token := myboard.Token
	realboard, found := s.boards[Token[:9]]
	if !found {
		http.Error(w, "unknown token", http.StatusNotFound)
		return
	}
	if realboard.Token != Token {
		http.Error(w, "broken token", http.StatusNotFound)
		return
	}
	if myboard.LastUpdate != realboard.LastUpdate {
		http.Error(w, "missed update", http.StatusNotFound)
		return
	}
	myboard.TimePop = Addwaittime(time.Now())
	myboard.LastUpdate++
	s.boards[Token[:9]] = myboard
}

