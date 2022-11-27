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
package config

import (
	"fmt"
	"os"

	"github.com/cnative-dev/cli/internal"
	"github.com/cnative-dev/cli/internal/flags"
	"github.com/spf13/cobra"
)

func NewConfigRemoveCommand() *cobra.Command {
	var project string
	// configListCmd represents the list command
	var cmd = &cobra.Command{
		Use:   "remove",
		Short: "删除一条配置项",
		Run: func(cmd *cobra.Command, args []string) {
			internal.WithAuthorized(func() {
				if resp, err := internal.R().
					SetPathParam("key", args[0]).
					SetPathParam("project", project).
					Delete("/api/projects/{project}/configs/{key}"); err == nil &&
					resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
					fmt.Printf("配置 %s 已删除\n", args[0])
				} else {
					fmt.Fprintln(os.Stderr, resp.Error().(*internal.ErrResp).Details)
					return
				}
			})
		},
	}
	flags.AddProject(cmd, &project)
	internal.InitCommand(cmd)
	internal.UsageTemplateWithArgs(cmd, "KEY")

	return cmd
}
