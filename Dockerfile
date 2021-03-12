FROM fedora:30
RUN dnf -y install glibc.i686 cyrus-sasl-lib cyrus-sasl-plain libuuid openssl python3 gettext hostname iputils libwebsockets-devel gdb valgrind && dnf -y update && dnf clean all
ADD qpid-proton-image.tar.gz qpid-dispatch-image.tar.gz /
WORKDIR /home/qdrouterd/etc
WORKDIR /home/qdrouterd/bin
COPY ./scripts/* /home/qdrouterd/bin/
ARG version=latest
ENV VERSION=${version}
ENV QDROUTERD_HOME=/home/qdrouterd

EXPOSE 5672 55672 5671
CMD ["/home/qdrouterd/bin/launch.sh"]
