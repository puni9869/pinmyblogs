# pinmyblogs

[![Pipeline](https://github.com/puni9869/pinmyblogs/actions/workflows/go.yml/badge.svg)](https://github.com/puni9869/pinmyblogs/actions/workflows/go.yml)

### How to start
1. You need latest Go runtime. https://go.dev
2. `export ENVIRONMENT=local`
3. Run `make watch server`. This is a live hot reloading server config. Uses https://github.com/air-verse/air
4. For unit test run `make test`.
5. For vulnerability check run `make govulncheck  && make vet`.
6. To check static errors run `make lint`.


### Create a super user in posgtres:

> CREATE USER <username> WITH SUPERUSER PASSWORD '';

> CREATE USER puni9869 SUPERUSER;

### How to contact
1. Create an issue and raise your concern, bug and feature request. I will be happy to address run-time.
2. Reachout to me on [Linkedin](https://www.linkedin.com/in/punitinani1/)
