# mailman

Microservice for sending emails

## Requirements

- Go 1.18
- Docker
- Minikube

## Quickstart

#### Build mailman image

```shell
make
```

#### Start minikube and setup everything

```shell
make minikube
```

#### Start minikube tunnelling to access cluster services from local machine

```shell
minikube tunnel
# Or 
make minikube-tunnel
```

#### Stop minikube and clean up everything

```shell
make minikube-clean
```

#### Build and deploy mailman with configuration, services, etc.

```shell
make mailman
```

#### Rebuild and redeploy mailman, and update the configuration

```shell
make mailman-rebuild
```

#### Clean mailman and its configuration, services, etc.

```shell
make mailman-clean
```

#### Deploy postgres DB with configuration, services, etc.

```shell
make postgres
```

#### Reset and redeploy postgres DB (start from scratch)

```shell
make postgres-reset
```

#### Clean postgres DB and its configuration, services, etc.

```shell
make postgres-clean
```

## Sample requests

#### Create a mailing entry

```shell
curl localhost:8080/api/messages -X POST -d '{"email":"jan.kowalski@example.com","title":"Interview","content":"simple text","mailing_id":2, "insert_time": "2022-03-30T15:42:38.72512917Z"}'
# {"id":23}
```

#### Delete a mailing entry

```shell
curl localhost:8080/api/messages/23 -X DELETE
```

#### Send mailing entries with a mailing ID

```shell
curl localhost:8080/api/messages/send -X POST -d '{"mailing_id": 2}'
```
