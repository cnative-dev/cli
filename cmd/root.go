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
	"time"

	"github.com/cnative-dev/cli/cmd/auth"
	"github.com/cnative-dev/cli/cmd/build"
	"github.com/cnative-dev/cli/cmd/config"
	"github.com/cnative-dev/cli/cmd/project"
	"github.com/cnative-dev/cli/cmd/repo"
	"github.com/cnative-dev/cli/cmd/secret"
	"github.com/cnative-dev/cli/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var version = "dev"

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
		silenceUpdateCkTime := viper.GetTime("silence_update_cktime")
		if silenceUpdateCkTime.Before(time.Now()) { //first time check every day，force and background
			go internal.UpdateClient(version, true)
			tomorrow := time.Now().AddDate(0, 0, 1)
			viper.Set("silence_update_cktime", tomorrow.Unix()/(3600*24)*(3600*24)) //第二天零点
		} else {
			internal.UpdateClient(version, false)
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
	internal.SetVersion(version)
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
