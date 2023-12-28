package utils

import "encoding/json"

func DeepCopy(src, dst any) {
	bytes, _ := json.Marshal(src)
	_ = json.Unmarshal(bytes, dst)
}
