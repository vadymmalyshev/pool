FROM golang
RUN mkdir -p /root/.ssh /pool
COPY ./hiveon_ci_rsa /root/.ssh/.
RUN ls /root/.ssh && \
    echo 'StrictHostKeyChecking no' >> /etc/ssh/ssh_config
