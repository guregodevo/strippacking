#include "strip_packing.h"

#define PRNG_BUFSZ 128

void Context::ThreadInit() {
  char* statebuf = new char[PRNG_BUFSZ];
  initstate_r(seed, statebuf, PRNG_BUFSZ, &__random_data);
}

Context::Context(int argc, char** argv) {
  seed = rand();
  opt.Parse(argc, argv);
}

void Context::InitAlgo() {
  if ("kp1" == opt.algo) {
    algo = new Kp1Algo;
  } else if ("kp2_msp_b" == opt.algo) {
    algo = new Kp2MspBalanced;
  } else if ("pyramid" == opt.algo) {
    algo = new PyramidAlgo;
  } else if ("simple_pyramid" == opt.algo) {
    algo = new SimplePyramidAlgo;
  }
}

void Context::DestroyAlgo() {
  delete algo;
}

void Context::Render() {
  if (!opt.render) {
    return;
  }
  ren = new Renderer;
  ren->ShowAll(this);
}

struct Segment {
  double x, y;
  Segment(double x, double y) : x(x), y(y) {}
};

bool point_inside_segment(double a, Segment* b) {
  return double_less(b->x, a) && double_less(a, b->y);
}

bool segments_overlap(Segment* a, Segment* b) {
  return point_inside_segment(a->x, b) ||
    point_inside_segment(a->y, b) ||
    point_inside_segment(b->x, a) ||
    point_inside_segment(b->y, a);
}

bool rects_overlap(Rect* a, Rect* b) {
  Segment s1(a->x, a->x + a->w), s2(b->x, b->x + b->w);
  Segment s3(a->y, a->y + a->h), s4(b->y, b->y + b->h);
  return segments_overlap(&s1, &s2) && segments_overlap(&s3, &s4);
}

bool rect_inside_strip(Rect* r, int m, double H) {
  if (double_less(r->x, 0)
      || double_less(double(m), r->x + r->w) 
      || double_less(r->y, 0) 
      || double_less(H, r->y + r->h)) {
    return false;
  }
  for (int y = 1; y < m; ++y) {
    if (double_less(r->x, double(y)) && double_less(double(y), r->x + r->w)) {
      return false;
    }
  }
  return true;
}

bool Context::Validate() {
  if (!opt.validate) {
    return true;
  }
  std::vector<SavedRect>::iterator i;
  for (i = algo->saved_rects.begin(); i != algo->saved_rects.end(); ++i) {
    if (i->color == "red") {
      continue;    
    }
    if (!rect_inside_strip(&i->r, opt.m, algo->solution_height)) {
      std::cout << "FAIL\n";
      std::cout << "Following rectangle intersects strip borders\n";
      i->r.PrintInfo();
      return false;
    }
  }
  for (int y = 0; y < algo->saved_rects.size() - 1; ++y) {
    if (algo->saved_rects[y].color == "red") {
      continue;
    }
    for (int j = y + 1; j < algo->saved_rects.size(); ++j) {
      if (algo->saved_rects[j].color == "red") {
        continue;
      }
      if (rects_overlap(&algo->saved_rects[y].r, &algo->saved_rects[j].r)) {
        std::cout << "FAIL\n";
        std::cout << "Following 2 rectangles overlap\n";
        algo->saved_rects[y].r.PrintInfo();
        algo->saved_rects[j].r.PrintInfo();
        return false;
      }
    }
  }
  return true;
}