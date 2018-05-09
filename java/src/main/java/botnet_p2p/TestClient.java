package botnet_p2p;

import botnet_p2p.MessageOuterClass.Message;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.net.Socket;

public class TestClient {
    private static final Logger logger = LogManager.getLogger(TestClient.class);

    public static void main(String args[]) {

        Message.Builder builder = Message.newBuilder();
        builder.setTYPE(Message.MessageType.JOIN);
        builder.setSender("sender name");
        Message message = builder.build();
        logger.info("message:\r\n" + message.toString());
        logger.info("sending message");

        Socket client = new Socket();

        try {
            client.connect(new InetSocketAddress("localhost", 3000));
            message.writeTo(client.getOutputStream());
            client.getOutputStream().flush();
            client.close();
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
