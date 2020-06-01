package addons

import (
	"hello-k8s/pkg/kubernetes/client"
)

// Add-Ons
type AddOn interface {
	Deploy(c *client.HelloK8SClient, options AddOnOptions) error
}

type AddOnOptions struct {
	AddonName string
}
