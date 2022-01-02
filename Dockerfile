FROM oraclelinux:7-slim

LABEL org.opencontainers.image.authors="moihn@hotmail.com"

ARG BS_GROUP
ARG BS_GID
ARG BS_USER
ARG BS_UID

RUN yum install -y oracle-instantclient-release-el7 sudo \
    && yum install -y oracle-instantclient-basic \
    && yum -y clean all \
# user and group
    && groupadd $BS_GROUP -g $BS_GID \
    && useradd -g $BS_GROUP -u $BS_UID -d /home/$BS_USER -s /bin/bash $BS_USER \
    && chown -R $BS_USER:$BS_GROUP /home/$INF_BS_USER \
    && echo "$BS_USER ALL=(root) NOPASSWD:ALL" > /etc/sudoers.d/user \
    && chmod 0440 /etc/sudoers.d/user
ADD cmd/whzy_rest_api/whzy_rest_api /usr/local/bin/
ADD config.yml /usr/local/etc/
ENTRYPOINT ["whzy_rest_api", "-config", "/usr/local/etc/config.yml"]