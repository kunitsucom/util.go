package cockroachdb

type Index struct {
	Name    *Ident
	Columns []*ColumnIdent
}

func (index *Index) GetPlainName() string {
	return index.Name.PlainString()
}

type Indexes []*Index

func (indexes Indexes) Append(index *Index) Indexes {
	for i := range indexes {
		if indexes[i].GetPlainName() == index.GetPlainName() {
			indexes[i] = index
			return indexes
		}
	}
	return append(indexes, index)
}
