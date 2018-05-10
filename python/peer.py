import hashlib
import json
import python.Message_pb2

class Peer(object):
    """
    Peer
    """

    def __init__(self, host, port):
        self.host, self.port = host, port

    def address(self):
        return (self.host, self.port)

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
