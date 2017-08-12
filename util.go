package hub

func (h *Hub) clamp(Q0, Q1 int64) (int64, int64) {
	nr := h.Buffer.Len()
	return clamp(Q0, 0, nr), clamp(Q1, 0, nr)
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func clamp(v, l, h int64) int64 {
	if v < l {
		return l
	}
	if v > h {
		return h
	}
	return v
}
