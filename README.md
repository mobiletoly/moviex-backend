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
7. 

# Introduction

Moviex is a template of Go based microservice infrastructure that you can use to start your own project.
To provide a very basic functionality I decided to build a relatively simple Moviex application that allows
external clients to fetch for film descriptions and perform some basic queries.
I try not to use any fancy libraries and use a very minimalistic approach with as smallest number of dependencies
as possible. This README is structured as a tutorial way so please follow the flow. Not only it explains a code
structure, but also how to deploy your code to kubernetes.

Here is a high level overview of what you can expect to see in this code.

- Example of API Gateway microservice that acts as a publicly facing microservice. It provides GraphQL interface
(I use a great gqlgen library for this) to communicate with external clients. API Gateway is a thin
service and its only job is to properly redirect requests to other microservices and properly handle and
federate responses received from business-logic microservices. Note that while API Gateway uses GraphQL to
interface with the outside world, but to communicate with other Moviex microservices - it uses gRPC.
- Also we have two simple business-logic microservices - Film Service and User Service. Film Service takes care of 
providing access to films and actors database, while User Service keeps user login information and very simple film
list of purchased movies.
- We use PostgreSQL database to store films and users. Go's sqlx library is used for this. Sample database is
provided and bootstrapped if needed when docker is launched.
- We use a hexagonal architecture (ports and adapters) to structure our app. I strongly encourage developer to
read up about hexagonal architecture, env if I tried to simplify it for our Go project to make sure that we don't
have any unnecessary abstractions and have a reasonable number of layers and data entities.
- Moviex supports Request Id provided by caller to API Gateway via HTTP header or auto-generated in case if not
provided by caller. It is passed all the way from API Gateway to other microservices to ensure that we can properly
track code flow, errors etc separately per each request. logrus logger that we use is set up to print request id
that is shared amongst multiple microservices. It is simple but efficient way to troubleshoot distributed calls.
- Very simple dependency injection patterns are used. We don't use any special libraries for this, just structure
our code in a way that it is easier to wire and pass dependencies around.
<br><br>
- Kubernetes setup to deploy application to kubernetes cluster

# Pre-requisites

Obviously you need Go tools to be available on your computer: https://go.dev/doc/install

- Out Moviex code requires **protobuf** compiler to generate Go files from protobuf schemas. Follow this link to find out how
to install protobuf on your computer: https://grpc.io/docs/languages/go/quickstart/ <br>
If you are on Mac and have brew install you can easily do it by entering this commands:
    ```
    $ brew install protobuf
    $ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    $ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    ```

- We use awesome **gqlgen** tool for Go to generate stub files and models for GraphQL schemas.
  You can become familiar with it by visiting: https://github.com/99designs/gqlgen


- **docker** is required to deploy your app. also you want to make sure that you have docker
  compose tool available if you are planning to use our docker-compose.yml files (e.g. to
  run a PostgreSQL database with initial data locally)


- **minikube** is required if you plan to deploy this app to kubernetes locally. For Mac you can install it with


- We use **helm** to simplify deployment of apps in kubernetes

# Hexagonal architecture

Before we dive into details of project structure, let's take a look at the common architectural pattern that
we use in our code. It's called Hexagonal Architecture (or Ports and Adapters) and you can find plenty of
information online. Here is an example of how we apply this architecture to our code:

API Gateway service - the flow from receiving incoming GraphQL requests to sending gRPC request to another
microservice:
```
(primary adapter: /apigateway/internal/adapters/primary/apiserver)
     ↓
   (input use case: /apigateway/internal/core/usecase)
        ↓
      [ output port interface: /apigateway/internal/core/output ]
      (implementation of port, secondary adapter: /apigateway/internal/adapters/remoting)
```

All platform specific calls and functionality (such as calling gRPC services in API Gateway or database
interaction in Film and User services) reside in folder `adapters/` (e.g. `adapters/remoting/` folder contains
gRPC clients). All the directories in `adapters/` called "secondary (or driven) adapters". Notable exception in
this layout is directory called `adapters/primary` - this is where our "primary (or driver) adapter" resides.
Primary adapter is an entry point for application, something that receives commands from the outside world and
passes it down-streams - it can be a web server or command-line app handler. In our case we have
`adapter/primary/apiserver` - this is our GraphQL-based server that can receive GraphQL-based POST requests
from external clients. Primary adapter should not communicate with secondary adapters directly (more over,
even secondary adapters should not communicate to each other). Instead, primary adapter should communicate with
business logic (that resides in `core/` directory) via what's called Use Cases (business logic layer)
and use cases will communicate with secondary adapters via Output Ports. This output ports are interfaces and
implemented by secondary adapters. Secondary adapters usually convert business logic entities into repository
entities and then call repository functionality.

So in a nutshell here is how it looks like

1. Primary adapter `adapters/primary/apiserver` receives GraphQL request to fetch all films. In our case call
   is processed by function in `adapters/primary/apiserver/graph/film.resolvers.go` (this is how gqlgen library
   works).
