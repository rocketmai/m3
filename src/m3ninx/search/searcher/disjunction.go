// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package searcher

import (
	"github.com/m3db/m3/src/m3ninx/index"
	"github.com/m3db/m3/src/m3ninx/postings"
	"github.com/m3db/m3/src/m3ninx/postings/roaring"
	"github.com/m3db/m3/src/m3ninx/search"
)

type disjunctionSearcher struct {
	searchers search.Searchers
}

// NewDisjunctionSearcher returns a new Searcher which matches documents which are matched
// by any of the given Searchers.
func NewDisjunctionSearcher(searchers search.Searchers) (search.Searcher, error) {
	if len(searchers) == 0 {
		return nil, errEmptySearchers
	}

	return &disjunctionSearcher{
		searchers: searchers,
	}, nil
}

func (s *disjunctionSearcher) Search(r index.Reader) (postings.List, error) {
	var (
		union = make([]postings.List, 0, len(s.searchers))
	)
	for _, sr := range s.searchers {
		pl, err := sr.Search(r)
		if err != nil {
			return nil, err
		}

		union = append(union, pl)
	}

	// Perform a lazy fast union.
	// TODO: Try and see if returns err, if so fallback to slower method?
	return roaring.UnionReadOnly(union)
}
