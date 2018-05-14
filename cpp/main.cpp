#include <iostream>

#include "Server.h"
#include "Client.h"

int main() {
    char c;
    std::cin >> c;
    if (c == 'c') {
        Client client = Client();
        client.run();
    } else if (c == 's') {
        Server server = Server("hello");
        server.run();
    }
    return 0;
}