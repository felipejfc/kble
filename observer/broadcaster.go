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

package observer

import "sync"

// Broadcaster struct
type Broadcaster struct {
	observersLock sync.RWMutex
	observers     []Observer
}

// NewBroadcaster ctor
func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		observers: []Observer{},
	}
}

// Add adds a observer to the broadcaster
func (b *Broadcaster) Add(observer Observer) {
	b.observersLock.Lock()
	defer b.observersLock.Unlock()
	b.observers = append(b.observers, observer)
}

// Notify notifies all observers of the nodeUpdate
func (b *Broadcaster) Notify(nodeUpdate *NodeUpdate) {
	b.observersLock.RLock()
	observers := b.observers
	b.observersLock.RUnlock()
	for _, observer := range observers {
		go observer.OnUpdate(nodeUpdate)
	}
}
