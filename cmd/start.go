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
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/felipejfc/kble/provider"
	"github.com/felipejfc/kble/util"
	"github.com/felipejfc/kble/watchers"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var selectedProvider string
var providerWatcherInterval int
var nodeResyncInterval int
var ifaceName string
var kubeconfig string
var incluster bool

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start starts kble",
	Long:  `start starts kble`,
	Run: func(cmd *cobra.Command, args []string) {
		l := log.New()
		l.Infoln("starting kble...")
		if selectedProvider == "" {
			l.Fatalf("no provider specified")
		}
		p := provider.NewLayer2Provider()
		providerWatcher := provider.NewWatcher(p, time.Duration(providerWatcherInterval)*time.Second)
		ch := make(chan os.Signal)
		defer close(ch)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		waitGroup := &sync.WaitGroup{}
		waitGroup.Add(1)
		go providerWatcher.Run(waitGroup)

		clientset, err := util.GetKubernetesClient(l, incluster, kubeconfig)
		if err != nil {
			l.Fatal("error getting kubernetes client: %s", err.Error())
		}
		nodeWatcher := watchers.NewNodeWatcher(clientset, time.Duration(nodeResyncInterval)*time.Second, l)
		// TODO must send wg?
		nodeWatcher.Start()
		<-ch
		providerWatcher.Shutdown()
		nodeWatcher.Stop()

		waitGroup.Wait()
	},
}

func init() {
	// TODO debug flag
	startCmd.Flags().StringVarP(&selectedProvider, "provider", "p", "", "")
	startCmd.Flags().IntVarP(&providerWatcherInterval, "routes watcher interval (s)", "t", 30, "")
	startCmd.Flags().IntVarP(&nodeResyncInterval, "node watcher resync interval (s)", "r", 30, "")
	startCmd.Flags().StringVarP(&ifaceName, "network interface to use for packet transport", "i", "eth0", "")
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	startCmd.Flags().StringVar(&kubeconfig, "kubeconfig", fmt.Sprintf("%s/.kube/config", home), "path to the kubeconfig file (not needed if using --incluster)")
	startCmd.Flags().BoolVar(&incluster, "incluster", false, "incluster mode (for running on kubernetes)")
	RootCmd.AddCommand(startCmd)
}
