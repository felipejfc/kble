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

package provider

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// Watcher updates networking provider
type Watcher struct {
	provider Provider
	interval time.Duration
	stop     bool
}

// NewWatcher ctor
func NewWatcher(provider Provider, watcherInterval time.Duration) *Watcher {
	return &Watcher{
		provider: provider,
		interval: watcherInterval,
		stop:     false,
	}
}

// Shutdown stops the watcher
func (w *Watcher) Shutdown() {
	w.stop = true
}

// Run runs the watcher
func (w *Watcher) Run(wg *sync.WaitGroup) {
	log.Infoln("starting provider watcher")
	for !w.stop {
		// TODO only if node is ready
		w.provider.EnsureRoutes()
		time.Sleep(w.interval)
	}
	log.Warnln("stopping provider watcher")
	wg.Done()
}
