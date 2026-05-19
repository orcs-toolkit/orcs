#pragma once

#include <string>

namespace orcs {

struct TelemetrySnapshot {
  std::string os_type;
  double cpu_load;
  unsigned long long total_mem_bytes;
  unsigned long long free_mem_bytes;
};

class SystemTelemetry {
 public:
  TelemetrySnapshot sample() const;
};

}  // namespace orcs
