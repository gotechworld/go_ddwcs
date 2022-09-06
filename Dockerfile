FROM golang:1.14-alpine3.12 as builder

# arguments
ARG SSH_KEY

#install dependencies
RUN apk update && apk upgrade && apk add --no-cache git openssh ca-certificates

# ENV vars
ENV GO111MODULE=on
ENV GOPROXY=direct
ENV GOSUMDB=off
ENV GIT_TERMINAL_PROMPT=1

# change the way e download from our private bitbucket
RUN git config --global --add url."git@gitlab.altex.ro:".insteadOf "https://gitlab.altex.ro/"

# create ssh dir
RUN mkdir -p ~/.ssh

# Copy SSH key for git private repos
RUN echo "${SSH_KEY}" >> ~/.ssh/id_rsa

# add bitbucket to known hosts
RUN touch ~/.ssh/known_hosts

RUN ssh-keyscan -t rsa gitlab.altex.ro >> ~/.ssh/known_hosts

### Skip Host verification for git user
RUN echo "Host gitlab.altex.ro" \
         "\n\t" \
         "HostName gitlab.altex.ro" \
         "\n\t" \
         "User git" \
         "\n\t" \
         "IdentityFile ~/.ssh/id_rsa" \
         "\n\t" \
         "AddKeysToAgent yes" \
         "\n\t" \
         "StrictHostKeyChecking no"  >> ~/.ssh/config

# set permissions
RUN chmod -R 600 ~/.ssh/

# add key to agent
RUN eval `ssh-agent -s` && ssh-add -k ~/.ssh/id_rsa

# create app dir
RUN mkdir /app

# copy files to current app dir
COPY . /app/

#move working directory
WORKDIR /app

# Download depemdencies
RUN go clean -modcache
RUN go get -v -d

# Build the source
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o go_app

# The second and final stage
FROM scratch

# Copy the certs from the builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary from the builder stage and settings
COPY --from=builder /app/config/ /app/config
COPY --from=builder /app/go_app /app/go_app

# start app
ENTRYPOINT ["/app/go_app"]

EXPOSE 8080