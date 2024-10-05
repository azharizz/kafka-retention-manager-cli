package cmd

import (
	"fmt"
	"os"

	"main/internal/bucket"
	"main/internal/redis"
	"main/internal/utils"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	bucketName   string
	destBucket   string
	srcPrefix    string
	destPrefix   string
	redisAddress string
	redisKey     string
)

var rootCmd = &cobra.Command{
	Use:   "cli-app",
	Short: "A CLI app for managing GCS buckets and Redis",
	Long: heredoc.Doc(`
		This CLI app provides functionality to manage Google Cloud Storage buckets
		and interact with Redis. It can move files between buckets or delete files.
	`),
	Run: func(cmd *cobra.Command, args []string) {
		// First, ask the user to choose between moving and deleting
		var operation string
		prompt := &survey.Select{
			Message: "Choose an operation:",
			Options: []string{"Move files", "Delete files"},
		}
		survey.AskOne(prompt, &operation)

		// Common questions for both operations
		questions := []*survey.Question{
			{
				Name: "bucketName",
				Prompt: &survey.Input{
					Message: "Enter the source bucket name:",
					Default: bucketName,
				},
			},
			{
				Name: "srcPrefix",
				Prompt: &survey.Input{
					Message: "Enter the source prefix:",
					Default: srcPrefix,
				},
			},
		}

		// Additional questions for move operation
		if operation == "Move files" {
			questions = append(questions,
				&survey.Question{
					Name: "destBucket",
					Prompt: &survey.Input{
						Message: "Enter the destination bucket name:",
						Default: destBucket,
					},
				},
				&survey.Question{
					Name: "destPrefix",
					Prompt: &survey.Input{
						Message: "Enter the destination prefix:",
						Default: destPrefix,
					},
				},
				&survey.Question{
					Name: "redisAddress",
					Prompt: &survey.Input{
						Message: "Enter the Redis address:",
						Default: redisAddress,
					},
				},
				&survey.Question{
					Name: "redisKey",
					Prompt: &survey.Input{
						Message: "Enter the Redis key:",
						Default: utils.GetRedisKey(),
					},
				},
			)
		}

		answers := struct {
			BucketName   string
			DestBucket   string
			SrcPrefix    string
			DestPrefix   string
			RedisAddress string
			RedisKey     string
		}{}

		err := survey.Ask(questions, &answers)
		if err != nil {
			color.Red("Failed to get user input: %v", err)
			os.Exit(1)
		}

		// Update variables with user input
		bucketName = answers.BucketName
		srcPrefix = answers.SrcPrefix

		if operation == "Move files" {
			destBucket = answers.DestBucket
			destPrefix = answers.DestPrefix
			redisAddress = answers.RedisAddress
			redisKey = answers.RedisKey

			// Execute the move operation
			count, err := bucket.CountFilesInBucket(bucketName, srcPrefix)
			if err != nil {
				color.Red("Failed to count files in bucket: %v", err)
				os.Exit(1)
			}
			color.Green("Total files in bucket '%s': %d\n", bucketName, count)

			rdb := redis.InitializeRedisClient(redisAddress)
			defer rdb.Close()

			countRedis, err := redis.GetValueFromRedis(rdb, redisKey)
			if err != nil {
				color.Red("Failed to get value from Redis: %v", err)
				os.Exit(1)
			}
			color.Green("Total rows in Redis: %d", countRedis)

			if count == countRedis {
				color.Green("The count file matches %v : %v", count, countRedis)
				err = bucket.MoveFilesToBucket(bucketName, destBucket, srcPrefix, destPrefix)
				if err != nil {
					color.Red("Failed to move files: %v", err)
					os.Exit(1)
				}
				color.Green("Files moved successfully")
			} else {
				color.Yellow("The count file does NOT match %v : %v", count, countRedis)
				color.Yellow("Files were not moved due to count mismatch")
			}
		} else {
			// Execute the delete operation
			err = bucket.DeleteFilesWithPrefix(bucketName, srcPrefix)
			if err != nil {
				color.Red("Failed to delete files: %v", err)
				os.Exit(1)
			}
			color.Green("Files deleted successfully")
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli-app.yaml)")
	rootCmd.Flags().StringVar(&bucketName, "bucket", "", "Source bucket name")
	rootCmd.Flags().StringVar(&destBucket, "dest-bucket", "", "Destination bucket name")
	rootCmd.Flags().StringVar(&srcPrefix, "src-prefix", "", "Source prefix")
	rootCmd.Flags().StringVar(&destPrefix, "dest-prefix", "", "Destination prefix")
	rootCmd.Flags().StringVar(&redisAddress, "redis-addr", "localhost:6379", "Redis address")
	rootCmd.Flags().StringVar(&redisKey, "redis-key", "", "Redis key")

	viper.BindPFlag("bucket", rootCmd.Flags().Lookup("bucket"))
	viper.BindPFlag("dest-bucket", rootCmd.Flags().Lookup("dest-bucket"))
	viper.BindPFlag("src-prefix", rootCmd.Flags().Lookup("src-prefix"))
	viper.BindPFlag("dest-prefix", rootCmd.Flags().Lookup("dest-prefix"))
	viper.BindPFlag("redis-addr", rootCmd.Flags().Lookup("redis-addr"))
	viper.BindPFlag("redis-key", rootCmd.Flags().Lookup("redis-key"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cli-app")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
