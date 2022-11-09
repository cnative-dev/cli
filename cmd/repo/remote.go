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
package repo

import (
	"fmt"

	"github.com/cnative-dev/cli/internal"
	"github.com/cnative-dev/cli/internal/flags"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/spf13/cobra"
)

func NewRepoRemoteCommand() *cobra.Command {
	var force bool
	var path string
	var project string
	var name string
	var cmd = &cobra.Command{
		Use:   "remote",
		Short: "链接本地 Git 仓库",
		Example: `  # 在当前目录下设定远程仓库 cnative
  cnative repo remote -p xxx

  # 在指定目录下设定远程仓库 cnative
  cnative repo remote -p xxx --path ~/project/demo
`,
		Run: func(cmd *cobra.Command, args []string) {
			internal.WithAuthorized(func() {
				if repo, err := git.PlainOpen(path); err == nil {
					if remote, e := repo.Remote(name); force || e == git.ErrRemoteNotFound {
						link := GitLink{}
						if resp, err := internal.R().
							SetResult(&link).
							SetPathParam("projectId", project).
							Get("/api/projects/{projectId}"); err == nil && resp.StatusCode() == 200 {
							if remote != nil {
								repo.DeleteRemote(name)
							}
							repo.CreateRemote(&config.RemoteConfig{
								Name: name,
								URLs: []string{link.URL},
							})
							fmt.Printf("远程仓库 %s 已配置。\n", name)
						} else {
							fmt.Printf("远程仓库 %s 不存在,请检查\n", project)
						}
					} else {
						fmt.Printf("当前仓库已经绑定远程仓库(%s)，如果需要覆盖请用 -f 参数\n", remote.Config().URLs[0])
					}
				} else if err == git.ErrRepositoryNotExists {
					fmt.Printf("该目录(%s)无 git 仓库\n", path)
				}
			})
		},
	}
	internal.InitCommand(cmd)
	flags.AddProject(cmd, &project)
	cmd.Flags().BoolVarP(&force, "force", "f", false, "强制覆盖本地绑定的 CNative 远程仓库")
	cmd.Flags().StringVarP(&name, "name", "n", "cnative", "远程分支名")
	cmd.Flags().StringVar(&path, "path", ".", "git 仓库路径")
	return cmd
}
