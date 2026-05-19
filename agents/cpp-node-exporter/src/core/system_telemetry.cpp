#include "orcs/system_telemetry.hpp"

#include <fstream>
#include <string>

namespace orcs {

TelemetrySnapshot SystemTelemetry::sample() const {
  TelemetrySnapshot snapshot{};
#ifdef __linux__
  snapshot.os_type = "linux";
  snapshot.cpu_load = 0;
  std::ifstream meminfo("/proc/meminfo");
  std::string key;
  unsigned long long value;
  std::string unit;
  while (meminfo >> key >> value >> unit) {
    if (key == "MemTotal:") {
      snapshot.total_mem_bytes = value * 1024;
    }
    if (key == "MemAvailable:") {
      snapshot.free_mem_bytes = value * 1024;
    }
  }
#else
  snapshot.os_type = "unknown";
  snapshot.cpu_load = 0;
  snapshot.total_mem_bytes = 0;
  snapshot.free_mem_bytes = 0;
#endif
  return snapshot;
}

}  // namespace orcs
