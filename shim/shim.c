#include "shim.h"

// Anchor symbol so the dylib always has at least one translation unit
// besides the linked Go archive. Public SPACESEE_mock_* symbols come from the Go c-archive.
void spacesee_mock_shim_anchor(void) {}
