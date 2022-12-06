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
package build

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/cnative-dev/cli/internal"
	"github.com/cnative-dev/cli/internal/flags"
	"github.com/spf13/cobra"
)

func NewBuildLogCommand() *cobra.Command {
	var project string
	var follow bool
	var lines int
	var cmd = &cobra.Command{
		Use:     "log",
		Aliases: []string{"logs", "log"},
		Short:   "查看构建日志",
		Run: func(cmd *cobra.Command, args []string) {
			internal.WithAuthorized(func() {
				ignoreCloseError := false
				s := internal.NewSpinner()
				resp, _ := internal.R().
					SetPathParam("projectId", project).
					SetPathParam("buildId", args[0]).
					SetDoNotParseResponse(true).
					Get("/api/projects/{projectId}/builds/{buildId}/log")
				s.Start()
				defer func() {
					s.Stop()
					ignoreCloseError = true
					resp.RawResponse.Body.Close()
				}()
				defer resp.RawResponse.Body.Close()
				reader := bufio.NewReader(resp.RawResponse.Body)
				buffer := []string{}
				for {
					line, prefix, err := reader.ReadLine()
					if line != nil {
						if prefix {
							buffer = append(buffer, string(line))
						} else {
							buffer = append(buffer, string(line))
							s.Disable()
							fmt.Println(strings.Join(buffer, ""))
							s.Enable()
							buffer = []string{}
						}
					} else if ignoreCloseError {
						break
					} else if err == io.EOF {
						if len(buffer) > 0 {
							s.Disable()
							fmt.Println(buffer)
							s.Enable()
						}
						break
					} else if err != nil {
						panic(err)
					} else {
						break
					}
				}
			})
		},
	}
	internal.InitCommand(cmd)
	internal.UsageTemplateWithArgs(cmd, "构建ID")
	flags.AddProject(cmd, &project)
	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "持续读取 log 流")
	cmd.Flags().IntVarP(&lines, "lines", "n", 10, "打印最后 n 行")
	return cmd
}
