FROM golang:1.18
COPY ./ /password-manager/server
WORKDIR /password-manager/server
RUN apt-get -y update
RUN apt-get -y install git
RUN mkdir /root/.ssh/
ADD ssh_password-manager_key /root/.ssh/id_rsa
RUN touch /root/.ssh/known_hosts
RUN ssh-keyscan gitlab.com >> /root/.ssh/known_hosts
RUN git clone git@gitlab.com:berizaryad_project/devops/servicepswd.git
WORKDIR ./servicepswd
RUN git checkout main
RUN git config --global user.email "password@berizaryad.ru"
RUN git config --global user.name "password-manager"
WORKDIR ../
RUN go build .
ENTRYPOINT ["./main"]
