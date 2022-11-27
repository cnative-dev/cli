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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cnative-dev/cli/cmd/auth"
	"github.com/cnative-dev/cli/cmd/build"
	"github.com/cnative-dev/cli/cmd/config"
	"github.com/cnative-dev/cli/cmd/project"
	"github.com/cnative-dev/cli/cmd/repo"
	"github.com/cnative-dev/cli/cmd/secret"
	"github.com/cnative-dev/cli/internal"
	"github.com/fatih/color"
	"github.com/kardianos/osext"
	"github.com/sanbornm/go-selfupdate/selfupdate"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/mod/semver"
)

var cfgFile string
var version = "dev"

func getExecRelativeDir(dir string) string {
	filename, _ := osext.Executable()
	path := filepath.Join(filepath.Dir(filename), dir)
	return path
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cnative",
	Short: "CNative 命令行工具",

	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd:    true,
		DisableDefaultCmd:   true,
		DisableNoDescFlag:   true,
		DisableDescriptions: true,
	},
	Version: version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		updateDir := ".cnative" + string(os.PathSeparator)
		os.MkdirAll(getExecRelativeDir(updateDir), 0777)
		var preupdater = &selfupdate.Updater{
			CurrentVersion: version,
			Dir:            updateDir,
			CmdName:        "cnative",
			ForceCheck:     false,
		}
		if preupdater.WantUpdate() {
			if endpoints, err := internal.Endpoints(); err == nil {
				updateBase := endpoints["update"]
				if !strings.HasSuffix(updateBase, "/") {
					updateBase = fmt.Sprintf("%s/", updateBase)
				}
				var updater = &selfupdate.Updater{
					CurrentVersion: version,
					ApiURL:         updateBase,
					BinURL:         updateBase,
					DiffURL:        updateBase,
					Dir:            updateDir,
					CheckTime:      24,
					CmdName:        "cnative",
					ForceCheck:     false,
				}
				if newVersion, err := updater.UpdateAvailable(); err == nil && newVersion != "" {
					// 有新版本
					if !semver.IsValid(newVersion) {
						return
					}
					// 不看 patch，只看功能更新
					var oldV string
					newV := semver.MajorMinor(newVersion)
					if semver.IsValid(updater.CurrentVersion) {
						oldV = semver.MajorMinor(updater.CurrentVersion)
					} else {
						oldV = semver.MajorMinor("v0.0.0")
					}
					if semver.Compare(newV, oldV) > 0 { //应该恒为 true
						color.HiGreen("cnative 有版本更新：%s -> %s\n请执行 sudo cnative update 来更新客户端\n", updater.CurrentVersion, newVersion)
						if semver.Compare(semver.Major(newV), semver.Major(oldV)) <= 0 {
							// 如果没有大版本更新（只有功能更新）那就跳过一段时间再检查，否则每次都提示
							updater.SetUpdateTime()
						}
					}
				} else if err == nil && newVersion == "" { //没有新版本
					updater.SetUpdateTime()
				}
			} else {
				fmt.Fprintln(os.Stderr, err.Error())
				return
			}
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		return
	}
}

func init() {
	if version != "dev" {
		log.SetOutput(ioutil.Discard)
	}

	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(auth.NewAuthCommand())
	rootCmd.AddCommand(build.NewBuildCommand())
	rootCmd.AddCommand(config.NewConfigCommand())
	rootCmd.AddCommand(project.NewProjectCommand())
	rootCmd.AddCommand(secret.NewSecretCommand())
	rootCmd.AddCommand(repo.NewRepoCommand())
	rootCmd.AddCommand(NewRuntimeLogCommand())
	rootCmd.AddCommand(NewUpdateCommand(version))

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cnative/config.yaml)")
	rootCmd.PersistentFlags().MarkHidden("config")
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.SetVersionTemplate(`{{printf "Version: %s\n" .Version}}`)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		path := fmt.Sprintf("%s/.cnative/", home)
		os.MkdirAll(path, 0755)
		// Search config in home directory with name ".cnative" (without extension).
		viper.AddConfigPath(path)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.ReadInConfig()
}
