# qpid-dispatch k8s tests


## k8s directory

The `k8s` directory holds test suites that have the ability
to run different tests to validate the qdrouterd image on
a Kubernetes cluster.

These test suites will run atomic test containers inside the
cluster and validate the results.

## container-images directory

At the `container-images` directory, there is a directory tree
organized by the programming language and atomit test name.

Each atomic test must provide a `Dockerfile` that builds a container
image that can be used inside a Kubernetes cluster.

Along with the `Dockerfile`, you must also add an entry in the
`test/container-images/Makefile` so the image is built and pushed
to the docker registry.

## setenv.sh file

All environment variables that can be customized externally before
running the test suites are defined at the `setenv.sh` file. This
script does not need to be sourced locally, as the variables stored
in it contains the default values, and so all you need to do is to
customize them as needed.
