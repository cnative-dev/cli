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
package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/cnative-dev/cli/internal"
	"github.com/google/uuid"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewAuthLoginCommand() *cobra.Command {
	var provider string
	var timeout int
	type Token struct {
		Token    string `json:"token"`
		GitToken string `json:"git"`
	}
	responseToken := &Token{}

	var cmd = &cobra.Command{
		Use:   "login",
		Short: "登录 CNative",
		Run: func(cmd *cobra.Command, args []string) {
			var host string
			session := uuid.NewString()

			if base, found := os.LookupEnv("CNATIVE_API"); found {
				host = base
			} else {
				host = "https://api.cnative.dev"
			}

			s := internal.NewSpinner()
			s.Suffix = " process login request..."
			s.FinalMSG = "  process login request...timeout."
			s.Start()
			browser.OpenURL(fmt.Sprintf("%s/auth/%s/authorize?session=%s&expires=%d", host, provider, session, timeout))
			intervals := []int{2, 1, 2, 3, 3, 3, 5, 5, 6}
			for _, interval := range intervals {
				if resp, err := internal.R().
					SetQueryParam("session", session).
					SetResult(responseToken).
					Get("/api/auth"); err == nil && resp.StatusCode() == 200 {
					s.FinalMSG = "  process login request...done."
					viper.Set("token", responseToken.Token)
					// 设置 netrc
					rcs, _ := internal.ReadNetrc()
					filtered := []internal.NetrcLine{}
					for _, rc := range rcs {
						if rc.Machine != "git.cnative.dev" {
							filtered = append(filtered, rc)
						}
					}
					filtered = append(filtered, internal.NetrcLine{
						Machine:  "git.cnative.dev",
						Login:    "cnative",
						Password: responseToken.GitToken,
					})
					internal.WriteNetrc(filtered)
					break
				}
				time.Sleep(time.Duration(interval) * time.Second)
			}
			s.Stop()
		},
	}
	internal.InitCommand(cmd)

	cmd.Flags().IntVarP(&timeout, "timeout", "t", 3600, "登录凭证过期时间，单位秒，默认1小时")
	cmd.Flags().StringVarP(&provider, "provider", "p", "", "第三方服务登陆，可选值：[github]")
	cmd.MarkFlagRequired("provider")

	return cmd
}
