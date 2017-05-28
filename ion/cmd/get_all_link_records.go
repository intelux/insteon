// Copyright © 2017 Julien Kauffmann
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

// getAllLinkRecordsCmd represents the get-all-link-records command.
var getAllLinkRecordsCmd = &cobra.Command{
	Use:   "get-all-link-records",
	Short: "Dumps the PLM's all-link records database",
	Long:  `Displays the list of all all-link records from the PowerLine Modem's database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, time.Second)
		records, err := powerLineModem.GetAllLinkRecords(ctx)

		if err != nil {
			return err
		}

		if records == nil {
			fmt.Printf("The all-link records database is empty.\n")
			return nil
		}

		w := &tabwriter.Writer{}
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)
		defer w.Flush()

		fmt.Fprintf(w, "#\tDevice\tGroup\tMode\tLink-data\n")

		for idx, record := range records {
			fmt.Fprintf(w, "%2d\t%s\t%d\t%s\t%s\n", idx, record.Identity, record.Group, record.Mode(), hex.EncodeToString(record.LinkData[:]))
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(getAllLinkRecordsCmd)
}
