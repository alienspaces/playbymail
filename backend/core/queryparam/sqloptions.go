package queryparam

import (
	"time"

	coresql "gitlab.com/alienspaces/playbymail/core/sql"
)

func ToSQLOptionsWithDefaults(qp *QueryParams) *coresql.Options {
	if len(qp.SortColumns) == 0 {
		qp.SortColumns = []SortColumn{
			{
				Col: "created_at",
			},
		}
	}
	return ToSQLOptions(qp)
}

func ToSQLOptions(qp *QueryParams) *coresql.Options {
	opts := &coresql.Options{
		Limit:  qp.PageSize + 1,
		Offset: (qp.PageNumber - 1) * qp.PageSize,
	}

	for _, sc := range qp.SortColumns {
		dir := coresql.OrderDirectionASC
		if sc.IsDescending {
			dir = coresql.OrderDirectionDESC
		}

		opts.OrderBy = append(opts.OrderBy, coresql.OrderBy{
			Col:       sc.Col,
			Direction: dir,
		})
	}

	for col, values := range qp.Params {
		qOpVals := map[Operator][]string{}
		var sqlParams []coresql.Param

		var sqlOp coresql.Operator
		for _, p := range values {
			if v, ok := p.Val.([]string); ok {
				switch p.Op {
				case OpLike:
					sqlOp = coresql.OpLikeAny
				case OpILike:
					sqlOp = coresql.OpILikeAny
				}

				for i, s := range v {
					v[i] = "%" + s + "%"
				}
				p.Val = v
			} else if v, ok := p.Val.(string); ok {
				switch p.Op {
				case "":
					// "" corresponds to =, IN, ANY, @>
					qOpVals[p.Op] = append(qOpVals[p.Op], v)
					continue
				case OpGreaterThan:
					sqlOp = coresql.OpGreaterThan
				case OpGreaterThanEqual:
					sqlOp = coresql.OpGreaterThanEqual
				case OpLessThan:
					sqlOp = coresql.OpLessThan
				case OpLessThanEqual:
					sqlOp = coresql.OpLessThanEqual

					if _, err := time.Parse(time.DateOnly, v); err == nil {
						// This is necessary because postgres maps a date of
						// 2022-01-01 to 2022-1-01T00:00:00Z
						v += "T23:59:59.999999999Z"
					}
					p.Val = v
				case OpNotEqual:
					sqlOp = coresql.OpNotEqual
				case OpLike:
					sqlOp = coresql.OpLike
					p.Val = "%" + v + "%"
				case OpILike:
					sqlOp = coresql.OpILike
					p.Val = "%" + v + "%"
				}
			}

			sqlParams = append(sqlParams, coresql.Param{
				Col: col,
				Op:  sqlOp,
				Val: p.Val,
			})
		}

		for op, v := range qOpVals {
			if op == "" {
				var val any
				val = v
				if len(v) == 1 {
					val = v[0]
				}

				sqlParams = append(sqlParams, coresql.Param{
					Col: col,
					Val: val,
				})
			}
		}
		opts.Params = append(opts.Params, sqlParams...)
	}

	return opts
}
