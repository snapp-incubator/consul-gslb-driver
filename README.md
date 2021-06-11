# Consul GSLB Driver

**This project is under active development and not usable yet**

Gslb driver for Hashicorp Consul.


## Development
* `make test` Run tests.
* `make build` builds golang app locally.
* `make docker-build` build docker image locally.
* `make docker-push` push container image to registry.

* `make run` spins up local Docker containers so you can access via web or command line.
* `make test` runs all tests locally.
* `make rsh` spins up a temporary container based on the built image and gives you a bash inside that container.
* `make debug` spins up a temporary container based on the build image with sleep entrypoint, and gives you a bash inside that container to debug in case container stops immediately.


## Usage

Fill a `config.yaml` file such as the one in `config.example.yaml`. You can also set ENV vars corresponding to those values. Then run with:

```bash
consul-gslb-driver -config=<path-to-config> <extra-flags>
```

Note that the order will be:

1. flags
2. env vars: `CONSULGSLB_*`
3. config: in current dir, and then in homedir
4. default cobra Pflag values
5. default value in newConfig()

## Security

### Reporting security vulnerabilities

If you find a security vulnerability or any security related issues, please DO NOT file a public issue, instead send your report privately to cloud@snapp.cab. Security reports are greatly appreciated and we will publicly thank you for it.

## License

Apache-2.0 License, see [LICENSE](LICENSE).
