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

#include <fstream>
#include <iostream>
#include <string>

class Client {
public:
    Client();

    ~Client();

    void run();

protected:
    void create();

    void close_socket();

    void echo();

    bool send_request(std::string);

    bool get_response();

    int server_;
    int buflen_;
    char *buf_;

private:
    const char *socket_name_;
};

#endif //CPP_CLIENT_H
