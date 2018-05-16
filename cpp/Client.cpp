#include <sys/un.h>
#include <arpa/inet.h>
#include "Client.h"

Client::Client(std::string host, uint16_t port) {
  // setup buffer
  buffer_size = 1024;
  buffer = new char[buffer_size + 1];

  setup_socket(host, port);
}

Client::~Client() {
  delete buffer;
}

void Client::run() {
  std::string line;

  // loop to handle_client user interface
  while(getline(std::cin, line)) {
    // append a newline
    line += "\n";
    // send request
    bool success = send_request(line);
    // break if an error occurred
    if(not success)
      break;
    // get a response
    success = get_response();
    // break if an error occurred
    if(not success)
      break;
  }
  close_socket();
}

void Client::setup_socket(std::string host, uint16_t port) {
  // setup socket address structure
  memset(&server_address, 0, sizeof(sockaddr_in));
  server_address.sin_family = AF_INET;
  server_address.sin_addr.s_addr = inet_addr(host.c_str());
  server_address.sin_port = htons(port);

  connection_socket = socket(AF_INET, SOCK_STREAM, 0);
  if(!connection_socket) {
    std::cerr << "ERROR: Failed to create socket." << std::endl;
    exit(-1);
  }

  if(connect(connection_socket, (const struct sockaddr *) &server_address, sizeof(sockaddr_in)) < 0) {
    std::cerr << "ERROR: Failed to connect to server." << std::endl;
    exit(-1);
  }
}

void Client::close_socket() {
  close(connection_socket);
}

bool Client::send_request(std::string request) {
  // prepare to send request
  const char *ptr = request.c_str();
  int nleft = request.length();
  int nwritten;
  // loop to be sure it is all sent
  while(nleft) {
    if((nwritten = send(connection_socket, ptr, nleft, 0)) < 0) {
      if(errno == EINTR) {
        // the socket call was interrupted -- try again
        continue;
      } else {
        // an error occurred, so break out
        perror("write");
        return false;
      }
    } else if(nwritten == 0) {
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
  while(response.find("\n") == std::string::npos) {
    int nread = recv(connection_socket, buffer, buffer_size, 0);
    if(nread < 0) {
      if(errno == EINTR)
        // the socket call was interrupted -- try again
        continue;
      else
        // an error occurred, so break out
        return "";
    } else if(nread == 0) {
      // the socket is closed
      return "";
    }
    // be sure to use append in case we have binary data
    response.append(buffer, nread);
  }
  // a better client would cut off anything after the newline and
  // save it in a cache
  std::cout << response;
  return true;
}