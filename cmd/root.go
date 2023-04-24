package cmd

import (
	"context"
	"os"
	"time"

	"github.com/danztran/sample-admission-validating-webhook/config"
	"github.com/danztran/sample-admission-validating-webhook/pkg/server"
	"github.com/danztran/sample-admission-validating-webhook/pkg/utils"
	"github.com/spf13/cobra"
)

var log = utils.MustGetLogger("cmd")

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "sample-admission-validating-webhook",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		Server := server.MustNew(server.Deps{
			Config: config.Values.Server,
		})

		ctx, cancel := context.WithCancel(context.Background())

		var err error
		go func() {
			err = Server.Run(ctx)
		}()

		utils.WaitToStop()
		log.Infof("terminating...")
		cancel()

		return err
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	var err error
	log.Debugf("%+v", config.Values)

	timeStart := time.Now()
	err = rootCmd.Execute()
	execTime := time.Since(timeStart)
	log.Infof("execution time: %v", execTime)

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
