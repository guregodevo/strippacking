#include "strip_packing.h"

void PyramidAlgo::Pack(int n, double xbe, double ybe, Context* context) {
  InitParams(n);
  
  Rect r;
  std::set<Rect> s[2];
  PyramidPos p[2];
  int best_ind;
  double best_pos;
  
  for (int y = 0; y < n; ++y) {
    NextRect(&r);
    
    p[0] = PackToPyramid(&s[0], &r);
    p[1] = PackToPyramid(&s[1], &r);
    if (p[0].pos > p[1].pos) {
      best_ind = 0;
      best_pos = p[0].pos;
    } else {
      best_ind = 1;
      best_pos = p[1].pos;
    }
    if (double_eq(0, best_pos)) {
      PackOnTop(&r);
    } else {
      PerformPacking(&s[best_ind], &p[best_ind], &r);
      ConvertCooToComplPyramid(&r, best_ind);
    }
    
    if (context->opt.save_rects) {
      SaveRect(&r);
    }
    RecalcSolutionHeightSingle(&r);
  }
  if (context->opt.render_bins) {
    saved_rects.push_back(
        SavedRect(Rect(0, 0, 1, h_ + 2 * shift_), "red", false));
  }
}

void PyramidAlgo::InitParams(int n) {
  shift_ = powl(n, 0.5);
  h_ = n / 4;
  top_h_ = 0;
}

void AppendAtBottom(std::set<Rect>* s, std::set<Rect>::iterator* i, Rect* r) {
  Rect b = **i;
  s->erase(*i);
  b.h += r->h;
  b.y -= r->h;
  r->x = 0;
  r->y = b.y;
  s->insert(b);
}

void AppendAtTop(std::set<Rect>* s, std::set<Rect>::iterator* i, Rect* r) {
  Rect b = **i;
  s->erase(*i);
  r->x = 0;
  r->y = b.y + b.h;
  b.h += r->h;
  s->insert(b);
}

void PyramidAlgo::PerformPacking(std::set<Rect>* s, PyramidPos* p, Rect* r) {
  if (s->end() == p->i) {
    r->x = 0;
    r->y = p->pos - r->h;
    s->insert(*r);
    return;
  }
  if (double_eq(p->pos, p->i->y)) {
    AppendAtBottom(s, &p->i, r);
  } else {
    AppendAtTop(s, &p->i, r);
  }
}

// Returns iterator pointing to pyramid's rectangle and position of top edge of
// rectangle $r. Does not perform actual packing.
PyramidPos PyramidAlgo::PackToPyramid(std::set<Rect>* s, Rect* r) {
  std::set<Rect>::iterator i, j;
  double cur = (1 - r->w) * h_ + shift_;
  double be, en;

  for (i = s->lower_bound(Rect(0, cur, 0, 0)); ; ++i) {
    if (s->end() == i) {
      en = 0; // No bottom limit.
    } else {
      en = i->y + i->h;
    }
    j = i;
    if (s->begin() == j) {
      be = h_ + shift_; // No limit at top.
    } else {
      --j;
      be = j->y;
    }
    be = std::min<double>(be, cur);
    
    // [be, en] boundaries determined.
    if (double_less(r->h, be - en)) {
      if ((j != i) && double_eq(j->y, be)) {
        return PyramidPos(j, be);
        //AppendAtBottom(s, &j, r);
        //return true;
      }
      if ((i != s->end()) && double_eq(i->y + i->h, en) 
          && double_less(be - en, shift_)) {
        return PyramidPos(i, en + r->h);
        //AppendAtTop(s, &i, r);
        //return true;
      }
      return PyramidPos(s->end(), be);
      //r->x = 0;
      //r->y = be - r->h;
      //s->insert(*r);
      //return true;
    }
    if (s->end() == i) {
      break;
    }
  }
  return PyramidPos(s->end(), 0);
}

bool operator< (const Rect& a, const Rect& b) {
  return (a.y + a.h) > (b.y + b.h);
}

void PyramidAlgo::ConvertCooToComplPyramid(Rect* r, int ind) {
  if (0 == ind) {
    return;
  }
  r->x = 1 - r->w;
  r->y = h_ + 2 * shift_ - r->y - r->h;
}

void PyramidAlgo::PackOnTop(Rect* r) {
  r->x = 0;
  r->y = top_h_ + h_ + 2 * shift_;
  top_h_ += r->h;
}