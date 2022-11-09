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
package project

import (
	"fmt"
	"os"

	"github.com/cnative-dev/cli/internal"
	"github.com/spf13/cobra"
)

func NewProjectDescribeCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "describe",
		Aliases: []string{"desc"},
		Short:   "查看一个项目详情",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			internal.WithAuthorized(func() {
				project := &Project{}
				if resp, err := internal.R().
					SetResult(project).
					SetPathParam("projectId", args[0]).
					Get("/api/projects/{projectId}"); err == nil && resp.StatusCode() == 200 {
					internal.PrettyStruct(project)
				} else {
					fmt.Fprintln(os.Stderr, resp.Error().(*internal.ErrResp).Error)
					os.Exit(1)
				}
			})
		},
	}

	internal.InitCommand(cmd)
	internal.UsageTemplateWithArgs(cmd, "项目名")
	return cmd
}
