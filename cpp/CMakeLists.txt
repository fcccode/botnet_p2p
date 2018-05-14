cmake_minimum_required(VERSION 3.8)
project(cpp)

set(CMAKE_CXX_STANDARD 14)

set(SOURCE_FILES main.cpp Client.cpp Server.cpp ClientConnection.h)
add_executable(cpp ${SOURCE_FILES})
SET(CMAKE_CXX_FLAGS ${CMAKE_CXX_FLAGS} "-pthread")