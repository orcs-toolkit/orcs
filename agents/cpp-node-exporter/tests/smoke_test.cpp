#include <cassert>

#include "orcs/os_adapter.hpp"
#include "orcs/system_telemetry.hpp"

int main() {
  orcs::LinuxAdapter adapter;
  assert(adapter.platform_name() == "linux");

  orcs::SystemTelemetry telemetry;
  auto snapshot = telemetry.sample();
  assert(!snapshot.os_type.empty());

  return 0;
}
