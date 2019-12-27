# First stage: build the executable.
FROM golang:buster AS builder

# It is important that these ARG's are defined after the FROM statement
ARG SSH_PRIV="nothing"
ARG SSH_PUB="nothing"

# Create the user and group files that will be used in the running 
# container to run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'random:x:65534:65534:random:/:' > /user/passwd && \
    echo 'random:x:65534:' > /user/group


RUN eval $(ssh-agent -s); \
    mkdir -p ~/.ssh;  \
    echo "$SSH_PRIV" >> ~/.ssh/id_rsa; \
    echo "$SSH_PUB" >> ~/.ssh/id_rsa.pub;  \
    chmod 700 ~/.ssh;  \
    chmod 600 ~/.ssh/id_rsa;  \
    chmod 644 ~/.ssh/id_rsa.pub;  \
    git config --global url.git@gitlab.com:.insteadOf https://gitlab.com/;  \
    ssh-add ~/.ssh/id_rsa;  \
    ssh-add -l;  \
    ssh-keyscan -t rsa gitlab.com >> ~/.ssh/known_hosts;

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