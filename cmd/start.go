// Copyright (c) 2017 Felipe Cavalcanti (fjfcavalcanti@gmail.com)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cmd

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/felipejfc/kble/provider"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var selectedProvider string
var interval int

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start starts kble",
	Long:  `start starts kble`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infoln("starting kble...")
		if selectedProvider == "" {
			log.Fatalf("no provider specified")
		}
		p := provider.NewLayer2Provider()
		watcher := provider.NewWatcher(p, time.Duration(interval)*time.Second)

		ch := make(chan os.Signal)
		defer close(ch)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

		waitGroup := &sync.WaitGroup{}
		waitGroup.Add(1)
		go watcher.Run(waitGroup)
		<-ch
		watcher.Shutdown()
		waitGroup.Wait()
	},
}

func init() {
	startCmd.Flags().StringVarP(&selectedProvider, "provider", "p", "", "")
	startCmd.Flags().IntVarP(&interval, "watcher interval (s)", "t", 30, "")
	RootCmd.AddCommand(startCmd)
}
