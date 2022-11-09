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
	"strings"

	"github.com/cnative-dev/cli/internal"
	"github.com/cnative-dev/cli/internal/flags"
	"github.com/spf13/cobra"
)

func NewConfigApplyCommand() *cobra.Command {
	var vars []string
	var varFile string
	var project string
	var cmd = &cobra.Command{
		Use:   "apply",
		Short: "创建或替换项目配置文件，注意：会重新部署线上服务",
		Example: `  # 从本地文件创建配置（properties 格式）
  cnative config apply --from-file ./config.properties -p xxx

  # 指定变量创建配置文件
  cnative config apply -v HELLO=world -v name=CNATIVE -p xxx

  # 使用文件创建配置，然后使用变量覆盖其中部分配置
  cnative config apply --from-file ./config.properties -v HELLO=world -v name=CNATIVE -p xxx
`,
		Run: func(cmd *cobra.Command, args []string) {
			internal.WithAuthorized(func() {
				if properties, err := internal.ReadProperties(varFile); err != nil {
					panic(err)
				} else {
					for _, line := range vars {
						if equal := strings.Index(line, "="); equal >= 0 {
							if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
								value := ""
								if len(line) > equal {
									value = strings.TrimSpace(line[equal+1:])
								}
								// assign the config map
								properties[key] = value
							}
						}
					}
					config := &Config{}
					if resp, err := internal.R().
						SetResult(config).
						SetBody(properties).
						SetPathParam("project", project).
						Post("/api/projects/{project}/configs"); err == nil && resp.StatusCode() == 200 {
						internal.PrettyMapAsArray(*config)
						fmt.Println("完成")
					} else {
						fmt.Fprintln(os.Stderr, resp.Error().(*internal.ErrResp).Error)
						os.Exit(1)
					}
				}
			})
		},
	}
	internal.InitCommand(cmd)
	flags.AddProject(cmd, &project)
	cmd.Flags().StringVar(&varFile, "from-file", "", "从本地文件创建配置")
	cmd.Flags().StringArrayVarP(&vars, "var", "v", nil, "配置变量，该参数可以多次使用。")

	return cmd
}
