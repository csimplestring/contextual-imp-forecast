package main

import "sort"

type ImpressionForecastor struct {
	pages          map[int64]PageMeta
	featureToPages map[Feature]*PagePostingList
}

func (i *ImpressionForecastor) forecast(ad Ad) int {
	impressions := 0

	searcher := &Searcher{
		pages:         i.pages,
		feature2Pages: i.featureToPages,
		currentPage:   0,
		ad:            ad,
	}

	searcher.initCursors()
	pageID := searcher.nextCandidate()
	for pageID < lastPID {
		p := i.pages[pageID]
		score := i.score(ad, p)
		if score > p.minScore {
			impressions += p.impressions
		}
		pageID = searcher.nextCandidate()
	}

	return impressions
}

func (i *ImpressionForecastor) score(ad Ad, page PageMeta) float32 {
	pageWeight := make(map[Feature]float32, len(page.features))
	for i, f := range page.features {
		pageWeight[f] = page.weights[i]
	}

	sim := float32(0)
	for _, f := range ad.Features() {
		sim += pageWeight[f] * ad.Weight(f)
	}

	return sim * ad.Bid()
}

type Searcher struct {
	pages         map[int64]PageMeta
	feature2Pages map[Feature]*PagePostingList

	currentPage int64
	ad          Ad
	features    []*FeatureCursor
}

func (s *Searcher) initCursors() {
	s.currentPage = 0

	ad := s.ad
	features := ad.Features()
	cursors := make([]*FeatureCursor, len(features))
	for i, f := range features {

		cursors[i] = &FeatureCursor{
			feature: f,
			posting: s.feature2Pages[f],
			current: 0,
		}
	}

	s.features = cursors
}

func (s *Searcher) nextCandidate() int64 {
	for {

		s.sortFeatures(s.features)
		pivot := s.findPivotFeature(s.features)

		if pivot == -1 {
			return lastPID
		}
		pivotPID := s.features[pivot].CurrentPage().PID
		if pivotPID == lastPID {
			return lastPID
		}

		if pivotPID == s.currentPage {
			f := s.pickFeature(s.features[0 : pivot+1])
			s.features[f].Beyond(pivotPID + 1)
		} else {
			if s.features[0].CurrentPage().PID == pivotPID {
				s.currentPage = pivotPID
				return s.currentPage
			} else {
				f := s.pickFeature(s.features[0 : pivot+1])
				s.features[f].Beyond(pivotPID)
			}
		}
	}
}

func (s *Searcher) sortFeatures(features []*FeatureCursor) {
	sort.Slice(features, func(i, j int) bool {
		return features[i].CurrentPage().PID < features[j].CurrentPage().PID
	})
}

func (s *Searcher) findPivotFeature(cursors []*FeatureCursor) int {
	n := len(cursors)
	ub := float32(0)
	ad := s.ad

	minScoreP := s.pages[s.currentPage].minScore
	threshold := minScoreP / ad.Bid()

	for i := 0; i < n; i++ {

		f := cursors[i].feature
		maxWeight := s.feature2Pages[f].MaxWeight
		ub += ad.Weight(f) * maxWeight
		if ub >= threshold {
			return i
		}
	}

	return -1
}

func (s *Searcher) pickFeature(cursors []*FeatureCursor) int {
	return 0
}
