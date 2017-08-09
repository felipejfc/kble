// MIT License
//
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

package watchers

import (
	"time"

	"github.com/felipejfc/kble/observer"
	"github.com/sirupsen/logrus"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	api "k8s.io/client-go/pkg/api/v1"
	cache "k8s.io/client-go/tools/cache"
)

// NodeWatcher is a watcher for nodes
type NodeWatcher struct {
	clientset      kubernetes.Interface
	nodeIndex      cache.Indexer
	nodeController cache.Controller
	resyncPeriod   time.Duration
	log            logrus.FieldLogger
	stopCh         chan struct{}
	broadcaster    *observer.Broadcaster
}

// NewNodeWatcher ctor
func NewNodeWatcher(clientset kubernetes.Interface, resyncPeriod time.Duration, broadcaster *observer.Broadcaster, log logrus.FieldLogger) *NodeWatcher {
	l := log.WithFields(logrus.Fields{
		"source": "nodeWatcher",
	})
	return &NodeWatcher{
		clientset:    clientset,
		resyncPeriod: resyncPeriod,
		log:          l,
		broadcaster:  broadcaster,
	}
}

func (n *NodeWatcher) nodeAddEventHandler(obj interface{}) {
	n.log.WithFields(log.Fields{
		"node":  obj,
		"event": "add",
	}).Debug("node watcher event")
	nodeUpdate := &observer.NodeUpdate{
		Node: obj.(*api.Node),
		Op:   observer.ADD,
	}
	n.broadcaster.Notify(nodeUpdate)
}
func (n *NodeWatcher) nodeDeleteEventHandler(obj interface{}) {
	n.log.WithFields(log.Fields{
		"node":  obj,
		"event": "delete",
	}).Debug("node watcher event")
	nodeUpdate := &observer.NodeUpdate{
		Node: obj.(*api.Node),
		Op:   observer.REMOVE,
	}
	n.broadcaster.Notify(nodeUpdate)

}
func (n *NodeWatcher) nodeUpdateEventHandler(oldObj, newObj interface{}) {
	// Not interested in updates
	return
}

// List lists the nodes
func (n *NodeWatcher) List() []*api.Node {
	objList := n.nodeIndex.List()
	nodeInstances := make([]*api.Node, len(objList))
	for i, ins := range objList {
		nodeInstances[i] = ins.(*api.Node)
	}
	return nodeInstances
}

// Start starts the nodeWatcher
func (n *NodeWatcher) Start() error {
	n.log.Infoln("starting node watcher")
	eventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc:    n.nodeAddEventHandler,
		DeleteFunc: n.nodeDeleteEventHandler,
		UpdateFunc: n.nodeUpdateEventHandler,
	}
	lw := cache.NewListWatchFromClient(n.clientset.Core().RESTClient(), "nodes", api.NamespaceAll, fields.Everything())
	n.nodeIndex, n.nodeController = cache.NewIndexerInformer(
		lw,
		&api.Node{}, n.resyncPeriod, eventHandler,
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	)
	n.stopCh = make(chan struct{})
	go n.nodeController.Run(n.stopCh)
	return nil
}

// Stop stops the nodeWatcher
func (n *NodeWatcher) Stop() {
	// TODO verificar esse graceful shutdown
	n.stopCh <- struct{}{}
}
