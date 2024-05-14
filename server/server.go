package server

import (
	"context"
	"errors"
	"github.com/dmzlingyin/utils/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Serve(port string, handler http.Handler) {
	s := &http.Server{
		Addr:           port,
		Handler:        handler,
		ReadTimeout:    5 * time.Minute,
		WriteTimeout:   5 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Infof("server listening on %s", port)
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("server exited!!! err: %s", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Errorf("server shutdown err: %s", err)
	}
	log.Info("server exiting")
}
