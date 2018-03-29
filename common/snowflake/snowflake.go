package snowflake

import (
	"github.com/imkuqin-zw/uuid_generator/common"
)

const (
	TS_MASK         = 0x1FFFFFFFFFF // 41bit
	MACHINE_ID_MASK = 0x3FF         // 10bit
	SN_MASK         = 0xFFF         // 12bit
)

func CreateUUID(chProc chan chan uint64, machineID uint64) int64 {
	var sn uint64     // 12-bit serial no
	var last_ts int64 // last timestamp
	for {
		ret := <-chProc
		t := common.Ts()
		if t < last_ts {
			t = common.WaitMs(last_ts)
		}
		if last_ts == t {
			sn = (sn + 1) & SN_MASK
			if sn == 0 {
				t = common.WaitMs(last_ts)
			}
		} else {
			sn = 0
		}
		last_ts = t
		ret <- (uint64(t) & TS_MASK) << 22 | machineID | sn
	}
}