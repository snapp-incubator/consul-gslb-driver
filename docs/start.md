# Docs

## Test gRPC server


Do not implement Reflection. It has some known issues with gogo-protobuf, and there are many dependencies to that. Use proto file instead.

Sample test:

CreateGSLB:

```bash
grpcurl -proto gslbi.proto -d '{"name": "test-grpc-service", "service_name": "learn", "host":"spcld-health2-be.apps.private.teh-1.snappcloud.io", "weight":"1", "parameters": {"probe_timeout":"3", "probe_scheme":"http","probe_address":"spcld-health2-be.apps.private.teh-1.snappcloud.io","probe_interval":"5"} }' -plaintext -unix /users/my/gitlab/consul-gslb-driver/socket gslbi.v1.Controller.CreateGSLB
```

DeleteGSLB:

```bash
grpcurl -proto gslbi.proto -d '{"gslb_id": "test-grpc-service"}' -plaintext -unix /users/my/gitlab/consul-gslb-driver/socket gslbi.v1.Controller.DeleteGSLB
```
