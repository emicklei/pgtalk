package xs

type ReadWrite interface {
	Name() string
	Value(entity interface{}, fieldValue interface{})
}
type Int8Access struct {
	name   string
	writer func(dest interface{}, i *int64)
}

func NewInt8Access(name string, writer func(dest interface{}, i *int64)) Int8Access {
	return Int8Access{name: name, writer: writer}
}

func (a Int8Access) Value(entity interface{}, fieldValue interface{}) {
	var i int64 = fieldValue.(int64)
	a.writer(entity, &i)
}

type ScanToWrite struct {
	RW     ReadWrite
	Entity interface{}
}

func (s ScanToWrite) Scan(fieldValue interface{}) error {
	s.RW.Value(s.Entity, fieldValue)
	return nil
}

func (a Int8Access) Name() string { return a.name }

type TextAccess struct {
	name   string
	writer func(dest interface{}, i *string)
}

func NewTextAccess(name string, writer func(dest interface{}, i *string)) TextAccess {
	return TextAccess{name: name, writer: writer}
}

func (a TextAccess) Value(entity interface{}, fieldValue interface{}) {
	var i string = fieldValue.(string)
	a.writer(entity, &i)
}

func (a TextAccess) Name() string { return a.name }
