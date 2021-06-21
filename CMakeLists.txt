cmake_minimum_required(VERSION 3.16)
project(k8deployment C)

set(CMAKE_C_STANDARD 11)

file(GLOB_RECURSE LIB_SOURCES_SO RELATIVE ${CMAKE_SOURCE_DIR} "libs/*.so")
file(GLOB_RECURSE LIB_SOURCES RELATIVE ${CMAKE_SOURCE_DIR} "libs/*.c")
file(GLOB_RECURSE SRC_SOURCES RELATIVE ${CMAKE_SOURCE_DIR} "src/*.c")

message("-- Cmake Source Dirs: ${CMAKE_SOURCE_DIR}")
message("-- Including Libs Directories: ${LIB_SOURCES}")
message("-- Including Libs Directories: ${LIB_SOURCES_SO}")
message("-- Including Src Directories: ${SRC_SOURCES}")



#target_include_directories(k8deployment PUBLIC "${PROJECT_SOURCE_DIR}/libs/kubernetes")
#target_link_directories(k8deployment PUBLIC "${PROJECT_SOURCE_DIR}/libs/kubernetes")

add_executable(k8deployment main.c ${LIB_SOURCES} ${SRC_SOURCES})
target_link_libraries(k8deployment git2 ${CMAKE_SOURCE_DIR}/libs/kubernetes/libkubernetes.so)