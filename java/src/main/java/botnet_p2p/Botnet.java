package botnet_p2p;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class Botnet {

    private static final Logger logger = LogManager.getLogger(Server.class);
    private Server server;


    Botnet(int port) {
        Runtime.getRuntime().addShutdownHook(new ShutdownHandler());
        server = new Server(port);
        server.start();
    }

    public static void main(String args[]) {
        logger.info("starting");
        Botnet botnet = new Botnet(3000);
    }

    class ShutdownHandler extends Thread {
        @Override
        public void run() {
            super.run();
            logger.info("closing requested");
            server.interrupt();
        }
    }

}