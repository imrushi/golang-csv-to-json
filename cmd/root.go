package cmd

import (
	"github.com/spf13/cobra"
)

type inputFile struct {
	filepath  string
	separator string
	pretty    bool
}

var (
	//This is used for config file
	rootCmd = &cobra.Command{
		Use:   "csv2json",
		Short: "csv2json: It is simple way to covert Csv to Json file.",
		Long:  `csv2json: It is a CLI application for convert Csv file to Json file with multiple options.`,
	}
)

//Excute the root commands
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("separator", "s", "comma", "Column Separator")
	rootCmd.MarkPersistentFlagRequired("separator")
	rootCmd.PersistentFlags().BoolP("pretty", "p", false, "Generate pretty JSON")
}
