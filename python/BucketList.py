import heapq
import threading

from python.peer import Peer

def largest_differing_bit(value1, value2):
    """
    Returns index(from 0 to 127) of largest differing bit: Eg. for argument 011010...0 and 011110...0 it returns 3.
    :param value1: First id
    :param value2: Second id
    :return: index(from 0 to 127) of largest differing bit.
    """
    distance = value1 ^ value2
    length = -1
    while (distance):
        distance >>= 1
        length += 1
    return max(0, length)

class BucketList(object):
    """
    Data structure of BucketList. Basically it is list of list of size k.
    """
    def __init__(self, bucket_size, buckets_number, id):
        """
        :param bucket_size: Size of every bucket
        :param buckets_number: How many buckets to create
        """
        self.bucket_size = bucket_size
        self.buckets = [[] for i in range(buckets_number)]
        self.id = id
        self.lock = threading.Lock()

    def insert(self, peer):
        """
        Insert peer into appropriate bucket
        :param peer: Peer to insert
        """
        if peer.id != self.id:
            bucket_number = largest_differing_bit(self.id, peer.id)
            peer_info = peer.get_info()
            with self.lock:
                bucket = self.buckets[bucket_number]

                if len(bucket) >= self.bucket_size:
                    bucket.pop(0)
                if peer_info not in bucket:
                    bucket.append(peer_info)

    def nearest_nodes(self, key, limit=None):
        num_results = limit if limit else self.bucket_size
        with self.lock:
            def keyfunction(peer):
                return key ^ peer.get_info()[2]
            peers = (peer for bucket in self.buckets for peer in bucket)
            best_peers = heapq.nsmallest(num_results, peers, keyfunction)
        return [Peer(peer.get_info()[0:2]) for peer in best_peers]