package config

// 端口转发服务
type PortForward struct {
	Listen  string   `json:"listen"`
	Forward []string `json:"forward"`
}

// DNS 代理服务
type DnsProxy struct {
	DnsServer string `json:"dns_server"`
}

// 端口连接检测
type PortConnCheck struct {
	Ports string `json:"ports"`
}

// 默认设置
type Config struct {
	NodeName      string        `json:"node_name"`
	NodeIp        string        `json:"node_ip"`
	PortForward   []PortForward `json:"port_forward"`
	DnsProxy      DnsProxy      `json:"dns_proxy"`
	PortConnCheck PortConnCheck `json:"port_conn_check"`
}
