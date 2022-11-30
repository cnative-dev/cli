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
	"encoding/json"
	"fmt"
	"strings"

	"os"

	"github.com/cnative-dev/cli/internal"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/spf13/cobra"
)

func NewRepoPushCommand() *cobra.Command {
	var force bool
	var cmd = &cobra.Command{
		Use:   "push",
		Short: "推送本地代码变更",
		Example: `  # 将本地代码推送至 cnative 进行部署
  cnative repo push demo
`,
		Run: func(cmd *cobra.Command, args []string) {
			internal.WithAuthorized(func() {
				if resp, err := internal.R().Get("/api/group/token"); err == nil {
					token := map[string]string{}
					json.Unmarshal(resp.Body(), &token)
					if repo, err := git.PlainOpen("."); err == nil {
						//打开了仓库，拿到了token
						auth := &http.BasicAuth{
							Username: "cnative",
							Password: token["token"],
						}
						//准备远程分支
						remoteName := fmt.Sprintf("cnative-%s", internal.RandSeq(8))
						link := GitLink{}
						if resp, err := internal.R().
							SetResult(&link).
							SetPathParam("projectName", args[0]).
							Get("/api/projects/{projectName}"); err == nil && resp.StatusCode() == 200 {
							//创建远程分支
							repo.CreateRemote(&config.RemoteConfig{
								Name: remoteName,
								URLs: []string{link.URL},
							})
							defer repo.DeleteRemote(remoteName)
							currentReference, _ := repo.Head()
							refSpec := fmt.Sprintf("refs/heads/%s:refs/heads/main", currentReference.Name().Short())
							options := &git.PushOptions{
								RefSpecs:   []config.RefSpec{config.RefSpec(refSpec)},
								RemoteName: remoteName,
								Auth:       auth,
								Force:      force,
							}
							fmt.Printf("正在推送本地分支 %s, 到项目 %s\n", currentReference.Name().Short(), args[0])
							spinner := internal.NewSpinner()
							internal.AddCallback(func() {
								repo.DeleteRemote(remoteName)
							})
							spinner.Start()
							//推送
							if err := repo.Push(options); err != nil {
								if err == git.ErrForceNeeded || strings.HasPrefix(err.Error(), "non-fast-forward") {
									spinner.Disable()
									fmt.Fprintln(os.Stderr, "远程提交历史记录与本地不同步，请确认本地代码和项目的对应关系，强制推送请用 force 参数")
									spinner.Enable()
								} else {
									spinner.Disable()
									fmt.Fprintln(os.Stderr, err.Error())
									spinner.Enable()
								}
								spinner.Stop()
							} else {
								spinner.Stop()
								fmt.Println("完成.")
							}
						}

					} else if err == git.ErrRepositoryNotExists {
						fmt.Println("当前目录无 git 仓库")
					}
				} else {
					fmt.Fprintln(os.Stderr, err.Error())
				}
			})
		},
	}
	internal.InitCommand(cmd)
	internal.UsageTemplateWithArgs(cmd, "项目名")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "强制推送")
	return cmd
}
