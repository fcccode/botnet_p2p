package botnet_p2p;

import java.net.InetSocketAddress;
import java.nio.channels.SocketChannel;

public class BotnetNode {
    SocketChannel socketChannel;
    InetSocketAddress address;
    NodeStatus status;

    public BotnetNode(InetSocketAddress address, NodeStatus status) {
        this.address = address;
        this.status = status;
    }

}
