cmake_minimum_required(VERSION 3.16)
project(orcs_cpp_node_exporter VERSION 0.1.0 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

add_library(orcs_process STATIC
  src/core/process_discovery.cpp
  src/core/process_terminator.cpp
)
target_include_directories(orcs_process PUBLIC include)

add_library(orcs_telemetry STATIC
  src/core/system_telemetry.cpp
)
target_include_directories(orcs_telemetry PUBLIC include)

add_library(orcs_platform STATIC
  src/platform/linux/os_adapter_linux.cpp
  src/platform/windows/os_adapter_windows.cpp
  src/platform/macos/os_adapter_macos.cpp
)
target_include_directories(orcs_platform PUBLIC include)

add_executable(orcs-node-agent src/main.cpp)
target_link_libraries(orcs-node-agent PRIVATE orcs_process orcs_telemetry orcs_platform)

enable_testing()
add_executable(orcs_cpp_tests tests/smoke_test.cpp)
target_link_libraries(orcs_cpp_tests PRIVATE orcs_process orcs_telemetry orcs_platform)
add_test(NAME orcs_cpp_smoke COMMAND orcs_cpp_tests)
