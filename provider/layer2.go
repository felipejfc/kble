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
	"github.com/felipejfc/kble/observer"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/pkg/api/v1"
)

//Layer2Provider struct
type Layer2Provider struct {
	machineID string
}

// NewLayer2Provider ctor
func NewLayer2Provider(machineID string) *Layer2Provider {
	l := &Layer2Provider{
		machineID: machineID,
	}
	return l
}

// EnsureRoutes check if routes are set for all nodes
func (l *Layer2Provider) EnsureRoutes() {

}

// AddRouteForNode setup the route for a node
func (l *Layer2Provider) AddRouteForNode(node *v1.Node) {
	cidr := node.Spec.PodCIDR
	ll := log.WithFields(log.Fields{
		"node":     node.Name,
		"pod cidr": cidr,
	})
	if cidr == "" {
		ll.Warn("node has no pod cidr")
		return
	}
	if node.Status.NodeInfo.MachineID == l.machineID {
		ll.Debug("not adding route as node is self")
		return
	}
	ll.Info("layer2 provider ensuring routes for node")
}

// OnUpdate receives a node update
func (l *Layer2Provider) OnUpdate(nodeUpdate *observer.NodeUpdate) {
	node := nodeUpdate.Node
	log.WithFields(log.Fields{
		"node":      nodeUpdate.Node,
		"operation": nodeUpdate.Op,
		"pod cidr":  node.Spec.PodCIDR,
	}).Debug("updating node")
	l.AddRouteForNode(node)
}
