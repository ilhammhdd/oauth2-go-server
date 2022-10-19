package test

import (
	"encoding/base64"
	"encoding/json"
	"math"
	"math/bits"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ilhammhdd/go-toolkit/errorkit"
	"golang.org/x/crypto/blake2b"
	"ilhammhdd.com/oauth2-go-server/entity"
)

type TryJsonTime struct {
	DateTime time.Time `json:"date_time"`
}

func TestInitiateRegisterClient(t *testing.T) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		t.Fatalf("uuid new random error: %s", err.Error())
	}
	uuidBlake2b256 := blake2b.Sum256([]byte(uuid.String()))
	t.Logf("\nuuid: %s, base64 uuidBlake2b256: %s", uuid, base64.RawURLEncoding.EncodeToString(uuidBlake2b256[:]))

	icr, detailedErr := entity.NewInitiateClientRegistration(nil)
	if errorkit.IsNotNilThenLog(detailedErr) {
		t.Fatalf("\n%s", detailedErr.Error())
	}
	t.Logf("\nsession expired at zero-value RFC3339nano: %s", icr.SessionExpiredAt.Format(time.RFC3339Nano))

	icr.SessionExpiredAt = time.Now().UTC()
	t.Logf("\nsession expired at UTC RFC3339nano: %s", icr.SessionExpiredAt.Format(time.RFC3339Nano))

	tjt := TryJsonTime{DateTime: time.Now().UTC()}
	tjtJson, err := json.Marshal(tjt)
	if err != nil {
		t.Fatalf("\nerror marshalling TryJsonTime: %s", err.Error())
	}
	t.Logf("\nTryJsonTime Json: %s", tjtJson)

	tjtJsonTest := `{"date_time": "2022-09-01T18:00:00Z"}`
	var tjtUnmarshaled TryJsonTime
	err = json.Unmarshal([]byte(tjtJsonTest), &tjtUnmarshaled)
	if err != nil {
		t.Fatalf("\nerror unmarshalling TryJsonTime: %s", err.Error())
	}
	t.Logf("\nTryJsonTime struct: %v", tjtUnmarshaled)
}

func TestBitsArithmatic(t *testing.T) {
	var left1 uint32 = 4_294_967_294
	var right1 uint32 = 2
	var carry1 uint32 = 0
	sum1, carryOut1 := bits.Add32(left1, right1, carry1)
	t.Logf("\nsum1: %d, carryOut1: %d", sum1, carryOut1)

	var seedLeft2 float32 = 20.2
	var seedRight2 float32 = 40.4
	var left2 uint32 = math.Float32bits(seedLeft2)
	var right2 uint32 = math.Float32bits(seedRight2)
	var carry2 uint32 = 0
	sum2, carryOut2 := bits.Add32(left2, right2, carry2)
	var sum2Float float32 = math.Float32frombits(sum2)
	t.Logf("\nsum2 bits: %d, sum2 float: %f, carryOut2: %d", sum2, sum2Float, carryOut2)
}
