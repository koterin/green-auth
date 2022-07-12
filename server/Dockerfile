FROM golang:1.18
COPY ./ /green-auth/server

WORKDIR /green-auth/server
RUN apt-get -y update
RUN apt-get -y install git
RUN mkdir /root/.ssh/
ADD ssh_green-auth_key /root/.ssh/id_rsa
RUN touch /root/.ssh/known_hosts
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts
RUN git clone git@github.com:koterin/green-auth.git

WORKDIR ./green-auth
RUN git checkout service
RUN git config --global user.email "green.auth@ktrn.com"
RUN git config --global user.name "green-auth"

WORKDIR ../
RUN go build .
ENTRYPOINT ["./main"]