2. API server (specifically function in `film.resolvers.go` file) routes this request to business logic that
   resides in `core/usecase/film.go` (`FetchFilms` function). Usually you can have some business logic here, but
   we don't do too much in API Gateway and simply call remoting service to communicate with another service.
   We do this via calling output port - function `FetchFilms` declared in interface `outport.FilmRemoting` that
   resides in core business logic file `core/outports/remoting.go`.
3. In our case secondary adapter `adapters/remoting` will be called because it has functionality implementing
   FetchFilms port. Normally adapter's job is to convert business logic entities into repository entities. In
   our case the input for adapter will be GraphQL models generated by gqlgen library according to GraphQL
   specification and adapter will convert it to protobuf models that can be consumed by gRPC calls. Then adapter
   calls gRPC repository function that suppose to perform remote call to Films service. Once response is received
   from Films service - then adapter will convert protobuf response message into GraphQL model and return it back
   to business logic layer that called the adapter. In our case Use Case receives a response and returns it
   back to primary adapter (GraphQL server) that sends it back to a caller.

# Project structure

Since all services are written in Go - we have decided to keep them in one parent directory. It greatly simplifies
a development and debugging and makes sharing common code much easier. While code is in a single directory, the
build process still generates separate executable files for microservices, separate Dockerfiles are used etc.
Note, that we heavily use `internal` directories. This is to separate source code from different service
and make it inaccessible to each other.

Let's examine project structure to simplify our navigation.

Overview:

`. cmd/`  - command-line apps for services<br>
`.  | apigateway/` - command-line app to launch API Gateway service<br>
`.  | filmsrv/` - command-line app to launch Film service<br>
`.  | usersrv/` - command-line app to launch User service<br>
`. configs/` - configuration files<br>
`.  | apigateway/` - config files for API Gateway service<br>
`.  | filmsrv/` - config files for Film service<br>
`.  | usersrv/` - config files for User service<br>
`. internal/` - main code directory<br>
`.  | common/` - common code in use by more than one service<br>
`.  |  | db/` - common database operations<br>
`.  |  |  | config.go` - database configuration struct<br>
`.  |  |  | connector.go` - connector to PostgreSQL database<br>
`.  | proto/` - protobuf definitions to communicate between microservices<br>
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
servers and use very similar architecture.

# Running locally

## Database docker

First step is to run a PostgreSQL database and import data that we have in `/deploy/docker/moviex-db/init.sql` file.
This backup file contains dump of data with some films, actors, categories. If you already have your own
database instance running - you can restore this backup to your own database, but probably the easiest
approach is to use a docker-compose. In your terminal:

```
$ cd deploy/docker/moviex-db
$ docker-compose up
```

it will run PostgreSQL with default login "postgres" and password "postgres" on port 5432. also
`init.sql` will be automatically executed and `moviex` database with some initial data will be
created.

Note that while we have a single database for a multiple microservices, but our database has few
schemas and each service has access ownly to its own schema. You can easily change this and point
our services to access different databases, all you have to do is to open service config file
and change parameters. For example, for Film Service you should open */config/filmsrv/config-local.yaml*/<br>
By default it is going to look like:

```
database:
    host: 127.0.0.1
    port: 5432
    name: moviex
    user:
    password:
    sslmode: disable
```

Note that `user` and `password` fields are empty. You can put your username and password here,
but usually it is not a good idea. Instead this code allows you to override database parameters
(or any other parameters that can be added in config file) by using environment variables with
syntax such as "SERVICE_CONFIGNODE_CONFIGNAME" in all-capital letters form, e.g. to change 
parameters for Film service we use a name `FILMSRV` (it will be `USERSRV` for User service),
the CONFIGNODE will be `DATABASE` and let's say password will be `PASSWORD`. So all together
to pass database password to Film service, you have to set up environment variable
`FILMSRV_DATABASE_PASSWORD=mypassword`.

## IntelliJ IDEA

Now open moviex-backend directory in IntelliJ and let's setup projects. Choose Run / Edit Configuration
in the menu and add few new **Go Build** tasks:

- Name: `apigateway`

  Package: `github.com/mobiletoly/moviex-backend/cmd/apigateway`<br>

  Program arguments: `server --port=8080`


- Name: `filmsrv`

  Package: `github.com/mobiletoly/moviex-backend/cmd/filmsrv`

  Program arguments: `server --port=8081`

  Environment: `FILMSRV_DATABASE_USER=postgres;FILMSRV_DATABASE_PASSWORD=postgres;FILMSRV_DATABASE_SSLMODE=disable`


- Name: `usersrv`

  Package: `github.com/mobiletoly/moviex-backend/cmd/usersrv`

  Program arguments: `server --port=8082`

  Environment: `USERSRV_DATABASE_USER=postgres;USERSRV_DATABASE_PASSWORD=postgres;USERSRV_DATABASE_SSLMODE=disable`


As we already mentioned, using a special conventions for environment variables, you can override
default settings in `/configs/[service]/config-[env].yaml` files and that is exactly what we
were doing here.

