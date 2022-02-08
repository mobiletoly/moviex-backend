# Table of Contents
1. [Introduction](#introduction)
2. [Pre-requisites](#pre-requisites)
3. [Hexagonal architecture](#hexagonal-architecture)
4. [Project structure](#project-structure)
5. [Running locally](#running-locally)
    1. [Database docker](#database-docker)
    2. [IntelliJ IDEA](#intellij-idea)
    3. [VS Code](#vs-code)
    4. [Run as command-line applications](#run-as-command-line-applications)
6. [Build and deploy to kubernetes](#build-and-deploy-to-kubernetes)
    1. [Database helm](#database-helm)
    2. [Microservices deployment](#microservices-deployment)
7. [Access API Gateway](#access-api-gateway)

# Introduction

Moviex is an example of Go based microservice code and infrastructure that you can use to bootstrap your own project.
To provide a very basic functionality I decided to build a relatively simple Moviex application that allows
external clients to fetch movie descriptions and perform some basic queries.
We use a very minimalistic approach with small number of dependencies. This README is structured as a tutorial so
please follow the flow. Not only it explains a code structure, but also how to build docker images and deploy your
code to kubernetes.

Here is a high-level overview of what you can expect to see in this code.

- Example of API Gateway project that acts as a publicly facing microservice. It provides GraphQL interface
(we use great `gqlgen` library to generate Go code from graphql schemas) to communicate with external clients.
API Gateway's only job is to properly redirect requests to other microservices and handle and federate responses
received from business-logic microservices. Note that API Gateway uses GraphQL only to interact with the outside
world, but to communicate with other Moviex microservices - it uses gRPC (with aim to great performance).
If you  don't care about GraphQL - it will be very easy for you to swap GraphQL with REST, since we already run
net/http server to handle HTTP requests.


- In addition to API Gateway service we have two simple business-logic microservices - Film Service and User Service.
Film Service provides access to film and actor database, while User Service keeps user login information
and simple list of purchased films.


- We use PostgreSQL database to store films and users. `sqlx` library for Go is used for this. Sample data are
provided and bootstrapped if needed when we deploy PostgreSQL.


- We use a hexagonal architecture (ports and adapters) to structure our app. I strongly encourage developers to
read about hexagonal architecture, it really helps to create modular projects. But keep in mind that in Go
we try to avoid unnecessary levels of abstraction, and it is reflected in Moviex code - hexagonal architecture
is pretty "lightweight" here.


- Moviex supports Request Id provided by caller of API Gateway via HTTP header or auto-generated in case it is not
provided by caller. It is passed all the way from API Gateway to other microservices to ensure that we can properly
trace code flow and errors for each request separately. Logger that we use is set up to print request id
that is shared amongst multiple microservices. It is simple but efficient way to troubleshoot distributed calls.


- Very simple dependency injection patterns are used. We don't use any libraries for this, just structure
our code in a way that it is easier to wire and pass dependencies around.


- We provide docker files to create service images as well as Helm chart to deploy our code to kubernetes. Everything
is simplified to make it easy to understand, so while steps that we describe are sufficient to deploy and run
code in kubernetes, but you might need to do more tuning to make it production quality and secure (well, it depends
on your requirements).

  
#### Client/Server architecture

```
                              Kubernetes (not necessary)
                           +------------------------------------------------------------------
                           |
client -----(GraphQL)------|--> [API Gateway service] --(gRPC)--> [User service]
                           |       |                                      .
                           |       |                                      .
                           |       |                                      .
                           |       +--(gRPC)--> [Film service] . . . [PostgreSQL database] 
                           |   
```


# Pre-requisites

- You need Go tools to be available on your computer: https://go.dev/doc/install


- Moviex code requires **protobuf** compiler to generate Go files from protobuf schemas. Follow this link to find
out how to install protobuf on your computer: https://grpc.io/docs/languages/go/quickstart/ <br>
If you are on Mac and have brew installed - you can easily do it by entering these commands:
    ```shell
    $ brew install protobuf
    $ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    $ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    ```

- We use awesome **gqlgen** tool for Go to generate stub files and models for GraphQL schemas.
  You can become familiar with it by visiting: https://github.com/99designs/gqlgen


- **docker** is required to deploy your app. Also, you want to make sure that you have docker
  compose tool available if you are planning to use our docker-compose.yml files (e.g. to
  run a PostgreSQL database with initial data locally)


- **minikube** is required if you plan to deploy this app to kubernetes locally. For Mac/brew you can install
it with `$ brew install minikube`. Make sure to run `$ minikube start` to launch a local cluster,


- We use **helm** to simplify deployment of apps to kubernetes


# Hexagonal architecture

Before we dive into details of project structure, let's take a look at the common architectural pattern that
we use in our code. It's called Hexagonal Architecture (or Ports and Adapters) and you can find plenty of
information online. Here is an example of how we apply this architecture to our code:

API Gateway service - the flow from receiving incoming GraphQL requests to sending gRPC request to another
microservice:
```
(primary adapter: /apigateway/internal/adapters/primary/apiserver)
     ↓
   (input use case - business logic: /apigateway/internal/core/usecase)
        ↓
      [ output port interface: /apigateway/internal/core/output ]
      (implementation of port - secondary adapter: /apigateway/internal/adapters/remoting)
```

All platform-specific calls and functionalities (such as calling gRPC services in API Gateway or database
access in Film and User services) reside in folder `adapters/` (e.g. `adapters/remoting/` folder contains
gRPC clients). All directories in `adapters/` called "secondary (or driven)" adapters. Notable exception in
this layout is directory called `adapters/primary` - this is where our "primary (or driver)" adapter resides.
Primary adapter is an entry point for an application, something that receives commands from the outside world and
passes it downstream - it can be a web server or command-line app handler. In our case we have
`adapter/primary/apiserver` - this is our GraphQL-based server that receives GraphQL-based POST requests
from external clients. Primary adapter should not communicate with secondary adapters directly (moreover,
secondary adapters should not communicate with each others directly). Instead, primary adapter should communicate
with business logic (that resides in `core/` directory) via what's called Use Cases (business logic layer)
and use cases will communicate with secondary adapters via Output Ports. These output ports are interfaces
declared in `core/` directory and are implemented by secondary adapters. Secondary adapters usually convert business
logic entities into platform-specific repository entities and call repository functionality after that.

So here is an example of how it looks for our API Gateway:

1. Primary adapter `adapters/primary/apiserver` receives GraphQL request to fetch all films. In our case call
   is processed by function in `adapters/primary/apiserver/graph/film.resolvers.go` (this is how gqlgen library
   works).
2. Handler in `film.resolvers.go` file routes this request to business logic that resides in
   `core/usecase/film.go` (`FetchFilms` function). Usually you can provide more business logic here (e.g.
   to add caching requests), but we don't do too much in API Gateway and simply call remoting service to
   communicate with another service. We do this via calling output port - function `FetchFilms` declared in
   interface `outport.FilmRemoting` that resides in core business logic file `core/outports/remoting.go`.
4. In our case secondary adapter `adapters/remoting` will be called because it has functionality implementing
   FetchFilms port. Normally adapter's job is to convert business logic entities into repository entities. In
   our case the input for adapter will be GraphQL model generated by gqlgen library according to GraphQL
   specification and adapter will convert it to protobuf model that can be consumed by gRPC call. Then adapter
   calls gRPC repository function that suppose to perform remote gRPC call to Films service. Once response is received
   from Films service - then adapter will convert protobuf response message into GraphQL model and return it back
   to business logic layer that called the adapter. In our case Use Case receives a response and returns it
   back to primary adapter (GraphQL server) that sends it back to a caller.


# Project structure

Since all our services are written in Go - we have decided to keep them in one parent directory. It greatly
simplifies a development and debugging and makes sharing common code much easier. Even if code is in a single
directory, the build process still generates separate executable files for microservices, separate Dockerfiles
are used etc. Note, that we heavily use `internal` directories. This is to separate source code from different
services and make it isolated from each other.

Let's examine project structure to simplify our navigation.

Overview:

`. cmd/`  - command-line apps for services<br>
`.  | apigateway/` - command-line app to launch API Gateway service<br>
`.  | filmsrv/` - command-line app to launch Film service<br>
`.  | usersrv/` - command-line app to launch User service<br>

`. configs/` - configuration files<br>
`.  | apigateway/` - config files for API Gateway service<br>
`.  |   | config-local.yaml` - service configuration for local deployment<br>
`.  |   | config-k8s.yaml` - service configuration for kubernetes deployment<br>
`.  | filmsrv/` - config files for Film service<br>
`.  |   | ...`<br>
`.  | usersrv/` - config files for User service<br>
`.  |   | ...`<br>

`. deploy/` - deployment tools<br>
`.  | docker/` - docker files<br>
`.  |   | apigateway` - contains Dockerfile for API Gateway<br>
`.  |   | filmsrv` - contains Dockerfile for Film service<br>
`.  |   | usersrv` - contains Dockerfile for User service<br>
`.  |   | moviex-db` - database docker tools<br>
`.  |   |   | docker-compose.yml` - to run PostgreSQL with docker compose<br>
`.  |   |   | init.sql.gz` - gzipped script to create database with initial data<br>

`. internal/` - main code directory<br>
`.  | common/` - common code in use by more than one service<br>
`.  |  | db/` - common database operations<br>
`.  |  |  | config.go` - database configuration struct<br>
`.  |  |  | connector.go` - connector to PostgreSQL database<br>
`.  | proto/` - protobuf definitions to communicate between microservices<br>
`.  |  | filmsrv.proto` - models and rpc interfaces to access Film service<br>
`.  |  | usersrv.proto` - models and rpc interfaces to access Film service<br>
`.  | apigateway/` - API Gateway service<br>
`.  | filmsrv/` - Film service<br>
`.  | usersrv/` - User service<br>
<br>
**API Gateway** layout (only important files and directories are shown):<br>
`. apigateway/`<br>
`.  | infra/` - entry point (service launcher and dependency initialization)<br>
`.  |  | apiserver.go` - API server launcher<br>
`.  |  | deps.go` - instantiate and wire dependencies<br>
`.  | internal/`<br>
`.  |  | adapters/` - primary and secondary adapters<br>
`.  |  |  | appconfig/` - adapter to access service yaml configuration<br>
`.  |  |  |  | internal/`<br>
`.  |  |  |  |  | mapper` - mappers from yaml configs to core's **app.Config** struct<br>
`.  |  |  |  |  | loader.go` - loads configurations from **/config/apigateway/\*.yaml**<br>
`.  |  |  |  | appconfig.go` - config loader to return **app.Config** struct with service config<br>
`.  |  |  | remoting/` - adapter to perform gRPC calls to other services<br>
`.  |  |  |  | mapper/` - mappers from gRPC response to GraphQL entities<br>
`.  |  |  |  | film.go` - adapter to call Film service via gRPC<br>
`.  |  |  |  | user.go` - adapter to call User service via gRPC<br>
`.  |  |  | primary/`<br>
`.  |  |  |  | apiserver/` - primary adapter to run GraphQL server<br>
`.  |  |  |  |  | generated/` - GraphQL generated server code<br>
`.  |  |  |  |  | graph/` - GraphQL handlers to call use cases from **/core/usecase** dir<br>
`.  |  |  |  |     | film.resolver.go` - GraphQL film handlers<br>
`.  |  |  |  |     | user.resolver.go` - GraphQL user handlers<br>
`.  |  | core/` - core business functionality<br>
`.  |  |  | app/` - service-wide application code<br>
`.  |  |  |  | appconfig.go` - app config structs for adapters and core<br>
`.  |  |  | model/` - models to flow between core business logic and adapters<br>
`.  |  |  |  | models_gen.go` - GraphQL models generated by gqlgen from **/graph/\*.graphqls** files<br>
`.  |  |  |  | pagination.go` - model and utilities for GraphQL pagination**<br>
`.  |  |  | outport/` - output ports to be implemented by secondary adapters<br>
`.  |  |  |  | remoting.go` - output ports to perform remote calls<br>
`.  |  |  | usecase/` - use cases<br>
`.  |  |  |  | film.go` - film use cases to orchestrate access to remoting output ports<br>
`.  |  |  |  | user.go` - user use cases to orchestrate access to remoting output ports<br>
`.  |  | di/`<br>
`.  |  |  | di.go` - holds dependencies<br>
`.  |  | graph/` - place to keep GraphQL schema definitions and gqlgen settings<br>
`.  |  |  | gqlgen.yml` - gqlgen generator settings<br>
`.  |  |  | film.graphqls` - GraphQL schema definitions to access films<br>
`.  |  |  | user.graphqls` - GraphQL schema definitions to access users<br>
`.  |  |  | root.graphqls` - common definitions<br>
<br>
**Film service** layout (only important files and directories are shown):<br>
`. filmsrv/`<br>
`.  | infra/` - entry point (service launcher and dependency initialization)<br>
`.  |  | apiserver.go` - API server launcher<br>
`.  |  | deps.go` - instantiate and wire dependencies<br>
`.  | internal/`<br>
`.  |  | adapters/` - primary and secondary adapters<br>
`.  |  |  | appconfig/` - adapter to access service yaml configuration<br>
`.  |  |  |  | internal/`<br>
`.  |  |  |  |  | mapper` - mappers from yaml configs to core's **app.Config** struct<br>
`.  |  |  |  |  | loader.go` - loads configurations from **/config/filmsrv/\*.yaml**<br>
`.  |  |  |  | appconfig.go` - config loader to return **app.Config** struct with service config<br>
`.  |  |  | persist/` - adapter to perform database queries<br>
`.  |  |  |  | internal/`<br>
`.  |  |  |  |  | mapper/` - mappers from database entities to protobuf models<br>
`.  |  |  |  |  | repo/` - database repository to access PostgreSQL database<br>
`.  |  |  |  |     | actor.go` - repository to query actor table<br>
`.  |  |  |  |     | category.go` - repository to query category table<br>
`.  |  |  |  |     | film.go` - repository to query film table<br>
`.  |  |  |  |     | repo.go` - common repository functions<br>
`.  |  |  |  | dbconn.go` - database connector<br>
`.  |  |  |  | film.go` - adapter to call database repository<br>
`.  |  |  | primary/`<br>
`.  |  |  |  | apiserver/` - primary adapter to run GraphQL server<br>
`.  |  |  |  |  | apiserver.go` - gRPC server runner<br>
`.  |  | core/` - core business functionality<br>
`.  |  |  | app/` - service-wide application code<br>
`.  |  |  |  | appconfig.go` - app config structs for adapters and core<br>
`.  |  |  | outport/` - output ports to be implemented by secondary adapters<br>
`.  |  |  |  | persist.go` - output ports to perform database calls<br>
`.  |  |  | usecase/` - use cases<br>
`.  |  |  |  | film.go` - film use cases to orchestrate access to database output ports<br>
`.  |  | di/`<br>
`.  |  |  | di.go` - holds dependencies<br>
<br>

**User service** layout is very similar to **Film service** layout, because both of them are gRPC
servers and use very similar architecture. we ski

# Running locally

## Database docker

First step is to run a PostgreSQL database and import data that we have in `/deploy/docker/moviex-db/init.sql.gz` file.
This file contains dump of data with a sample of films, actors, categories. If you already have your own
database instance running - you can restore this backup to your own instance, but probably the easiest
approach is to use a docker-compose we have supplied. In terminal window enter:

```shell
$ docker-compose -f deploy/docker/moviex-db/docker-compose.yml up
```

to run PostgreSQL with default login "postgres" and password "postgres" on port 5432. also
`init.sql.gz` will be executed and `moviex` database with some initial data will be
created.

Note that while we have a single database for a multiple microservices, but our database has few
schemas and each service has access to its own database schema. There are some other approaches possible, e.g.
each microservice will be using its own database. There are config files where you can
find service options, including database options. For example, Film Service has its configuration in
*/config/filmsrv/config-local.yaml*/<br> for local deployment. By default, it looks like:

```
database:
    host: 127.0.0.1
    port: 5432
    name: moviex
    user: _
    password: _
    sslmode: disable
```

Note that `user` and `password` fields are empty (`_`). You can put your username and password here,
but usually it is not a good idea. Instead, the configuration loader in our code allows overriding database
parameters (or any other parameters that can be added to config file) by using environment variables with
syntax such as `[SERVICENAME]_[CONFIGNODE]_[CONFIGNAME]`. To change parameters for Film service we use a
name `FILMSRV` (it will be `USERSRV` for User service), the CONFIGNODE part will be `DATABASE` and password
will be `PASSWORD`. So all together to pass database password to Film service, you have to set up environment
variable `FILMSRV_DATABASE_PASSWORD=mypassword`.

## IntelliJ IDEA

In IntelliJ open moviex-backend directory. Choose **Run** / **Edit Configuration** in the menu and add few
new **Go Build** tasks:

- Name: `apigateway`

  Package: `github.com/mobiletoly/moviex-backend/cmd/apigateway`

  Program arguments: `server --deployment=local`


- Name: `filmsrv`

  Package: `github.com/mobiletoly/moviex-backend/cmd/filmsrv`

  Program arguments: `server --deployment=local`

  Environment: `FILMSRV_DATABASE_USER=postgres;FILMSRV_DATABASE_PASSWORD=postgres;`


- Name: `usersrv`

  Package: `github.com/mobiletoly/moviex-backend/cmd/usersrv`

  Program arguments: `server --deployment=local`

  Environment: `USERSRV_DATABASE_USER=postgres;USERSRV_DATABASE_PASSWORD=postgres;`


As we already mentioned, using a special conventions for environment variables, you can override
default settings in `/configs/[service]/config-[env].yaml` files and this is exactly what we
did here.

Now you can run all three services one by one (order does not matter) for Run or Debug.

**Tip:** if you want to run all three apps at the same time, you can open Run / Edit Configuration,
select new **Compound** configuration and add all three previously added configs into it. When you run or debug
Compound configuration - it runs all three services one by one.

## VS Code

(documentation is not available yet, but it is coming...)

## Run as command-line applications

If you want to run code from command line, you can build all three apps:

```
$ go build -ldflags="-w" -o .  ./cmd/...
```

As a result three command-line executables will be created: `apigateway`, `filmsrv`, `usersrv`.
To run API Gateway, enter this command:
```
$ ./apigateway server --deployment=local
```
To run Film and User service:
```
$ FILMSRV_DATABASE_USER=postgres FILMSRV_DATABASE_PASSWORD=postgres ./filmsrv server --deployment=local
$ USERSRV_DATABASE_USER=postgres USERSRV_DATABASE_PASSWORD=postgres ./usersrv server --deployment=local
```

# Build and deploy to kubernetes

For local kubernetes deployment we will be using minikube. Please read up on what is minikube and how to install 
and run it. If you want to deploy code to a cloud - you are on your own (as of now).

We decided not to add all dependencies into a single helm deployment, instead split process into two steps -
deployment of PostgreSQL and deployment of all three microservices - API gateway, Film Service, and User Service.

## Database helm

Let's start with PostgreSQL (we will be using bitnami deployment). Note that this is not going to be a production
quality set up, and it will up to you to have a properly configured PostgreSQL instances. To do so you have
multiple options: properly setup Bitnami PostgreSQL cluster
(https://engineering.bitnami.com/articles/create-a-production-ready-postgresql-cluster-bitnami-kubernetes-and-helm.html),
use AWS RDS etc.

We start by adding bitnami PostgreSQL repository:

```shell
$ helm repo add bitnami https://charts.bitnami.com/bitnami
```

Before we proceed with installing postgresql pod - let's copy our initial script with SQL commands into kubernetes
config map. This will provide us with some initial database table structures and data, such as films and actors.
Since we already have script in `/deploy/docker/moviex-db` directory - we will use it:

```shell
$ kubectl create configmap db-init-schema --from-file=deploy/docker/moviex-db/init.sql.gz
```

Here we uploaded SQL script into config map with key name `db-init-schema`. You can check that operation was
successful by running `kubectl get configmaps`.

Next let's install and run PostgreSQL with helm chart:

```shell
$ helm install my-postgresql bitnami/postgresql \
    --set postgresqlUsername=postgres \
    --set postgresqlDatabase=moviex \
    --set postgresqlPostgresPassword=postgres \
    --set postgresqlUsername=postgres \
    --set postgresqlPassword=postgres \
    --set initdbScriptsConfigMap=db-init-schema
```

Here we used simple password, you don't necessary want to do this for production.
Also note `--set initdbScriptsConfigMap=db-init-schema` parameter - this instructs postgres to run SQL script
we have previously uploaded to configuration map.

For a convenience you can run these commands: 
```
$ export POSTGRES_PASSWORD=$(kubectl get secret --namespace default my-postgresql -o \
    jsonpath="{.data.postgresql-password}" | base64 --decode)
$ kubectl run my-postgresql-client --rm --tty -i --restart='Never' --namespace default \
    --image docker.io/bitnami/postgresql:11.14.0-debian-10-r28 \
    --env="PGPASSWORD=$POSTGRES_PASSWORD" \
    --command -- psql --host my-postgresql -U postgres -d moviex -p 5432
```

First command extracts PostgreSQL password from running my-postgresql pod and copies it into `POSTGRES_PASSWORD`
environment variable. In our case it will `postgres`, but without passing `--set postgresqlPassword=postgres` 
random password will be generated and then `POSTGRES_PASSWORD` variable will contain newly generated password.
Second command creates a new kubernetes pod based on bitnami with postgres tools installed and logs you
into `psql` on this pod. Once you logged into psql - you will be able to execute commands on database running
in my-postgresql pod. For example, enter `\l` command to see the list of available database, if everything is OK,
you should be able to see `moviex` in a list. You can also try to enter `\dt *.*` to see all tables available
(including postgres system tables).

Now when everything is up and running, you can execute 
```shell
$ kubectl get all -o wide
``` 
and the output should be something similar to

```
NAME                             READY   STATUS    RESTARTS   AGE
pod/my-postgresql-postgresql-0   1/1     Running   0          39m

NAME                             TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)    AGE
service/kubernetes               ClusterIP   10.96.0.1      <none>        443/TCP    30d
service/my-postgresql            ClusterIP   10.96.58.149   <none>        5432/TCP   39m
service/my-postgresql-headless   ClusterIP   None           <none>        5432/TCP   39m

NAME                                        READY   AGE
statefulset.apps/my-postgresql-postgresql   1/1     39m
```

This is pretty good. We already have a my-postgresql service running as well as Persistent Volume, and a StatefulSet
installed. Sure, we have only one pod running for database, which could be OK or not OK for your production
environment, but as I mentioned in the beginning of this section - there are alternatives.

**Tip:** If by some reasons while accessing my-postgresql you are seeing authentication error such as
`password authentication failed for user "postgres"` - then there is a chance that you have previously tried
to deploy postgres with different credentials and volume was not properly removed. You can try to uninstall it
fully and delete all leftovers:
```
$ helm uninstall my-postgresql
$ kubectl delete pvc -l app.kubernetes.io/name=postgresql
```
After that you have to redeploy PostgreSQL.

## Microservices deployment

Now we are ready to deploy all 3 services (API Gateway, Film service, User service) to kubernetes. For this
we have a Helm chart in `deploy/helm` directory. Will start with building docker images for our services (don't
forget to replace `DOCKER_USERNAME` with your actual username registered in Docker Hub), then push it
to Docker Hub

```shell
$ docker build -f deploy/docker/apigateway/Dockerfile -t DOCKER_USERNAME/apigateway:latest .
$ docker push ptolik/apigateway
$ docker build -f deploy/docker/filmsrv/Dockerfile -t DOCKER_USERNAME/filmsrv:latest .
$ docker push ptolik/filmsrv
$ docker build -f deploy/docker/usersrv/Dockerfile -t DOCKER_USERNAME/usersrv:latest .
$ docker push ptolik/usersrv
```

For now, we will be dealing with `latest` tag instead of versions, makes it easier for us to develop and test.

**TIP:** if you want to avoid pushing images to Docker Hub and instead want to keep docker image locally
(it is faster process and simplify redeployments), you can run `eval $(minikube -p minikube docker-env)` in your
active terminal window. Now you can start using `docker build` without `docker push`. Remember to enter this command
for every terminal session that you use (or when you open a new one).

It is time to run our helm chart and deploy our microservices to kubernetes cluster:

```shell
$ helm install moviex ./deploy/helm/ --set database.user=postgres --set database.password=postgres
```

This is it, services are deployed and ready to be used. If we run `kubectl get all` we should see something
similar to:

```shell
NAME                                             READY   STATUS    RESTARTS   AGE
pod/moviex-backend-apigateway-77f4d994c7-84pnk   1/1     Running   0          54m
pod/moviex-backend-apigateway-77f4d994c7-mlmnr   1/1     Running   0          54m
pod/moviex-backend-filmsrv-67f566c798-4v5xx      1/1     Running   0          54m
pod/moviex-backend-filmsrv-67f566c798-l24h6      1/1     Running   0          54m
pod/moviex-backend-usersrv-5749fc4cdd-9vw2f      1/1     Running   0          54m
pod/moviex-backend-usersrv-5749fc4cdd-xjlnp      1/1     Running   0          54m
pod/my-postgresql-postgresql-0                   1/1     Running   0          56m

NAME                                TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
service/kubernetes                  ClusterIP   10.96.0.1        <none>        443/TCP          63m
service/moviex-backend-apigateway   NodePort    10.101.133.166   <none>        8080:31849/TCP   54m
service/moviex-backend-filmsrv      ClusterIP   10.96.49.104     <none>        8080/TCP         54m
service/moviex-backend-usersrv      ClusterIP   10.97.26.152     <none>        8080/TCP         54m
service/my-postgresql               ClusterIP   10.100.221.77    <none>        5432/TCP         56m
service/my-postgresql-headless      ClusterIP   None             <none>        5432/TCP         56m

NAME                                        READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/moviex-backend-apigateway   2/2     2            2           54m
deployment.apps/moviex-backend-filmsrv      2/2     2            2           54m
deployment.apps/moviex-backend-usersrv      2/2     2            2           54m

NAME                                                   DESIRED   CURRENT   READY   AGE
replicaset.apps/moviex-backend-apigateway-77f4d994c7   2         2         2       54m
replicaset.apps/moviex-backend-filmsrv-67f566c798      2         2         2       54m
replicaset.apps/moviex-backend-usersrv-5749fc4cdd      2         2         2       54m

NAME                                        READY   AGE
statefulset.apps/my-postgresql-postgresql   1/1     56m
```

As you can see - we have 2 pods per each microservice. Each set of microservices sits behind service. E.g. to call
Film service from API Gateway - you will be using `http://moviex-backend-filmsrv:8080` endpoint and call will be routed
to one of the `moviex-backend-filmsrv-...` pods. Our configuration file `/configs/apigateway/configs-k8s.yaml` already
contains settings with proper host name and port:
```yaml
services:
  filmsrv:
    host: moviex-backend-filmsrv
    port: 8080
  usersrv:
    ...
```

This configuration allows API Gateway to communicate with Film and User services.

Now we can try to verify connection to API Gateway from within a cluster. Let's create busybox pod inside our
kubernetes cluster and perform `/version` HTTP GET request to our API Gateway:

```shell
$ kubectl run access-client --rm --tty -i --restart='Never' --namespace default --image busybox --command -- sh
```
Once pod is created and we are inside:
```shell
# wget -qO- http://apigateway-moviex-backend:8080/version
# cat version
```
it should print API Gateway service version (we have a very simple HTTP GET handler in our API Gateway code
to return version and also in our helm chart we use it for liveness and readiness probe).

# Access API Gateway

Next step is to send GraphQL request to API Gateway service. If you have launched API gateway via IntelliJ or
from command line, then your request URL will be `http://127.0.0.1:8080/query`. If you use minikube and have
microservices already deployed to kubernetes, then one of the solutions could be to expose moviex-backend-apigateway
service to our localhost, e.g.
```shell
$ minikube service --url moviex-backend-apigateway
```
The output will be something like `http://127.0.0.1:52583` (some other random port will be used in your case).
Let's try to perform GraphQL request to API Gateway to test it (this will also test gRPC interaction of API Gateway 
with Film service).

```shell
$ curl --location --request POST 'http://127.0.0.1:52583/query' \
--header 'Content-Type: application/json' \
--data-raw '{"query":"query films($first: Int!, $after: String) {films(first: $first, after: $after) { totalCount pageInfo { hasNextPage endCursor } edges { node { id title category { id name lastUpdate } actors { id firstName lastName lastUpdate } } }}}","variables":{"first":10,"after":null}}'
```

If everything is OK, you will see JSON response with 10 films. cURL is probably not the best tool to author
GraphQL requests, you can try some UI tools instead (I personally like Insomnia).
