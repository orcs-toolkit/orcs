#include "orcs/process_terminator.hpp"

#ifdef __linux__
#include <signal.h>
#endif

namespace orcs {

bool ProcessTerminator::kill_pid(int pid) const {
#ifdef __linux__
  return ::kill(pid, SIGKILL) == 0;
#else
  (void)pid;
  return false;
#endif
}

}  // namespace orcs
