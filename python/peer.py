import python.Message_pb2
import random

class Peer(object):
    """
    Peer
    """

    def __init__(self, host, port):
        self.host, self.port = host, port
        local_random = random.Random()
        local_random.seed(int(''.join(host.split('.')))*int(port))
        self.id = local_random.getrandbits(128)

    def address(self):
        return (self.host, self.port)

    def get_info(self):
        return self.host, self.port, self.id

    def ping(self, socket=None):
        data = "Hello"
        self._sendmessage(data, socket)
        self._receivemessage(socket)

    def _sendmessage(self, message, sock=None):
        # HERE ENCODE MESSAGE
        mess = python.Message_pb2.Message()
        mess.TYPE = mess.JOIN
        if sock:
            sock.sendall(mess.SerializeToString())

    def _receivemessage(self, socket):
        mess = socket.recv(12000)
        message = python.Message_pb2.Message()
        message.ParseFromString(mess)

        print(mess)
