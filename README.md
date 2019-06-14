# qdrouterd-image
Builds an image of the Apache qpid dispatch router designed for use with kubernetes and openshift

e.g. to build:

```
make && docker build -t quay.io/interconnectedcloud/qdrouterd:latest . && docker push quay.io/interconnectedcloud/qdrouterd:latest
```

to run:

```
docker run -it quay.io/interconnectedcloud/qdrouterd:latest
```
