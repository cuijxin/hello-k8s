package database

import (
	"hello-k8s/pkg/model/common"
)

var DB Database

type Database interface {
	Init(option DBInitOptions) error
	Store(addon *common.AtomApplication) error
	Get(option RecordOptions) (*common.AtomApplication, error)
	List(option RecordOptions) ([]common.AtomApplication, error)
	ListAddons(option RecordOptions) ([]common.AtomApplication, error)
	Exist(option RecordOptions) (bool, error)
	Delete(option RecordOptions) error
	Update(app *common.AtomApplication) error
}

type DBInitOptions struct {
	User     string
	Password string
	Address  string
}

type RecordOptions struct {
	Name      string
	Namespace string
	ClusterID string
	Type      common.AppType
}
