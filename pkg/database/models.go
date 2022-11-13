package database

type InputModel struct {
	Mqtt      string `csv:"mqtt"`
	Invid     string `csv:"invid"`
	UnitGuid  string `csv:"unit_guid"`
	MsgId     string `csv:"msg_id"`
	Text      string `csv:"text"`
	Context   string `csv:"context"`
	Class     string `csv:"class"`
	Level     int    `csv:"level"`
	Area      string `csv:"area"`
	Addr      string `csv:"addr"`
	Block     string `csv:"block"`
	DataType  string `csv:"type"`
	Bit       string `csv:"bit"`
	InvertBit string `csv:"invert_bit"`
}
