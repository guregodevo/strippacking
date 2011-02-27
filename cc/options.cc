#include <cstdlib>

#include "strip_packing.h"

void Options::Parse(int argc, char** argv) {
  for (int y = 1; y < argc; y++) {
    if (0 == strcmp("-n", argv[y])) {
      ++y;
      n = atoi(argv[y]);
    } else if (0 == strcmp("-m", argv[y])) {
      ++y;
      m = atoi(argv[y]);
    } else if (0 == strcmp("-r", argv[y])) {
      render = true;
    } else if (0 == strcmp("-a", argv[y])) {
      ++y;
      algo = argv[y];
    } else if (0 == strcmp("-v", argv[y])) {
      validate = true;
    } else if (0 == strcmp("-t", argv[y])) {
      ++y;
      t = atoi(argv[y]);
    }
  }
  save_rects = render || validate;
}