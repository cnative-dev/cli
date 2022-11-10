/*
Copyright © 2022 Jingsi <lookisliu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/cnative-dev/cli/internal"
	"github.com/sanbornm/go-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

func NewUpdateCommand(version string) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "update",
		Short: "更新 CLI 客户端",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			//override root update check
		},
		Run: func(cmd *cobra.Command, args []string) {
			if endpoints, err := internal.Endpoints(); err == nil {
				updateBase := endpoints["update"]
				if !strings.HasSuffix(updateBase, "/") {
					updateBase = fmt.Sprintf("%s/", updateBase)
				}
				if home, err := os.UserHomeDir(); err == nil {
					var updater = &selfupdate.Updater{
						CurrentVersion: version,
						ApiURL:         updateBase,
						BinURL:         updateBase,
						DiffURL:        updateBase,
						Dir:            fmt.Sprintf("%s/.cnative/", home),
						CheckTime:      24,
						CmdName:        "cnative",
						ForceCheck:     true,
					}
					if newV, e := updater.UpdateAvailable(); e == nil && newV != "" {
						spinner := internal.NewSpinner()
						spinner.Suffix = " 正在更新..."
						spinner.FinalMSG = " 完成"
						spinner.Start()
						if err := updater.BackgroundRun(); err != nil {
							spinner.FinalMSG = "更新出错，请稍后再试"
						}
						spinner.Stop()
					}
				}
			} else {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
		},
	}
	internal.InitCommand(cmd)
	return cmd
}
