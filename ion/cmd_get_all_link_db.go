package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/intelux/insteon"
	"github.com/spf13/cobra"
)

var getAllLinkDBCmd = &cobra.Command{
	Use:  "get-all-link-db",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		records, err := insteon.DefaultPowerLineModem.GetAllLinkDB(rootCtx)

		if err != nil {
			return err
		}

		w := &tabwriter.Writer{}
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintf(w, "#\tDevice\tGroup\tMode\tLink-data\n")

		for i, record := range records {
			fmt.Fprintf(w, "%d\t%s\t%d\t%s\t%s\n", i, record.ID, record.Group, record.Mode(), hex.EncodeToString(record.LinkData[:]))
		}

		return w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(getAllLinkDBCmd)
}
