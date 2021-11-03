package main

var endOfPosting PagePostingItem = PagePostingItem{
	PID: lastPID,
}

type PagePostingItem struct {
	PID    int64
	Weight float32
}

type PagePostingList struct {
	// the end of posting items must be endOfPostingItem
	Items []PagePostingItem
	// meta info
	MaxWeight float32
}

type Feature string // 特征名称，包含特征值？

type FeatureCursor struct {
	feature Feature
	posting *PagePostingList
	current int
}

func (c *FeatureCursor) CurrentPage() PagePostingItem {
	return c.posting.Items[c.current]
}

func (c *FeatureCursor) Beyond(PID int64) {
	// page IDs are sorted in non-decreasing order
	for c.CurrentPage().PID < PID {
		c.current++
	}
}
