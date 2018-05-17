#ifndef CPP_SERVER_H
#define CPP_SERVER_H

#include <errno.h>
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <unistd.h>
#include <iostream>
#include <string>
#include <thread>
#include <vector>

#include "ClientConnection.h"

class Server {
 public:
  Server(uint16_t port);
  ~Server();
  Server(Server &&) = default;

  void run();

 protected:
  void setup_socket(uint16_t port);
  void close_socket();

  void handle_client(ClientConnection client);

  std::string get_request(ClientConnection client);
  bool send_response(ClientConnection client, std::string);

 private:
  int connection_socket;
  struct sockaddr_in server_address, client_address;
  std::vector<std::thread> client_connections;
};

#endif  // CPP_SERVER_H
