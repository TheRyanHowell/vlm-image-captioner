package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/theryanhowell/vlm-image-captioner/pkg/captioner"

	"github.com/spf13/cobra"
)

var csvOutput bool

var rootCmd = &cobra.Command{
	Use:   "vlm-image-captioner [image paths...]",
	Short: "A CLI tool to caption images using vision language models",
	Long:  `A CLI tool to caption images using vision language models.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := os.Getenv("OPENAI_API_KEY")

		baseURL := os.Getenv("OPENAI_BASE_URL")

		model := os.Getenv("OPENAI_MODEL")

		c := captioner.New(apiKey, baseURL, model)

		var csvWriter *csv.Writer
		if csvOutput {
			csvWriter = csv.NewWriter(os.Stdout)
			defer csvWriter.Flush()
			if err := csvWriter.Write([]string{"imagepath", "caption"}); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
		}

		for _, imagePath := range args {
			caption, err := c.Caption(context.Background(), imagePath)
			if err != nil {
				log.Printf("failed to get caption for %s: %v", imagePath, err)
				continue
			}

			caption = strings.TrimSpace(caption)

			if csvOutput {
				if err := csvWriter.Write([]string{imagePath, caption}); err != nil {
					log.Fatalln("error writing record to csv:", err)
				}
			} else {
				if len(args) > 1 {
					fmt.Fprintf(os.Stdout, "%s: %s\n", imagePath, caption)
				} else {
					fmt.Println(caption)
				}
			}
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&csvOutput, "csv", "c", false, "Output as CSV")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
