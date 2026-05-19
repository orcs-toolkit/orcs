#pragma once

#include <string>

namespace orcs {

class OSAdapter {
 public:
  virtual ~OSAdapter() = default;
  virtual std::string platform_name() const = 0;
};

class LinuxAdapter : public OSAdapter {
 public:
  std::string platform_name() const override;
};

class WindowsAdapter : public OSAdapter {
 public:
  std::string platform_name() const override;
};

class MacOSAdapter : public OSAdapter {
 public:
  std::string platform_name() const override;
};

}  // namespace orcs
