module github.com/jay16213/radio_simulator

go 1.14

require (
	git.cs.nctu.edu.tw/calee/sctp v1.0.0
	github.com/free5gc/MongoDBLibrary v1.0.0
	github.com/free5gc/UeauCommon v1.0.0
	github.com/free5gc/aper v1.0.0
	github.com/free5gc/milenage v1.0.0
	github.com/free5gc/nas v1.0.0
	github.com/free5gc/ngap v1.0.1
	github.com/free5gc/openapi v1.0.0
	github.com/golang/protobuf v1.4.3
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.1.3
	github.com/stretchr/testify v1.6.1
	github.com/urfave/cli/v2 v2.3.0
	go.mongodb.org/mongo-driver v1.4.4
	google.golang.org/grpc v1.35.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	git.cs.nctu.edu.tw/calee/sctp => /home/jay/thesis/sctp
	github.com/free5gc/MongoDBLibrary => /home/jay/thesis/MongoDBLibrary
	github.com/free5gc/nas v1.0.0 => github.com/jay16213/nas v1.0.1
)
