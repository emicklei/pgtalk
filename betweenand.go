package pgtalk

type BetweenAnd struct {
}

func MakeBetweenAnd(reader ColumnAccessor, begin, end SQLWriter) BetweenAnd { return BetweenAnd{} }
