#include <iostream>

#include "Client.h"
#include "Server.h"


int main() {
    char c;
    std::cin >> c;
    if (c == 'c') {
        Client client = Client("127.0.0.1", 8999);
        client.run();
    } else if (c == 's') {
        Server server = Server(8999);
        server.run();
    }
    return 0;
}