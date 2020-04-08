module hello-k8s

go 1.13

require (
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/cuijxin/postgres-operator-atom v0.0.1
	github.com/docker/distribution v2.7.1+incompatible
	github.com/emicklei/go-restful v0.0.0-20170410110728-ff4f55a20633
	github.com/fsnotify/fsnotify v1.4.7
	github.com/gin-gonic/gin v1.5.0
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/gorilla/websocket v1.4.0
	github.com/jinzhu/gorm v1.9.12
	github.com/lexkong/log v0.0.0-20180607165131-972f9cd951fc
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/oracle/mysql-operator v0.0.0-20190515081336-9aebcc37a080
	github.com/prometheus/client_golang v0.9.3
	github.com/shirou/gopsutil v2.20.1+incompatible
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.2
	github.com/spotahome/redis-operator v1.0.0-rc.4
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.5
	github.com/unknwon/com v1.0.1
	golang.org/x/crypto v0.0.0-20200214034016-1d94cc7ab1c6 // indirect
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	golang.org/x/text v0.3.2
	gopkg.in/igm/sockjs-go.v2 v2.0.1
	gopkg.in/square/go-jose.v2 v2.4.1
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.17.3
	k8s.io/apiextensions-apiserver v0.0.0-20191204090421-cd61debedab5
	k8s.io/apimachinery v0.17.3
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/heapster v1.5.4
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20190313235455-40a48860b5ab
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190409022649-727a075fdec8
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go => k8s.io/client-go v11.0.0+incompatible
	k8s.io/utils => k8s.io/utils v0.0.0-20190809000727-6c36bc71fc4a
)
