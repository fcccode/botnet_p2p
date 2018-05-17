#ifndef CPP_ROUTINGTABLE_H
#define CPP_ROUTINGTABLE_H


#include <cstdint>
#include <netinet/in.h>
#include <mutex>

class RoutingTable {

  class Entry {
  public:
    in_addr_t host;
    uint16_t port;
    uint64_t id;
  };

  std::mutex mutex;

  void insert(Entry);
  uint64_t *nearest_nodes();

  static const int buckets_number = 64;

  static unsigned int largest_differing_bit(uint64_t a, uint64_t b) {
    uint64_t distance = a ^ b;
    uint64_t length = 0;

    while(distance) {
      distance >>= 1;
      ++length;
    }
    return length;
  }
};


#endif //CPP_ROUTINGTABLE_H
