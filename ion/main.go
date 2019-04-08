package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

var (
	rootCtx, rootCtxCancel = withInterrupt(context.Background())
)

var rootCmd = &cobra.Command{
	Use: "ion",
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
