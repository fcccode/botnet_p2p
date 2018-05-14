#include <sys/un.h>
#include "Client.h"

Client::Client() {
    // setup variables
    buflen_ = 1024;
    buf_ = new char[buflen_ + 1];
    socket_name_ = "/tmp/unix-socket";
}

Client::~Client() {
    delete buf_;
}

void Client::run() {
    // connect to the server and run echo program
    create();
    echo();
}

void Client::create() {
    struct sockaddr_un server_addr;

    // setup socket address structure
    bzero(&server_addr, sizeof(server_addr));
    server_addr.sun_family = AF_UNIX;
    strncpy(server_addr.sun_path, socket_name_, sizeof(server_addr.sun_path) - 1);

    // create socket
    server_ = socket(PF_UNIX, SOCK_STREAM, 0);
    if (!server_) {
        perror("socket");
        exit(-1);
    }

    // connect to server
    if (connect(server_, (const struct sockaddr *) &server_addr, sizeof(server_addr)) < 0) {
        perror("connect");
        exit(-1);
    }
}

void Client::close_socket() {}

void Client::echo() {
    std::string line;

    // loop to handle_client user interface
    while (getline(std::cin, line)) {
        // append a newline
        line += "\n";
        // send request
        bool success = send_request(line);
        // break if an error occurred
        if (not success)
            break;
        // get a response
        success = get_response();
        // break if an error occurred
        if (not success)
            break;
    }
    close_socket();
}

bool Client::send_request(std::string request) {
    // prepare to send request
    const char *ptr = request.c_str();
    int nleft = request.length();
    int nwritten;
    // loop to be sure it is all sent
    while (nleft) {
        if ((nwritten = send(server_, ptr, nleft, 0)) < 0) {
            if (errno == EINTR) {
                // the socket call was interrupted -- try again
                continue;
            } else {
                // an error occurred, so break out
                perror("write");
                return false;
            }
        } else if (nwritten == 0) {
            // the socket is closed
            return false;
        }
        nleft -= nwritten;
        ptr += nwritten;
    }
    return true;
}

bool Client::get_response() {
    std::string response = "";
    // read until we get a newline
    while (response.find("\n") == std::string::npos) {
        int nread = recv(server_, buf_, 1024, 0);
        if (nread < 0) {
            if (errno == EINTR)
                // the socket call was interrupted -- try again
                continue;
            else
                // an error occurred, so break out
                return "";
        } else if (nread == 0) {
            // the socket is closed
            return "";
        }
        // be sure to use append in case we have binary data
        response.append(buf_, nread);
    }
    // a better client would cut off anything after the newline and
    // save it in a cache
    std::cout << response;
    return true;
}