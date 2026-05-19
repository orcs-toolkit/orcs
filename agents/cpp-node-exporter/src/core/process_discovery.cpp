#include "orcs/process_discovery.hpp"

#include <algorithm>
#include <filesystem>
#include <fstream>

namespace orcs {

std::vector<ProcessInfo> ProcessDiscovery::list_by_name(const std::string& process_name) const {
  std::vector<ProcessInfo> out;
#ifdef __linux__
  for (const auto& entry : std::filesystem::directory_iterator("/proc")) {
    if (!entry.is_directory()) {
      continue;
    }
    const auto name = entry.path().filename().string();
    if (!std::all_of(name.begin(), name.end(), ::isdigit)) {
      continue;
    }
    std::ifstream cmdline(entry.path() / "comm");
    std::string comm;
    if (!cmdline.good() || !std::getline(cmdline, comm)) {
      continue;
    }
    if (comm == process_name) {
      out.push_back(ProcessInfo{std::stoi(name), comm});
    }
  }
#else
  (void)process_name;
#endif
  return out;
}

}  // namespace orcs
