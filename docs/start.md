# Docs

## Test gRPC server


Do not implement Reflection. It has some known issues with gogo-protobuf, and there are many dependencies to that. Use proto file instead.

Sample test:

```bash
grpcurl -proto gslbi.proto -d '{"serviceID": "hi123"}' -plaintext -unix /users/my/gitlab/consul-gslb-driver/socket gslbi.v1.Controller.DeleteGSLB
```
