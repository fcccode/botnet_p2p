from python.peer import Peer
import python.Message_pb2
import socket
import socketserver
import threading

class Server(socketserver.ThreadingMixIn, socketserver.TCPServer):
    def __init__(self, server_address, request_handler_class):
        socketserver.TCPServer.__init__(self, server_address, request_handler_class)
        self.kademlia_node = None

class RequestHandler(socketserver.BaseRequestHandler):
    def handle(self):
        # self.request is the TCP socket connected to the client
        self.data = self.request.recv(12000).strip()
        message = python.Message_pb2.Message().ParseFromString(self.data)
        print("{} wrote: ".format(self.client_address[0]))
        print(message)



class KademliaNode(object):
    def __init__(self, host, port, id=None, seeds=[], requesthandler=RequestHandler):
        # Create kademlia node info object
        self.peer = Peer(host, port)
        self.other_peers = []
        # Create Server
        self.server = Server(self.peer.address(), requesthandler)
        self.server.kademlia_node = self

        self.server_thread = threading.Thread(target=self.server.serve_forever)
        self.server_thread.daemon = True
        self.server_thread.start()

        self.peers_connections = []
        self.bootstrap(seeds)

    # Boostrap the network with a list of bootstrap nodes
    def bootstrap(self, bootstrap_nodes = []):
        for bnode in bootstrap_nodes:
            print("Bootstraping")
            sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            self.peers_connections.append(sock)

            boot_peer = Peer(bnode[0], bnode[1])
            self.other_peers.append(boot_peer)

            sock.connect(boot_peer.address())
            boot_peer.ping(socket=sock)