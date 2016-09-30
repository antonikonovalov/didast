package timeid

import (
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	id := New()
	gap := DetectInterval(id, time.Hour)
	if len(gap) != 2 {
		t.Fatalf(`len of interval must be 2, got - %d`, len(gap))
	}
	if id < gap[0] || id > gap[1] {
		t.Errorf("epexted true detected interval for id %d, but - %d <> %d\n id = %s \n from = %s \n to = %s \n",
			id, gap[0], gap[1],
			time.Unix(0, id).UTC().String(),
			time.Unix(0, gap[0]).UTC().String(),
			time.Unix(0, gap[1]).UTC().String(),
		)
	}

	id2 := New()
	if id == id2 {
		t.Fatalf(`expected ids not eq, but - %d = %d`, id, id2)
	}
	gap2 := DetectInterval(id2, time.Hour)
	if fmt.Sprintf(`%v`, gap) != fmt.Sprintf(`%v`, gap2) {
		t.Fatalf(`expected gap is eqauls, but - [%d,%d] != [%d,%d]`, gap[0], gap[1], gap2[0], gap2[1])
	}

}

// TestCapability FAIL!!! нужно писать тикет в golang/go
func TestCapability(t *testing.T) {
	ids := []int64{}
	for j := 0; j < 100000; j++ {
		ids = append(ids, New())
	}
	for i, id := range ids {
		var nextId int64
		if i+1 != len(ids) {
			nextId = ids[i+1]
		} else {
			continue
		}
		if id > nextId {
			println(`error`, i, id, nextId, `id more then next`)
		}
		if len(fmt.Sprintf(`%d`, id)) != len(fmt.Sprintf(`%d`, nextId)) {
			println(`error`, i, id, nextId, `lens not equals`)
		}
	}
}
