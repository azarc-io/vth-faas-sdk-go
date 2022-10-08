# vth-faas-sdk-go
Golang SDK for Verathread Sparks/Connectors

value, err := ctx.GetStage("name").GetValue().Value() 

brew install gh jq
gh repo list azarc-io --limit 9999 --json sshUrl | jq '.[]|.sshUrl' | xargs -n1 git clone

brew install grpcurl
grpcurl -plaintext -vv localhost:7777 sdk_v1.AgentService/ExecuteJob

grpc server config 
```go
svr := grpc.NewServer(grpc.ConnectionTimeout(time.Second * 10))

sdk_v1.RegisterAgentServiceServer(svr, AgentService{})
reflection.Register(svr) // <<<<<<<<<<<

listener, err := net.Listen("tcp", "localhost:7777")
```

```shell
~/dev/code 
✦11 ❯ grpcurl -plaintext -vv localhost:7777 list                          
grpc.reflection.v1alpha.ServerReflection
sdk_v1.AgentService


~/dev/code
✦11 ❯ grpcurl -plaintext -vv -d '{"key": "job_key", "transaction_id": "transaction_id", "correlation_id": "correlation_id" }' localhost:7777 sdk_v1.AgentService/ExecuteJob

Resolved method descriptor:
rpc ExecuteJob ( .sdk_v1.ExecuteJobRequest ) returns ( .sdk_v1.Void );

Request metadata to send:
(empty)

Response headers received:
content-type: application/grpc

Estimated response size: 0 bytes

Response contents:
{
  
}

Response trailers received:
(empty)
Sent 1 request and received 1 response
```