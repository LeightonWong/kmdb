package kmdb

type Command struct {
	Type  string `json:"type"`
	Key   []byte `json:"key"`
	Value []byte `json:"value"`
	Sync  bool   `json:"sync"`
}

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Rst  []byte `json:"rst"`
}
