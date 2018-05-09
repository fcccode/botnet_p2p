package botnet_p2p;

import botnet_p2p.MessageOuterClass.Message;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.nio.ByteBuffer;
import java.nio.channels.*;
import java.util.Iterator;
import java.util.Set;


public class Server extends Thread {

    private static final Logger logger = LogManager.getLogger(Server.class);
    private final int port;
    private ServerSocketChannel serverSocketChannel;
    private Selector selector;

    Server(int port) {
        this.port = port;
    }

    @Override
    public void run() {
        try {
            selector = Selector.open();
            serverSocketChannel = ServerSocketChannel.open();
            serverSocketChannel.configureBlocking(false);
            serverSocketChannel.bind(new InetSocketAddress("localhost", port));

            SelectionKey selectionKey = serverSocketChannel.register(selector, SelectionKey.OP_ACCEPT);

            while (true) {
                // blocking call, waiting for at least one ready channel
                int channels = selector.select();

                Set<SelectionKey> selectedKeys = selector.selectedKeys();
                Iterator<SelectionKey> it = selectedKeys.iterator();
                while (it.hasNext()) {
                    SelectionKey key = it.next();

                    if (key.isAcceptable()) {
                        logger.info("new connection is possible");
                        handleNewConnection(selector);
                    }

                    if (key.isReadable()) {
                        handleNewMessage(key.channel(), selector);
                    }
                    it.remove();
                }

            }
        } catch (IOException e) {
            if (isInterrupted()) {
                logger.info("thread interrupted");
                Thread.currentThread().interrupt();
            } else {
                e.printStackTrace();
            }
        }
    }

    private void handleNewConnection(Selector selector) throws IOException {
        SocketChannel clientSocket = serverSocketChannel.accept();
        clientSocket.configureBlocking(false);
        clientSocket.register(selector, SelectionKey.OP_READ);
    }

    private void handleNewMessage(SelectableChannel channel, Selector selector) throws IOException {
        SocketChannel client = (SocketChannel) channel;
        ByteBuffer inputBuffer = ByteBuffer.allocate(512);
        if (client.read(inputBuffer) == -1) {
            client.close();
            return;
        }
        ByteBuffer messageBuffer = ByteBuffer.wrap(inputBuffer.array(), 0, inputBuffer.position());
        Message message = Message.parseFrom(messageBuffer);
        inputBuffer.clear();

        logger.info("message parsed");
        logger.info("message content:\r\n" + message.toString());
    }

    @Override
    public void interrupt() {
        super.interrupt();

        try {
            if (serverSocketChannel != null) {
                serverSocketChannel.close();
                logger.info("server socket closed");
            }
            if (selector != null) {
                selector.keys().forEach(selectionKey -> {
                    try {
                        selectionKey.channel().close();
                    } catch (IOException e) {
                        e.printStackTrace();
                    }
                });
                logger.info("clients sockets closed");
            }
        } catch (IOException e) {
        }
    }
}
