# Containerizing Go and a Use Case for Multi-Stage Docker Builds
###### Colton J. McCurdy (@McCurdyColton), Clayton Northey (@claythegreat11), Ofiliojo Ichaba (@Ofiliojo)
---

## Adoption of Go in the Service-Oriented Architecture at StockX

StockX is built on a service-oriented, or
microservice architecture. Microservice architecture is a system design paradigm
in which a application system is made up of smaller services, each performing specific task e.g., providing product
information, generating shipping labels, calculating pricing information, etc.
This suite of smaller services, often communicating over HTTP, constitute the
application.
While this post alludes to the architecture at StockX, we will not be discussing
in detail, the concept of service-oriented architectures. A post by Martin Fowler,
titled ["Microservices: a definition of this new architectural term"](https://martinfowler.com/articles/microservices.html) already does this sufficiently.

Most of the backend services at StockX are currently written in Golang (or just
Go), with a strong recommendation that any new services also be written in Go, unless
there is a rational reason to do otherwise. We chose to adopt Go for many
of our services at StockX
because it focuses on simplicity, while also providing
performance benefits over server-side languages. That being said, the choice to adopt Go as our primary language was not without
industry influence. An official Go Blog post titled, ["Eight Years of Go"](https://blog.golang.org/8years) referenced
analyst Donnie Berkholz's [post](http://redmonk.com/dberkholz/2014/03/18/go-the-emerging-language-of-cloud-infrastructure/)
deeming Go "the emerging language of cloud infrastructure". In
Berkholz's post, he uses commit data from an [Ohloh](https://www.openhub.net/p/ohloh) data set which contains information
about more than 600,000 free and open-source software (FOSS) projects to support his claim.

One of the larger (9,227,099 lines of code (LOC) According to [OpenHub](https://www.openhub.net/p/docker/analyses/latest/languages_summary)),
well-known open-source Go Projects is Docker. According
to Docker, they are the "world's leading containerization platform." amd StockX uses Docker to containerize our services. Containerization is a way to
house an application and its dependencies so that the application can be run on
a wider range of hosts, yet still behave as intended. Additionally, containerizing
an application facilitates running it in the cloud.

## Building and Running a Go Service

**Note: this post will not include the necessary steps to install Go or set up
a Go development environment, for that, refer to the [Getting Started](https://golang.org/doc/install)
page on the official website for Go.**

**Conventions Used in this Post:** In this post, I will use `$` to indicate a command
that needs to be run in your terminal. If this symbol is not present, it will
most-likely mean it is the response or output from the previous command.

Since a Go application's source code and dependencies can be compiled into a single,
static binary, it is easily containerized. This post will demonstrate process of building a self-contained
binary and containerizing a Go application. The purpose of this
post is not to discuss how to program in Go, therefore, the application that we
will be containerizing is
provided in this public [GitHub repository](https://github.com/mccurdyc/examples).

The code is for a simple server, running on port 8080 by default,
that will respond with the message "hello". You can run this server on your local machine
(assuming you have Go installed), with the following command:

```
$ go run main.go
```

To interact with the server that is running, you can use the Unix command, [cURL](https://github.com/curl/curl),
or `curl` in another window of your terminal. For this example, the URL we will use is `localhost:8080`.
Running the command below in your terminal, you should receive a response.

```
$ curl localhost:8080/hello
```

Response:
```
hello
```

### Running Go in Environment Without Go Installed

However, it is often the case that the application that you are building will be run
in an environment different than the environment in which it was developed, possibly one where
Go is not installed. For this reason, in our application, we can use the Go Toolchain to build a single, static binary with the
following command in the root of the project:

```
$ go build -o bin/hello .
```

In the command above, the `-o` flag, which is only allowed when compiling a single package,
forces build to write the resulting executable or object to the named output file,
instead of the default behavior.

Now, we can run the same server by executing the binary with the following command
in the root of the project:

```
$ ./bin/hello
```

In another window:
```
$ curl localhost:8080/hello
```

Response:
```
hello
```

### Building and Running Go on a Different Operating System

To build a binary that would be executable on an operating
system or processor architecture other than the that of the development environment, we can again leverage the
Go Toolchain's `build` command, except this time setting two environment variables.
In this example, we will be build a binary for [Alpine](https://alpinelinux.org/), a system running a security-oriented,
lightweight distribution of Linux. In order
to build the binary for this distribution of Linux, we need to set `GOOS=linux` and
`GOARCH=amd64`. Now, our build command is:

```
$ GOOS=linux GOARCH=amd64 go build -o bin/hello .
```

While this application won't run on your local machine --- unless you are running Alpine Linux ---
it will now run on a system that is running Alpine. An easy way to test this out
is to use Docker's [Alpine Linux image](https://hub.docker.com/_/alpine/).

## Containerizing a Go Application

**Note: this post will not include the steps necessary to install and run Docker on your
machine. For the steps needed to install and get started with Docker on your machine, refer to the
[Install Docker](https://docs.docker.com/engine/installation/) page on Docker's
website.**

### Building an Image

A Docker image can be built using a [Dockerfile](https://docs.docker.com/engine/reference/builder/).
A Dockerfile is a text file that contains the necessary commands to assemble
a docker image; an alternative to running the individual commands on the command line. The `docker build`
command allows users to automate the image build process by executing the commands
in the Dockerfile. The contents of our Dockerfile --- named `Dockerfile_1` in order
to differentiate between Dockerfiles throughout this tutorial --- are as follows:

```
FROM golang:1.9

CMD cd $GOPATH

RUN go get github.com/mccurdyc/goblogs/docker-multi-stage-builds/hello

WORKDIR $GOPATH/src/github.com/mccurdyc/goblogs/docker-multi-stage-builds/hello

ENTRYPOINT ["go", "run", "main.go"]
```

The following command will build a Docker image tagged as `stockx/hello`:

```
$ docker build --rm -f Dockerfile_1 -t stockx/hello .
```

Output:
```
Step 1/5 : FROM golang:1.9
 ---> 138bd936fa29
Step 2/5 : RUN mkdir /app
 ---> Using cache
 ---> fd1cf9736c21
Step 3/5 : ADD . /app
 ---> b698af017d04
Step 4/5 : WORKDIR /app
Removing intermediate container d8b5ad763bf7
 ---> 0518e56b3702
Step 5/5 : ENTRYPOINT ["go", "run", "main.go"]
 ---> Running in 14d350c4d643
Removing intermediate container 14d350c4d643
 ---> 807e62c84619
Successfully built 807e62c84619
Successfully tagged stockx/hello:latest
```

To see a list of images immediately available on your machine --- including the image tagged `stockx/hello` --- run the following:

```
$ docker images
```

Output:
```
REPOSITORY                            TAG                 IMAGE ID            CREATED             SIZE
stockx/hello                          latest              49b293df7e71        9 hours ago         734MB
```

Note that the size of the image is ~730MB. We will use this number
as a point of comparison later, so keep it in mind.

### What is the Difference Between a Docker Image and a Container?

While this post will not go into detail on the topic of the differences between
a Docker image and container, the high-level concept is that a container is a running
instance of an image. To start a container and map the container ports to your
local machine, use the following command:

```
$ docker run -p 8080:8080 stockx/hello
```

Output:
```
2018/01/23 02:12:01 started server on localhost:8080
```

To see the running container, use the following command:

```
$ docker ps
```

Output:
```
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS               NAMES
4492d4f94465        stockx/hello        "go run main.go /bin…"   1 second ago        Up 3 seconds        0.0.0.0:8080->8080/tcp   romantic_neumann
```

### Building a Lightweight Image

Currently, we are using the latest Golang image as our base image, this is unnecessarily
large image for what we are trying to do. We can significantly reduce the size of the image by taking
advantage of Go's ability to be compiled into a single, static binary and run that binary
on a different distribution of Linux, such as the aforementioned, lightweight Alpine distribution.

To do this, build the runtime environment specific binary, in this case, it will be Alpine
Linux, so we can use the command from earlier:

```
$ GOOS=linux GOARCH=amd64 go build -o bin/hello .
```

For this build, let's use `Dockerfile_2`, which runs Alpine Linux as the base image, copies the
Go binary to our container, and runs the binary instead of invoking Go. Our
`Dockerfile_2` contains the following commands:

```
FROM alpine:latest

RUN mkdir /app

ADD /bin/hello /app

WORKDIR /app

ENTRYPOINT ["/app/hello"]
```

Run:
```
$ docker build --rm -f Dockerfile_2 -t stockx/hello .
```

Now, if we run `docker images` and look at the size of the image, it is significantly
smaller, over 700MB smaller. It is now ~10MB as compared to the previous ~730MB.

Output:
```
REPOSITORY          TAG                 IMAGE ID            CREATED                  SIZE
stockx/hello        latest              a585021da814        Less than a second ago   10.7MB
```

Again, to start the container, and once started, list the running containers, use the following
commands, respectively:

```
$ docker run -p 8080:8080 stockx/hello
```

```
$ docker ps
```

Output:
```
REPOSITORY          TAG                 IMAGE ID            CREATED                  SIZE            PORTS                    NAMES
750208ddc2b4        stockx/hello        "/app/hello"        Less than a second ago   Up 1 second     0.0.0.0:8080->8080/tcp   cocky_shannon
```

## Difficulty Using 3rd-Party Libraries

At StockX, many of our Go services have their own database.
In an attempt to make working with the database more repeatable and reproducible,
we use an object-relational-mapping (ORM) library that includes a migration tool.
However, to use this migration tool, the source must also be added to the container.
Since the source for the migration tool is in Go, all we need is the binary
for the migration tool in the container. To do this, we need to build the
binary, but since our containers are running a different operating system
and processor architecture than our development machines, cross-compiling is necessary.
An added complexity that makes simply cross-compiling for the runtime environment
non-trivial, is that the migration tool imports C libraries.

One way to successfully to build and run the migration tool in the container
would be to use the `golang:alpine` base image --- similar to when we first containerized
our application --- and then run a `go get` and `go install`
in the container. This is less than ideal because our container is
doing more work than it should have too and is getting bloated. An alternative that we
ultimately identified as the best way to add this migration tool was to implement a [multi-stage docker build](https://docs.docker.com/engine/userguide/eng-image/multistage-build/#before-multi-stage-builds).
What this allows us to build a preliminary image that can be as bloated as we need it
be and then feed data from that into a smaller image.

To do this, all that needs done is add the contents of two Dockerfiles
into a single Dockerfile, like so:

```
FROM golang:alpine AS prelim

RUN mkdir /app

RUN apk add --no-cache git build-base

RUN go get github.com/markbates/pop/...
RUN GOOS=linux GOARCH=amd64 \
	go build -o /go/bin/soda github.com/markbates/pop/soda

FROM alpine:latest

RUN mkdir /app

ADD bin/hello /app

COPY --from=prelim /go/bin/soda bin/

ENTRYPOINT ["/app/hello"]
```

Running our server will behave the same, but what we are interested in is that the
`soda` binary is accessible in the container. . To check this, `exec` into the running
container with the following command (For the following command, you will
need to grab the `<CONTAINER_ID>` from the `docker ps` command):

```
docker exec -it <CONTAINER_ID> /bin/sh
```

Once in the container, ensure that the binary that we added to `/bin` is accessible
by running `./bin/soda`. If everything works out as expected, you should see something
similar to the following:

Output:
```
v3.51.1

A tasty treat for all your database needs

Usage:
  soda [flags]
  soda [command]

Available Commands:
  create      Creates databases for you
  drop        Drops databases for you
  generate
  help        Help about any command
  migrate     Runs migrations against your database.
  schema      Tools for working with your database schema

Flags:
  -c, --config string   The configuration file you would like to use.
  -d, --debug           Use debug/verbose mode
  -e, --env string      The environment you want to run migrations against. Will use $GO_ENV if set. (default "development")
  -h, --help            help for soda
  -p, --path string     Path to the migrations folder (default "./migrations")
  -v, --version         Show version information

Use "soda [command] --help" for more information about a command.
```
