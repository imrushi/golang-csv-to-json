package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serve)

	serve.PersistentFlags().StringP("host", "H", "localhost", "Host to listen on")
	serve.PersistentFlags().StringP("port", "P", "8090", "Port on which api to be served")
}

var serve = &cobra.Command{
	Use:   "serve",
	Short: "Starts the HTTP/2 REST API",
	Run: func(cmd *cobra.Command, args []string) {
		host, err := cmd.Flags().GetString("host")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		http.HandleFunc("/health", Health)
		log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil))
	},
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{status:ok}"))
}
