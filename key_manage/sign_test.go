package key_manage

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestSignMessage(t *testing.T) {
	message, _ := hex.DecodeString("a90abb62e8c7614e8f8af810083db2b4d1ddb47c2839457ea539fd232feaeed3")
	signature, err := Sign("aead75071f4a9437df36d08acdcbb78b8dca55d02f0631f33f72274e9ee45a98", message, false)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("signature: ", signature)

	if signature != "168d885cf992d7cc1b7831a65ae184b23f21548897477cf9b452243baa21bbf1216734d953a549c7b5014e33cfd6d84d4c0df1cae80c0725a12e884ac442a52e1c" {
		t.Fatal("sign message failed")
	}
}
