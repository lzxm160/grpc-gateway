module github.com/tronprotocol/grpc-gateway

go 1.14

require (
	github.com/ethereum/go-ethereum v1.9.20
	github.com/go-resty/resty/v2 v2.3.1-0.20200915215012-608c8d777d0e
	github.com/fbsobreira/gotron-sdk v0.0.0-20200910163704-5dae825f6e2e
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.4.2
	github.com/grpc-ecosystem/grpc-gateway v1.14.8
	github.com/i9/bar v0.0.0-20191101181816-5c944ef12f32
	github.com/stretchr/testify v1.6.2-0.20200818115829-54d05a4e1844
	golang.org/x/net v0.0.0-20200625001655-4c5254603344
	google.golang.org/genproto v0.0.0-20200825200019-8632dd797987
	google.golang.org/grpc v1.31.1
	google.golang.org/protobuf v1.25.0
)

replace github.com/fbsobreira/gotron-sdk => github.com/lzxm160/gotron-sdk v1.0.24
