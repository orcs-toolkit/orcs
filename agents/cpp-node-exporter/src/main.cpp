#include <iostream>
#include <string>
#include <vector>

#include "orcs/os_adapter.hpp"
#include "orcs/process_discovery.hpp"
#include "orcs/process_terminator.hpp"
#include "orcs/system_telemetry.hpp"

int main(int argc, char** argv) {
  std::string banned = argc > 1 ? argv[1] : "firefox";

  orcs::SystemTelemetry telemetry;
  const auto snapshot = telemetry.sample();

  orcs::ProcessDiscovery discovery;
  const auto matches = discovery.list_by_name(banned);

  std::cout << "platform=" << snapshot.os_type << " total_mem=" << snapshot.total_mem_bytes
            << " free_mem=" << snapshot.free_mem_bytes << " matched_processes=" << matches.size()
            << "\n";

  orcs::ProcessTerminator terminator;
  for (const auto& proc : matches) {
    std::cout << "restricted process found: pid=" << proc.pid << " name=" << proc.name << "\n";
    (void)terminator;
  }

  return 0;
}
