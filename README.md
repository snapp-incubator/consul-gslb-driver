# consul-gslb-driver

**This project is under active development and not usable yet**

Gslb driver for Hashicorp Consul.


## Development
* `make test` Run tests.
* `make build` builds golang app locally.
* `make docker-build` build docker image locally.
* `make docker-push` push container image to registry.


* `make build` builds Docker image locally.
* `make run` spins up local Docker containers so you can access via web or command line.
* `make test` runs all tests locally.
* `make rsh` spins up a temporary container based on the built image and gives you a bash inside that container.
* `make debug` spins up a temporary container based on the build image with sleep entrypoint, and gives you a bash inside that container to debug in case container stops immediately.

## Security

### Reporting security vulnerabilities

If you find a security vulnerability or any security related issues, please DO NOT file a public issue, instead send your report privately to cloud@snapp.cab. Security reports are greatly appreciated and we will publicly thank you for it.

## License

GNU Affero General Public License v3.0, see [LICENSE](LICENSE).
