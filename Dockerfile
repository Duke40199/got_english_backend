FROM golang:1.16

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/golang/got_english_backend

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . $GOPATH/src/github.com/golang/got_english_backend

RUN go env -w GO111MODULE=auto
RUN go env
# Install the package
RUN chmod +x $GOPATH/src/github.com/golang/got_english_backend/install_dependencies.sh
RUN $GOPATH/src/github.com/golang/got_english_backend/install_dependencies.sh

# This container exposes port 80 to the outside world
EXPOSE 80

# Run the executable
CMD ["bash", "start_app_production.sh" ] 

