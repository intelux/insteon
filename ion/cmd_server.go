package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var (
	serverCmdEndpoint = ":7660"
	serverCmdOptimize = false
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a control server",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		webService := insteon.NewWebService(nil, rootConfig)

		if err := webService.Synchronize(rootCtx, serverCmdOptimize); err != nil {
			return err
		}

		server := &http.Server{
			Addr:    serverCmdEndpoint,
			Handler: webService.Handler(),
		}

		go func() {
			<-rootCtx.Done()

			// Try to be nice to the clients at least 1 second.
			shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			server.Shutdown(shutdownCtx)

			// Shutdown no matter what.
			server.Close()
		}()

		fmt.Fprintf(os.Stderr, "Started HTTP web-service on %s.\n", serverCmdEndpoint)

		go webService.Run(rootCtx)

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}

		fmt.Fprintf(os.Stderr, "Stop HTTP web-service.\n")

		return nil
	},
}

func init() {
	serverCmd.Flags().StringVarP(&serverCmdEndpoint, "endpoint", "e", serverCmdEndpoint, "The endpoint to listen on.")
	serverCmd.Flags().BoolVarP(&serverCmdOptimize, "optimize", "o", serverCmdOptimize, "Force optimizations or fail.")

	rootCmd.AddCommand(serverCmd)
}
