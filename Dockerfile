FROM fedora:33
RUN dnf -y install glibc cyrus-sasl-lib cyrus-sasl-plain cyrus-sasl-gssapi cyrus-sasl-md5 libuuid openssl gettext hostname iputils libwebsockets-devel libnghttp2 gdb valgrind && dnf -y update && dnf clean all
ADD qpid-proton-image.tar.gz qpid-dispatch-image.tar.gz /
WORKDIR /home/qdrouterd/etc
WORKDIR /home/qdrouterd/bin
COPY ./scripts/* /home/qdrouterd/bin/
ARG version=latest
ENV VERSION=${version}
ENV QDROUTERD_HOME=/home/qdrouterd

EXPOSE 5672 55672 5671
CMD ["/home/qdrouterd/bin/launch.sh"]
