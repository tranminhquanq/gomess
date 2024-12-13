package cmd

import (
	"github.com/spf13/cobra"
)

var (
	configFile = ""
)

var rootCmd = cobra.Command{
	Use: "gomess",
	Run: func(cmd *cobra.Command, args []string) {
		serve(cmd.Context())
	},
}

func RootCommand() *cobra.Command {
	rootCmd.AddCommand(&serveCmd, &migrateCmd, &versionCmd)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "base configuration file to load")

	return &rootCmd
}

// func loadGlobalConfig(ctx context.Context) *config.GlobalConfiguration {
// 	if ctx == nil {
// 		panic("context must not be nil")
// 	}

// 	cfg, err := config.LoadGlobal(configFile)
// 	if err != nil {
// 		logrus.Fatalf("Failed to load configuration: %+v", err)
// 	}

// 	// if err := observability.ConfigureLogging(&config.Logging); err != nil {
// 	// 	logrus.WithError(err).Error("unable to configure logging")
// 	// }

// 	// if err := observability.ConfigureTracing(ctx, &config.Tracing); err != nil {
// 	// 	logrus.WithError(err).Error("unable to configure tracing")
// 	// }

// 	// if err := observability.ConfigureMetrics(ctx, &config.Metrics); err != nil {
// 	// 	logrus.WithError(err).Error("unable to configure metrics")
// 	// }

// 	// if err := observability.ConfigureProfiler(ctx, &config.Profiler); err != nil {
// 	// 	logrus.WithError(err).Error("unable to configure profiler")
// 	// }

// 	return cfg
// }

// func execWithConfigAndArgs(cmd *cobra.Command, fn func(config *config.GlobalConfiguration, args []string), args []string) {
// 	fn(loadGlobalConfig(cmd.Context()), args)
// }
