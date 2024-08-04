package main

import (
	"Ultra-learn/internal/logger"
	"Ultra-learn/internal/server"
	"fmt"
)

func main() {

	s := server.NewServer()
	logger.Info("Server started")

	err := s.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
