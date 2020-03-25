// Copyright 2019 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package builtin

import (
	"log"
	"math/rand"

	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/model"

	"github.com/pingcap/tipocket/pkg/go-sqlsmith/types"
	"github.com/pingcap/tipocket/pkg/go-sqlsmith/util"
)

type hintClass struct {
	name     string
	minArg   int
	maxArg   int
	constArg bool
	mysql    bool
	stable   bool
}

type hintData struct {
	data interface{}
}

var hintKeywords = []*hintClass{
	// with no args
	{"hash_agg", 0, 0, false, false, false},
	{"stream_agg", 0, 0, false, false, false},
	{"agg_to_cop", 0, 0, false, false, false},
	{"read_consistent_replica", 0, 0, false, false, false},
	{"no_index_merge", 0, 0, false, false, false},
	{"qb_name", 0, 0, false, false, false},

	// these have been renamed
	// {"tidb_hj", 2, 3, false, false, true},
	// {"tidb_smj", 2, 3, false, false, true},
	// {"tidb_inlj", 2, 3, false, false, true},
	// with 2 or more args
	{"hash_join", 2, -1, false, true, false},
	{"merge_join", 2, -1, false, false, false},
	{"inl_join", 2, -1, false, false, false},

	// with table name and at least one idx name
	{"use_index", 2, -1, false, false, false},
	{"ignore_index", 2, -1, false, false, false},
	{"use_index_merge", 2, -1, false, false, false},

	// with bool (TRUE or FALSE)
	{"use_toja", 1, 1, false, false, false},
	{"enable_plan_cache", 1, 1, false, false, false},
	{"use_cascades", 1, 1, false, false, false},

	// with int (MB)
	{"memory_quota", 1, 1, false, false, false},
	// with int (ms)
	{"max_execution_time", 1, 1, false, false, false},
}

// these will not be generated for some reason
var disabledHintKeywords = []*hintClass{
	// not released?
	{"time_range", 2, -1, false, false, false},
	// storage type with tablename: TIKV[t1]
	{"read_from_storage", 2, -1, false, false, false},
	// not released?
	{"query_type", 1, 1, false, false, false},
}

func GenerateHintExpr(table *types.Table) (h *ast.TableOptimizerHint) {
	rd := util.Rd(1000) % len(hintKeywords)
	hintKeyword := hintKeywords[rd]
	h.HintName = model.NewCIStr(hintKeyword.name)

	if hintKeyword.maxArg == 0 {
		return
	}

	if hintKeyword.maxArg == 1 {
		switch hintKeyword.name {
		case "use_toja", "enable_plan_cache", "use_cascades":
			h.HintData = hintData{
				data: util.Rd(2)%2 == 1,
			}
		case "memory_quota":
			h.HintData = hintData{
				data: util.Rd(100) + 50,
			}
		case "max_execution_time":
			h.HintData = hintData{
				data: util.Rd(1000) + 500,
			}
		default:
			log.Fatalf("unreachable hintKeyword.name:%s", hintKeyword.name)
		}
		return
	}

	shuffledTables := make([]ast.HintTable, len(h.Tables))
	copy(h.Tables, shuffledTables)
	rand.Shuffle(len(shuffledTables), func(i, j int) {
		shuffledTables[i], shuffledTables[j] = shuffledTables[j], shuffledTables[i]
	})

	shuffledIndexes := make([]model.CIStr, len(h.Indexes))
	copy(h.Indexes, shuffledIndexes)
	rand.Shuffle(len(shuffledIndexes), func(i, j int) {
		shuffledIndexes[i], shuffledIndexes[j] = shuffledIndexes[j], shuffledIndexes[i]
	})

	switch hintKeyword.name {
	case "hash_join", "merge_join", "inl_join":
		n := util.MinInt(util.Rd(4)+1, len(shuffledTables)) // avoid case n == 0
		for ; n > 0; n-- {
			h.Tables = append(h.Tables, shuffledTables[n-1])
		}
	case "use_index", "ignore_index", "use_index_merge":
		h.Tables = append(h.Tables, shuffledTables[util.Rd(len(shuffledTables))])
		n := util.MinInt(util.Rd(4)+1, len(shuffledIndexes)) // avoid case n == 0
		for ; n > 0; n-- {
			h.Indexes = append(h.Indexes, shuffledIndexes[n-1])
		}
	default:
		log.Fatalf("unreachable hintKeyword.name:%s", hintKeyword.name)
	}
	return
}
