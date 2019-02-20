# qdrouterd-image
Builds an image of the Apache qpid dispatch router designed for use with kubernetes and openshift

e.g. to build:

```
make && docker build -t quay.io/interconnectedcloud/qdrouterd:1.5.0 . && docker push quay.io/interconnectedcloud/qdrouterd:1.5.0
```

to run:

```
docker run -it quay.io/interconnectedcloud/qdrouterd:1.5.0
```
