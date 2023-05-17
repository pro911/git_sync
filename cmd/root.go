/*
Copyright © 2023 pro911 pro911@qq.com

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
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// Conf 全局变量 用来保存程序的所有配置信息
var Conf = new(Config)

type Config struct {
	*GitConfig `mapstructure:"git"`
}

type GitConfig struct {
	Cron    string `mapstructure:"cron"`
	Dir     string `mapstructure:"dir"`
	Email   string `mapstructure:"email"`
	Name    string `mapstructure:"name"`
	Address string `mapstructure:"address"`
	Commit  string `mapstructure:"commit"`
	Branch  string `mapstructure:"branch"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gitSync",
	Short: "这是一个git定时同步工具",
	Long: `这是一个git定时工具用来定时同步本地仓库数据到github上，例如我们本地的markdown笔记,我们可以利用github来做为资料的免费存储数据:
	
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {
	//
	//},
	Run: start,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.git_sync.yaml)")
	rootCmd.PersistentFlags().StringP("cron", "c", "", "Directory to sync (default is D:/wwwroot/private/sync-folder)")
	rootCmd.PersistentFlags().StringP("dir", "d", "", "Directory to sync (default is D:/wwwroot/private/sync-folder)")
	rootCmd.PersistentFlags().StringP("email", "e", "", "Directory to sync (default is D:/wwwroot/private/sync-folder)")
	rootCmd.PersistentFlags().StringP("name", "n", "", "Directory to sync (default is D:/wwwroot/private/sync-folder)")
	rootCmd.PersistentFlags().StringP("address", "a", "", "Directory to sync (default is D:/wwwroot/private/sync-folder)")
	rootCmd.PersistentFlags().StringP("commit", "t", "", "Directory to sync (default is D:/wwwroot/private/sync-folder)")
	rootCmd.PersistentFlags().StringP("branch", "b", "", "Directory to sync (default is D:/wwwroot/private/sync-folder)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		//// Find home directory.
		//home, err := os.UserHomeDir()
		//log.Println(home)
		//cobra.CheckErr(err)
		//
		//// Search config in home directory with name ".git_sync" (without extension).
		//viper.AddConfigPath(home)
		//viper.SetConfigType("yaml")
		//viper.SetConfigName(".git_sync")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// 添加配置项
	viper.SetDefault("git.cron", "*/1 * * * *")
	viper.SetDefault("git.dir", "D:/wwwroot/private/sync-folder")
	viper.SetDefault("git.email", "pro911@qq.com")
	viper.SetDefault("git.name", "pro911")
	viper.SetDefault("git.address", "git@github.com:pro911/sync-folder.git")
	viper.SetDefault("git.commit", "Auto-commit from Go")
	viper.SetDefault("git.branch", "main")

	//把读取到的配置信息,反序列化到Conf全局变量中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal feiled,err:%v\n", err)
	}
}

func start(cmd *cobra.Command, args []string) {
	fmt.Println(Conf.GitConfig)
	// 获取命令行传参或 Viper 配置文件参数的值
	if cmd.Flags().Changed("cron") {
		Conf.GitConfig.Cron, _ = cmd.Flags().GetString("cron")
	}
	if cmd.Flags().Changed("dir") {
		Conf.GitConfig.Dir, _ = cmd.Flags().GetString("dir")
	}
	if cmd.Flags().Changed("email") {
		Conf.GitConfig.Email, _ = cmd.Flags().GetString("email")
	}
	if cmd.Flags().Changed("name") {
		Conf.GitConfig.Name, _ = cmd.Flags().GetString("name")
	}
	if cmd.Flags().Changed("address") {
		Conf.GitConfig.Address, _ = cmd.Flags().GetString("address")
	}
	if cmd.Flags().Changed("commit") {
		Conf.GitConfig.Commit, _ = cmd.Flags().GetString("commit")
	}
	if cmd.Flags().Changed("branch") {
		Conf.GitConfig.Branch, _ = cmd.Flags().GetString("branch")
	}

	// 设置 Git 配置
	gitConfig := exec.Command("git", "config", "--global", "user.email", Conf.GitConfig.Email)
	if err := gitConfig.Run(); err != nil {
		log.Fatal(err)
	}
	gitConfig = exec.Command("git", "config", "--global", "user.name", Conf.GitConfig.Name)
	if err := gitConfig.Run(); err != nil {
		log.Fatal(err)
	}

	// 创建 cron 实例
	c := cron.New()

	// 添加定时任务
	_, err := c.AddFunc(Conf.GitConfig.Cron, func() {

		// 切换到主分支
		_, err := gitExec(Conf.GitConfig.Dir, "checkout", Conf.GitConfig.Branch)
		if err != nil {
			log.Fatal(err)
		}

		// 拉取最新代码
		_, err = gitExec(Conf.GitConfig.Dir, "pull", Conf.GitConfig.Address)
		if err != nil {
			log.Fatal(err)
		}

		// 检查 Git 状态
		output, err := gitExec(Conf.GitConfig.Dir, "status", "--porcelain")
		if err != nil {
			log.Fatal(err)
		}

		// 如果有变更，执行 git add 和 commit 命令
		if len(output) > 0 {
			_, err = gitExec(Conf.GitConfig.Dir, "add", ".")
			if err != nil {
				log.Fatal(err)
			}
			_, err = gitExec(Conf.GitConfig.Dir, "commit", "-m", Conf.GitConfig.Commit)
			if err != nil {
				log.Fatal(err)
			}

			// 推送到 GitHub 指定配置仓库
			_, err = gitExec(Conf.GitConfig.Dir, "push", Conf.GitConfig.Address)
			if err != nil {
				log.Fatal(err)
			}
		}

		log.Println("Done!")
	})

	if err != nil {
		log.Fatal(err)
	}

	// 启动 cron
	c.Start()

	// 等待信号中断
	select {}
}

func gitExec(dir, command string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", append([]string{command}, args...)...)
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return output, nil
}
