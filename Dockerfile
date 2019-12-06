# First stage: build the executable.
FROM golang:buster AS builder

# It is important that these ARG's are defined after the FROM statement
ARG ACCESS_TOKEN_USR="nothing"
ARG ACCESS_TOKEN_PWD="nothing"

# Create the user and group files that will be used in the running 
# container to run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'random:x:65534:65534:random:/:' > /user/passwd && \
    echo 'random:x:65534:' > /user/group

# Create a netrc file using the credentials specified using --build-arg
RUN printf "machine gitlab.com\n\
    login ${ACCESS_TOKEN_USR}\n\
    password ${ACCESS_TOKEN_PWD}\n\
    \n\
    machine api.gitlab.com\n\
    login ${ACCESS_TOKEN_USR}\n\
    password ${ACCESS_TOKEN_PWD}\n"\
    >> /root/.netrc

RUN chmod 600 /root/.netrc

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/gitlab.com/shitposting/discord-random

# Import the code from the context.
COPY . .

# Build the executable
RUN go install

# Final stage: the running container.
FROM debian:buster

# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/

# Copy the built executable
COPY --from=builder /go/bin/discord-random /home/random/discord-random

# Install dependencies and create home directory
RUN apt update && apt install -y ca-certificates; \ 
    chown -R random /home/random

# Set the workdir
WORKDIR /home/random

# Perform any further action as an unprivileged user.
USER random:random

# Run the compiled binary.
CMD ["./discord-random"]