Now you can run all three services one by one (order does not matter) for Run or Debug.
*Tip:* if you want to run all three apps at the same time, you can again open Run / Edit Configuration,
select new **Compound** configuration and add all three previously added configs into it.

## VS Code

(documentation is not available yet, but it is coming...)

## Run as command-line applications

If you want to run code from command line (and eventually you might need it anyway), you can build
all three apps at once:

```
$ go build -ldflags="-w" -o .  ./cmd/...
```

This will generate three executables: `apigateway`, `filmsrv`, `usersrv` and you can launch
them directly.

E.g. to run API Gateway, run this command:

```
$ ./apigateway server --port=8080
```

To run Film and User service: 

```
$ FILMSRV_DATABASE_USER=postgres FILMSRV_DATABASE_PASSWORD=postgres FILMSRV_DATABASE_SSLMODE=disable ./filmsrv server --port=8081
```

```
$ FILMSRV_DATABASE_USER=postgres FILMSRV_DATABASE_PASSWORD=postgres FILMSRV_DATABASE_SSLMODE=disable ./usersrv server --port=8082
```

# Build and deploy to kubernetes

For local kubernetes deployment we want to use minikube. Please read up on what is minikube and how to install it.

## Database helm

Let's start installing dependencies and applications. First what we want to deploy is PostgreSQL (we will be using
bitnami deployment). Note that this is not going to be a production set up, and it will up to you to have a
properly configured PostgreSQL instances. To do so you have multiple options, e.g. to properly setup Bitnami PostgreSQL
(https://engineering.bitnami.com/articles/create-a-production-ready-postgresql-cluster-bitnami-kubernetes-and-helm.html),
use AWS RDS etc.

We start with adding bitnami PostgreSQL repository:

```shell
$ helm repo add bitnami https://charts.bitnami.com/bitnami
```

Before we proceed with installing postgresql pod - let's copy our initial script with SQL commands into kubernetes. This
will provide us with some initial database table structures and data, such as films and actors. Since we already have
script in `/deploy/docker/moviex-db` directory - we will use it:

```shell
$ kubectl create configmap db-init-schema --from-file=deploy/docker/moviex-db/init.sql
```

Here we upload SQL script into config map with key name db-init-schema. You can check that operation was successful
by running `kubectl get configmaps`.

Next let's install and run postgresql helm chart:

```shell
$ helm install postgresql bitnami/postgresql \
    --set postgresqlUsername=postgres \
    --set postgresqlDatabase=moviex \
    --set postgresqlPostgresPassword=postgres \
    --set postgresqlUsername=postgres \
    --set postgresqlPassword=postgres \
    --set initdbScriptsConfigMap=db-init-schema
```

We use some very basic user name and password, you don't want to do it for production. Also note 
`--set initdbScriptsConfigMap=db-init-schema` parameter - this instructs postgresql to run a SQL script we have previously
uploaded to configuration map.

For a convenience you can run these commands: 
```
$ export POSTGRES_PASSWORD=$(kubectl get secret --namespace default my-postgresql -o \
    jsonpath="{.data.postgresql-password}" | base64 --decode)
$ kubectl run my-postgresql-client --rm --tty -i --restart='Never' --namespace default \
    --image docker.io/bitnami/postgresql:11.14.0-debian-10-r28 \
    --env="PGPASSWORD=$POSTGRES_PASSWORD" \
    --command -- psql --host my-postgresql -U postgres -d moviex -p 5432
```

First command copies postgresql password into `POSTGRES_PASSWORD` environment. In our case it will `postgres`, but let's
say if you would omit setting password via `--set postgresqlPassword=postgres` (and it's perfectly fine if you want
to generate a random password for your production system) - then this variable will contain newly generated password.
Second command creates a new kubernetes pod based and logs you into into `psql` on this pod. Once you logged into psql -
you can try to enter `\l` command to see the list of available database, if everything is OK, you should be able to
see `moviex` in a list. You can also try to enter `\dt *.*` to see all tables available (including postgres system tables).

If by some reasons you are seeing authentication error such as `password authentication failed for user "postgres"` - then
there is a chance that you have previously was trying to deploy postgres with different credentials and volume was not
properly removed. You can try to uninstall it fully and delete all leftovers (this is a good set of commands if
you want to start everything from scratch):

```
$ helm uninstall postgresql
$ kubectl delete pvc -l app.kubernetes.io/name=postgresql
```

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
installed. Sure, we have only pod running for database, which could be OK or not OK for your production environment,
but as we mentioned in the beginning of this section - there are alternatives.

One more thing we want to do is to verify that postgres service is accessible within our cluster. The easiest way
is just to temporary launch another pod instance in our cluster:

```shell
$ kubectl run access-client --rm --tty -i --restart='Never' --namespace default --image busybox --command -- sh
```

and once it is launched and you see remote shell, you can probe a service and a postgres port such, such as
`telnet my-postgresql 5432`. If you see that telnet is successfully connected, you can exit, everything seems
to be working great at this stage.
