FROM almalinux:latest
LABEL version="1.0"
LABEL description="Development environment"

# Configure openssh
RUN yum update -y
RUN yum install -y openssh openssh-server openssh-clients zsh tar wget sudo xz git
RUN mkdir -p /var/run/sshd
RUN ssh-keygen -t rsa -f /etc/ssh/ssh_host_rsa_key
RUN ssh-keygen -t dsa -f /etc/ssh/ssh_host_dsa_key
RUN /bin/echo 'root:huangqw'|chpasswd
RUN useradd -s /bin/zsh -G wheel huangqw
RUN /bin/echo 'huangqw:huangqw'|chpasswd
RUN /bin/echo 'huangqw ALL=(ALL) NOPASSWD: ALL' > /etc/sudoers.d/huangqw
RUN /bin/sed -i 's/.*session.*required.*pam_loginuid.so.*/session optional pam_loginuid.so/g' /etc/pam.d/sshd
RUN /bin/echo -e "LANG=\"en_US.UTF-8\"">/etc/default/local
# RUN /usr/bin/timedatectl set-timezone Asia/Shanghai

# Install Zig
RUN /usr/bin/wget https://ziglang.org/download/0.9.1/zig-linux-x86_64-0.9.1.tar.xz -O /opt/zig-linux-x86_64-0.9.1.tar.xz
RUN cd /opt && /usr/bin/tar -xf /opt/zig-linux-x86_64-0.9.1.tar.xz
RUN /usr/bin/ln -s /opt/zig-linux-x86_64-0.9.1/zig /usr/sbin/zig

# Install Go
RUN cd /tmp && /usr/bin/wget https://go.dev/dl/go1.18.1.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf /tmp/go1.18.1.linux-amd64.tar.gz
RUN echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
RUN go env -w GOPROXY=https://goproxy.cn
RUN rm -f /tmp/go1.18.1.linux-amd64.tar.gz

USER huangqw
# Uses some bundled plugins and installs some more from github
RUN sh -c "$(wget -O- https://github.com/deluan/zsh-in-docker/releases/download/v1.1.2/zsh-in-docker.sh | sed -r 's/centos \| amzn/centos \| amzn \| almalinux/g')" -- \
    -p git \
    -p ssh-agent \
    -p golang \
    -p z \
    -p https://github.com/zsh-users/zsh-autosuggestions \
    -p https://github.com/zsh-users/zsh-completions

USER root

EXPOSE 22

CMD ["/usr/sbin/sshd","-D"]