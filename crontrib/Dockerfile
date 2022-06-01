FROM centos:6.10

RUN minorver=6.10 \
    sed -e "s|^mirrorlist=|#mirrorlist=|g" \
    -e "s|^#baseurl=http://mirror.centos.org/centos/\$releasever|baseurl=https://mirrors.tuna.tsinghua.edu.cn/centos-vault/$minorver|g" \
    -i.bak \
    /etc/yum.repos.d/CentOS-*.repo
RUN yum groupinstall 'Development tools' -y
RUN yum install openssh openssh-server openssh-clients git tar wget vim -y
RUN mkdir -p /var/run/sshd
RUN ssh-keygen -t rsa -f /etc/ssh/ssh_host_rsa_key
RUN ssh-keygen -t dsa -f /etc/ssh/ssh_host_dsa_key
RUN /bin/echo 'root:root'|chpasswd
RUN /bin/sed -i 's/.*session.*required.*pam_loginuid.so.*/session optional pam_loginuid.so/g' /etc/pam.d/sshd
RUN /bin/echo -e "LANG=\"en_US.UTF-8\"">/etc/default/local

# Install Go 1.15
RUN cd /tmp && /usr/bin/wget https://go.dev/dl/go1.15.15.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && /bin/tar -C /usr/local -xzf /tmp/go1.15.15.linux-amd64.tar.gz
RUN echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
RUN /usr/local/go/bin/go env -w GOPROXY=https://goproxy.cn
RUN rm -f /tmp/go1.15.15.linux-amd64.tar.gz

EXPOSE 22

CMD ["/usr/sbin/sshd","-D"]