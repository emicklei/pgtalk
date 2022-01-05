package things

import (
	"log"
	"testing"
	"time"

	"github.com/emicklei/pgtalk"
	"github.com/emicklei/pgtalk/convert"
	"github.com/jackc/pgtype"
)

type structAccess struct {
	ID    pgtalk.FieldAccess[pgtype.UUID]
	Tdate pgtalk.FieldAccess[pgtype.Date]
}

var ThingAccess = structAccess{
	ID:    ID,
	Tdate: Tdate,
}

func TestStructAccess(t *testing.T) {
	ex := ThingAccess.Tdate.Less(convert.TimeToDate(time.Now()))
	log.Println(pgtalk.SQL(ex))
}
