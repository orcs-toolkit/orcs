#pragma once

#include <string>
#include <vector>

namespace orcs {

struct ProcessInfo {
  int pid;
  std::string name;
};

class ProcessDiscovery {
 public:
  std::vector<ProcessInfo> list_by_name(const std::string& process_name) const;
};

}  // namespace orcs
