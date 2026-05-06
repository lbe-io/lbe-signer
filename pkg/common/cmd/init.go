package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const InitCmdArgsNum = 2

func InitConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [numbers...]",
		Short: "project init before start",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < InitCmdArgsNum {
				return fmt.Errorf("At least need %d args", InitCmdArgsNum)
			}
			for _, arg := range args {
				// 验证是否为数字
				var num int
				_, err := fmt.Sscanf(arg, "%d", &num)
				if err != nil {
					return fmt.Errorf("无效的数字: %s", arg)
				}
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			result := 1
			for _, arg := range args {
				var num int
				fmt.Sscanf(arg, "%d", &num)
				result *= num
			}
			fmt.Printf("积: %d\n", result)
		},
	}
	cmd.Flags().StringP(FlagConf, "c", "./config", "path of config directory")
	return cmd
}
