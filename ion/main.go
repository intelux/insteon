package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var (
	rootCtx, rootCtxCancel = withInterrupt(context.Background())
	rootConfig             *insteon.Configuration
)

var rootCmd = &cobra.Command{
	Use: "ion",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		rootConfig, err = insteon.LoadDefaultConfiguration()

		return
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		rootCtxCancel()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func withInterrupt(ctx context.Context) (context.Context, func()) {
	ctx, cancel := context.WithCancel(ctx)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	go func() {
		<-ch
		cancel()
	}()

	go func() {
		<-ctx.Done()

		signal.Stop(ch)
		close(ch)
	}()

	return ctx, cancel
}
