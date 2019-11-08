package http

import (
	"context"
	"log"
	http2 "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

func NewCmd(router *Router) *cmd {
	return &cmd{
		router:  router,
		command: new(cobra.Command),
	}
}

type cmd struct {
	router  *Router
	command *cobra.Command
}

func (c *cmd) Cmd() *cobra.Command {
	c.command.Use = "http:serve"
	c.command.PersistentFlags().StringP("host", "", ":80", "Http serve address")
	c.command.Run = func(cmd *cobra.Command, args []string) {
		c.run(cmd, args)
	}

	return c.command
}

func (c *cmd) run(cmd *cobra.Command, args []string) {
	srv := &http2.Server{
		Addr:    cmd.Flag("host").Value.String(),
		Handler: c.router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http2.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")
}