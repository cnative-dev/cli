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
	"github.com/cnative-dev/cli/internal"
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
			spinner := internal.NewSpinner()
			spinner.Suffix = " 正在更新..."
			spinner.FinalMSG = " 完成\n"
			spinner.Start()
			if err := internal.UpdateClient(version, true); err != nil {
				spinner.FinalMSG = "更新出错，使用 sudo/管理员账号 尝试\n"
			}
			spinner.Stop()

		},
	}
	internal.InitCommand(cmd)
	return cmd
}
