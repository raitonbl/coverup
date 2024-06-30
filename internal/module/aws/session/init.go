package session

import "github.com/raitonbl/coverup/pkg/api"

const (
	ClientRegistryComponentType = "AwsClientRegistry"
)

type ClientRegistry interface {
	api.Component
	GetClient(clientType, name string) (any, error)
}
