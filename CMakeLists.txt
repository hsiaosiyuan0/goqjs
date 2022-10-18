project(goqjs LANGUAGES C)
cmake_minimum_required(VERSION 3.0)

set(CMAKE_EXPORT_COMPILE_COMMANDS ON CACHE BOOL "Enable/Disable output of compile commands during generation.")

set(CMAKE_C_STANDARD 99)
set(CMAKE_C_COMPILER "clang" CACHE STRING "C compiler" FORCE)

file(GLOB_RECURSE QUICKJS RELATIVE ${CMAKE_SOURCE_DIR} "quickjs/build/include/*.h")
add_library(libgoqjs OBJECT
  ${PROJECT_SOURCE_DIR}/bridge.c
  ${QUICKJS}
)

target_include_directories(libgoqjs
  PRIVATE
  quickjs/build/include/
)