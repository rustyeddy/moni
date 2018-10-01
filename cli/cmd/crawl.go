// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/rustyeddy/inventory/inv"
	"github.com/spf13/cobra"
)

// walkCmd represents the walk command
var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Crawl a Page at a URL make note of links",
	Long: `Crawl will crawl a webpage represented by URL, it records other links
response times, aliveness and history`,
	Run: func(cmd *cobra.Command, args []string) {
		// Crawl the URLs that have been supplied
		for _, url := range args {
			inv.Crawl(url)
		}
	},
}

func init() {
	crawlCmd.Flags().IntVarP(&inv.CrawlDepth, "depth", "d", 1, "crawl depth")
	rootCmd.AddCommand(crawlCmd)
}

func crawl(cmd *cobra.Command, args []string) {

}
