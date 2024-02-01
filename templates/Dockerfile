FROM ubuntu:20.04

RUN apt-get update -y
RUN apt-get upgrade -y

RUN apt-get install -y software-properties-common
RUN add-apt-repository ppa:rock-core/qt4 -y

RUN apt-get install -y git
RUN apt-get install -y curl wget unzip
RUN apt-get install -y build-essential
RUN apt-get install -y libtool
RUN apt-get install -y autotools-dev
RUN apt-get install -y automake
RUN DEBIAN_FRONTEND='noninteractive' apt-get install -y pkg-config
RUN apt-get install -y python3
RUN apt-get install -y bsdmainutils
RUN DEBIAN_FRONTEND='noninteractive' apt-get install -y cmake
RUN apt-get install -y rename
RUN apt-get install -y xxd
RUN apt-get install -y bison
RUN apt install -y g++-mingw-w64-x86-64 mingw-w64-x86-64-dev nsis
RUN apt install -y htop
RUN apt install -y gperf

RUN apt-get install -y libboost-all-dev libzmq3-dev libminiupnpc-dev  
RUN apt-get install -y libssl-dev libevent-dev
RUN apt-get install -y libqrencode-dev
RUN apt-get install -y libqt4-dev libprotobuf-dev protobuf-compiler
RUN apt-get install -y python3-pip

# configure tool networking
RUN apt-get install -y net-tools
RUN apt-get install -y iputils-ping
RUN apt-get install -y nano
RUN apt-get install -y nmap

# configure tool release
RUN apt-get install -y zip

# configure posix
RUN update-alternatives --set x86_64-w64-mingw32-gcc /usr/bin/x86_64-w64-mingw32-gcc-posix
RUN update-alternatives --set x86_64-w64-mingw32-g++ /usr/bin/x86_64-w64-mingw32-g++-posix




# RUN DEBIAN_FRONTEND='noninteractive' apt-get install -y --no-install-recommends xserver-xorg-input-evdev xserver-xorg-input-all

WORKDIR /root

RUN mkdir scripts
RUN mkdir qt
COPY ./scripts/* scripts
RUN ./scripts/install-berkeley.sh

# copying autostart
ADD ./tool/start_daemon.sh /root/start_daemon.sh
RUN chmod 777 /root/start_daemon.sh


# install ssh untuk ansible
# sshd
ENV ROOT_PASSWORD password
RUN mkdir /run/sshd; \
    apt install -y openssh-server; \
    sed -i 's/^#\(PermitRootLogin\) .*/\1 yes/' /etc/ssh/sshd_config; \
    sed -i 's/^\(UsePAM yes\)/# \1/' /etc/ssh/sshd_config; \
    apt clean;

# entrypoint
RUN { \
    echo '#!/bin/bash -eu'; \
    echo 'ln -fs /usr/share/zoneinfo/${TZ} /etc/localtime'; \
    echo 'echo "root:${ROOT_PASSWORD}" | chpasswd'; \
    echo 'nohup /usr/sbin/sshd -D -e &'; \
    echo '"$@"'; \
    } > /usr/local/bin/entry_point.sh; \
    chmod +x /usr/local/bin/entry_point.sh;




ENV LD_LIBRARY_PATH="$LD_LIBRARY_PATH:/usr/local/lib"
ENV TZ Asia/Bangkok


# installing github cli
RUN type -p curl >/dev/null || (apt update && apt install curl -y) 
RUN curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg \
&& chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg \
&& echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | tee /etc/apt/sources.list.d/github-cli.list > /dev/null \
&& apt update \
&& apt install gh -y



EXPOSE 22

ENTRYPOINT ["entry_point.sh"]
CMD [ "/root/start_daemon.sh" ]
