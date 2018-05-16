#ifndef CPP_CLIENT_H
#define CPP_CLIENT_H

#include <errno.h>
#include <netdb.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <unistd.h>
#include <sys/un.h>
#include <arpa/inet.h>

#include <fstream>
#include <iostream>
#include <string>

class Client {
 public:
  Client(std::string host, uint16_t port);
  ~Client();

  void run();

 protected:
  void setup_socket(std::string host, uint16_t port);
  void close_socket();

  bool send_request(std::string);
  std::string get_response();

 private:
  int connection_socket;
  struct sockaddr_in server_address;

  int buffer_size;
  char *buffer;
};

#endif //CPP_CLIENT_H
