package botnet_p2p;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import botnet_p2p.MessageOuterClass.Message;

import java.io.IOException;

public class TestNode {
    private static final Logger logger = LogManager.getLogger(TestNode.class);
    private Client client;

    TestNode() throws IOException, InterruptedException {
        Runtime.getRuntime().addShutdownHook(new ShutdownHandler());
        client = new Client();
        client.start();
        Thread.sleep(4000);

        Message.Builder builder = Message.newBuilder();
        builder.setTYPE(Message.MessageType.JOIN);
        builder.setSender("sender name");
        Message message = builder.build();

        client.sendMessage(message, "127.0.0.1", 3000);
        client.sendMessage(message, "127.0.0.1", 3000);
    }

    public static void main(String args[]) {
        logger.info("starting");
        try {
            TestNode testNode = new TestNode();
        } catch (IOException | InterruptedException e) {
            e.printStackTrace();
        }
    }

    class ShutdownHandler extends Thread {
        @Override
        public void run() {
            super.run();
            logger.info("closing requested");
            client.interrupt();
        }
    }
}
