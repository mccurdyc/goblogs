# Containerizing Go and a Use Case for Multi-Stage Docker Builds
###### Colton J. McCurdy (@McCurdyColton), Clayton Northey (@claythegreat11), Ofiliojo Ichaba (@Ofiliojo)
---

## Wide Adoption of Go in the Service-Oriented Architecture at StockX

At StockX, we have followed the industry in that we have a service-oriented, or
microservice, architecture. A microservice architecture is an architecture style
where each service has a very specific task, e.g., providing product
information, generating shipping labels, calculating pricing information, etc..
The suite of smaller services, often communicating over HTTP, constitute the
application as a whole.
While this post does allude to the architecture style at StockX, it will not discuss
the details about service-oriented architectures, a post that does, is written by Martin Fowler,
titled ["Microservices: a definition of this new architectural term"](https://martinfowler.com/articles/microservices.html).

A majority of the backend services at StockX are currently written in Golang, or just
Go, with a strong recommendation that new services also be written in Go, unless
there is a rational reason not to. One reason why we chose to adopt Go for many
of our services at StockX
was because of its ease of use and focus on simplicity, while also providing
performance benefits over other languages for building
servers. The choice to adopt Go as our primary backend language was not without
industry influence.

In an official Go Blog post titled, ["Eight Years of Go"](https://blog.golang.org/8years),
analyst Donnie Berkholz's [post](http://redmonk.com/dberkholz/2014/03/18/go-the-emerging-language-of-cloud-infrastructure/)
deeming Go "the emerging language of cloud infrastructure" is referenced. In
Berkholz's post, he uses commit data from an Ohloh data set which contains information
about more than 600,000 free and open-source software (FOSS) projects.

One of the larger (9,227,099 lines of code (LOC) according to [OpenHub](https://www.openhub.net/p/docker/analyses/latest/languages_summary)),
well-known open-source projects using Go is Docker. 83.3% or 7,846,057 LOC are written
in Go according to OpenHub, where the project is broken down by language. According
to Docker, they are the "world's leading containerization platform." Containerization is a way to
house an application and its dependencies so that the application can be run on
a wider range of hosts, yet still behave as intended. Additionally, containerizing
an application facilitates running that application in the cloud.

## Building and Running a Go Service

**Note: this post will not include the necessary steps to install Go or set up
a Go development environment, for that, refer to the [Getting Started](https://golang.org/doc/install)
page on the official website for Go.**

**Conventions Used in this Post:** In this post, I will use `$` to indicate a command
that needs to be run in your terminal. If this symbol is not present, it will
most-likely mean it is the response or output from the previous command.


Since a Go application's source and dependencies can be compiled into a single,
static binary, it is easily containerized. The process of building a self-contained
binary and containerizing a Go application will be demonstrated in this post. The purpose of this
post is not to discuss how to program in Go, therefore, the application that we
will be containerizing is
provided [in this public GitHub repository](https://github.com/mccurdyc/examples).

If you look through the code, this is a simple server, running on port 8080 by default,
that will respond with a message "hello".
This server can be run on your local machine, assuming you have Go installed, with the following command:

```
$ go run main.go
```

To interact with the server that is running, you can use the Unix command, cURL,
or `curl` in another window in your terminal. For our example, the URL that we will use is `localhost:8080`.
You should see a response similar to the following displayed in your terminal.

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
Go is not installed. For our application, we can use the Go Toolchain to build a single, static binary with the
following command in the root of the project:

```
$ go build -o bin/hello .
```

In the above command, the `-o` flag, which is only allowed when compiling a single package,
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

Additionally, if we wanted to build a binary that would be executable on an operating
system or processor architecture other than the that of the development environment, we can again leverage the
Go Toolchain's `build` command, except this time setting two environment variables.
In this example, we will be build a binary for a system running a security-oriented,
lightweight distribution of Linux, namely [Alpine](https://alpinelinux.org/). In order
to build the binary for this distribution of Linux, we need to set `GOOS=linux` and
`GOARCH=amd64`. Now, our build command looks as follows:

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

Building a Docker image can be done through the use of a [Dockerfile](https://docs.docker.com/engine/reference/builder/).
A Dockerfile is a text document that contains the commands necessary to assemble
an image, which alternatively could be called on the command line. The `docker build`
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

When a Dockerfile is present in the current directory, the following
command will build a Docker image tagged as `stockx/hello`:

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

Note that the size of the image is around 730MB. We will use this number
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

Now, to see the running container, use the following command:

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
large for what we are trying to do. We can reduce the size of the image significantly by taking
advantage of Go's ability to be compiled into a single, static binary and run the binary
on a different distribution of Linux, such as the aforementioned, lightweight Alpine distribution.

First, build the binary for the runtime environment, in this case, it will be Alpine
Linux, so we can use the command from above:

```
$ GOOS=linux GOARCH=amd64 go build -o bin/hello .
```

And let's use `Dockerfile_2` with Alpine Linux as the base image and copy the
Go binary to our container and run the binary instead of invoking Go. Our
`Dockerfile_2` contains the following:

```
FROM alpine:latest

RUN mkdir /app

ADD /bin/hello /app

WORKDIR /app

ENTRYPOINT ["/app/hello"]
```

```
$ docker build --rm -f Dockerfile_2 -t stockx/hello .
```

Now, if we run `docker images` and look at the size of the image, it is significantly
smaller, over 700MB smaller. It is now around 10MB instead of the previous 730MB.

Output:
```
REPOSITORY                            TAG                 IMAGE ID            CREATED             SIZE
stockx/hello        latest              a585021da814        Less than a second ago   10.7MB
```

Again, to start the container and once started, list the running containers, use the following
commands, respectively:

```
$ docker run -p 8080:8080 stockx/hello
```

```
$ docker ps
```

Output:
```
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
750208ddc2b4        stockx/hello        "/app/hello"        Less than a second ago   Up 1 second         0.0.0.0:8080->8080/tcp   cocky_shannon
```

## Difficulty Using 3rd-Party Libraries

At StockX, many of our Go services have their own database.
In an attempt to make working with the database more repeatable and reproducible,
we use an object-relational-mapping (ORM) library which also includes a migration tool.
However, the source for the migration tool must also be added to the container if
we want to actually use it in the container. Really, all that we need is the binary
for the migration tool in the container. This means that first, we must build the
binary, but considering that our containers are running a different operating system
and processor architecture than our development machines, so cross-compiling is necessary.
An added difficulty, one making simply cross-compiling for the runtime environment
not so simple is that the migration tool imports C libraries.

One way to achieve the ability to build and run the migration tool in the container
would be to use the `golang:alpine` base image --- similar to when we first containerized
our application --- and then run a `go get` and `go install`
in the container. Again, this is less than ideal because now our container is
doing much more work than it should have too and it is getting bloated. Another
way, and ultimately what we identified as the best way to add this migration
tool was to a [multi-stage docker build](https://docs.docker.com/engine/userguide/eng-image/multistage-build/#before-multi-stage-builds).
What this allows for is building a preliminary image where you can clutter it
up as much as you want and then feed data into the next image in the multi-stage build.

To do this, basically all that needs done is to add the contents of two Dockerfiles
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
`soda` binary is accessible in the container. For the proceeding command, you will
need to grab the `CONTAINER_ID` from `docker ps`. To check this, `exec` into the running
container with the following command:

```
docker exec -it CONTAINER_ID /bin/sh
```

Once in the container, ensure that the binary that we added to `/bin` is accessible
by running `./bin/soda`. If everything works out as expected, you should see something
like the following:

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
