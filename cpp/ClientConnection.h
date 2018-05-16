#ifndef CPP_CLIENTCONNECTION_H
#define CPP_CLIENTCONNECTION_H


#include <thread>

class ClientConnection {
public:
  int socket;
  char buffer[1024];
};


#endif //CPP_CLIENTCONNECTION_H
