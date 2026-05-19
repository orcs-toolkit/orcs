#pragma once

namespace orcs {

class ProcessTerminator {
 public:
  bool kill_pid(int pid) const;
};

}  // namespace orcs
