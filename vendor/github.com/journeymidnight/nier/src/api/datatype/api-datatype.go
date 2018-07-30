package datatype

const (
	ROLE_ROOT         = "ROOT"
	ROLE_ADMIN        = "ADMIN"
	ROLE_USER         = "USER"
	OP_TYPE_SET_SNMP  = "SET SNMP"
	OP_TYPE_SET_NTP   = "SET NTP"
	REQUEST_TOKEN_KEY = "TOKEN"
)

type QueryRequest struct {
	Name        string `json:"name,omitempty"`
	Password    string `json:"password,omitempty"`
	Type        string `json:"type,omitempty"`
	SnmpAddress string `json:"snmpaddress,omitempty"`
	HostName    string `json:"hostname,omitempty"`
	DiskName    string `json:"diskname,omitempty"`
}

type QueryRequestNtp struct {
	Name       string `json:"name,omitempty"`
	Password   string `json:"password,omitempty"`
	Type       string `json:"type,omitempty"`
	NtpAddress string `json:"ntpaddress,omitempty"`
}
type QueryResponse struct {
	RetCode int         `json:"retCode"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type UserRecord struct {
	UserName string
	Password string
	Type     string
}

type Roles struct {
	Roles []string
}

type TokenRecord struct {
	Token    string `json:"token"`
	UserName string `json:"userName"`
	Type     string `json:"type"`
	Created  string `json:"created"`
	Expired  string `json:"expired"`
}

type SnmpRecord struct {
	SnmpAddress string `json:"snmpAddress"`
}

type NtpRecord struct {
	NtpAddress string `json:"ntpAddress"`
}

type HistoryRecord struct {
	UserName string `json:"userName"`
	OpType   string `json:"opType"`
	Time     string `json:"time"`
}

type AlertRecord struct {
	EventId   string `json:"eventid"`
	EventName string `json:"eventname"`
	Node      string `json:"node"`
	SendTime  string `json:"sendtime`
	Status    string `json:"status`
}

type CreateViewRecords struct {
	ExportedFS string `json:"exportedfs"`
}

type RemoveViewRecord struct {
	ExportedFS string `json:"exportedfs"`
}

type VipRecord struct {
	Vip string `json:"vip"`
}