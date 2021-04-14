package dnsProxy

import (
	"HuanZhen/logger"
	"bufio"
	"fmt"
	"github.com/miekg/dns"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var blackDomain []string

func DnsServer(laddr string) error {

	serveMux := dns.NewServeMux()
	serveMux.HandleFunc(".", handleDnsRequest)

	e := make(chan error)
	for _, _net := range [...]string{"udp", "tcp"} {
		srv := &dns.Server{Addr: laddr, Net: _net, Handler: serveMux}
		go func(srv *dns.Server) {
			e <- srv.ListenAndServe()
		}(srv)
	}
	return <-e
}

func handleDnsRequest(w dns.ResponseWriter, req *dns.Msg) {
	dnsService := "114.114.114.114:53" // 默认的域名服务器
	//ip := "127.0.0.1"                     // 默认返回的IP地址
	// 处理域名请求

	// 首先记录请求的域名
	// 判断请求的域名是否在 domain cache 中
	//    是 -> 直接返回 cache 中的内容
	//    否 -> 判断域名是否在恶意域名中，包括广告域名、指定域名
	//        否 -> 请求正常的域名解析，返回内容
	//        是 -> 返回 DNS 污染
	// 如果不是，那么就让其请求 114.114.114.114， 或者自己请求一下，将结果返回


	quesFqdn := req.Question[0].Name
	domain := quesFqdn[:len(quesFqdn)-1]
	fmt.Println("正在请求的域名：", domain)

	if MatchDomain(domain) {
		logger.HZLogger.Warn("正在请求恶意域名", domain)
	}else {
		logger.HZLogger.Info("正在请求的域名：", domain)
	}

	server := domain + "."
	c := dns.Client{Timeout: 5 * time.Second}
	m, _, err := c.Exchange(req, dnsService)

	if err != nil {
		log.Println("DNS 请求失败")
	} else {
		// A 记录
		if m.Question[0].Qtype == 1 {
			if m.Question[0].Name == server {
				n := len(m.Answer)
				for i := 0; i < n; i++ {
					s := strings.Fields(m.Answer[i].String())
					//s[4] = ip
					m.Answer[i], _ = dns.NewRR(strings.Join(s, "\t"))
				}
			}
		}

		// CN 记录
		if m.Question[0].Qtype == 5 { //=5 表示CN记录
			if m.Question[0].Name == server {
				n := len(m.Answer)
				for i := 0; i < n; i++ {
					s := strings.Fields(m.Answer[i].String())
					if s[3] == "A" {
						//s[4] = ip
						m.Answer[i], _ = dns.NewRR(strings.Join(s, "\t"))
					}
				}
			}
		}

		// MX 记录
		if m.Question[0].Qtype == 15 { //15 表示MX记录
			n := len(m.Answer)
			for i := 0; i < n; i++ {
				s := strings.Fields(m.Answer[i].String())
				s[5] = server
				m.Answer[i], _ = dns.NewRR(strings.Join(s, "\t"))
			}
		}
		w.WriteMsg(m)
	}
}

func StartDnsProxy() {
	// 读取恶意域名
	fi, err := os.Open("./Asset/BlackDomain/blacklist.txt")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		blackDomain = append(blackDomain, string(a))
	}

	laddr := ":53"
	err = DnsServer(laddr)
	if err != nil {
		fmt.Println("DNS 服务启动失败", err)
	}
	fmt.Println("DNS 代理启动成功")

}

// 判断目标域名是否在恶意域名中
func MatchDomain(domain string) bool {
	// 读取文件中的恶意域名

	var flag = false
	for _, evilDomainItem := range blackDomain {
		if evilDomainItem == domain {
			return true
		}
	}
	return flag
}

func RequestDNS(src string, dnsService string) (dst []string, err error) {

	c := dns.Client{Timeout: 5 * time.Second}

	var lasetErr error
	// retry 3 times
	for i := 0; i < 3; i++ {
		m := dns.Msg{}
		// 最终都会指向一个ip 也就是typeA, 这样就可以返回所有层的cname.
		m.SetQuestion(src+".", dns.TypeA)
		r, _, err := c.Exchange(&m, dnsService+":53")
		log.Println("原始请求记录：", r)

		if err != nil {
			time.Sleep(1 * time.Second * time.Duration(i+1))
			continue
		}
		dst = []string{}
		for _, ans := range r.Answer {
			record, isType := ans.(*dns.CNAME)
			if isType {
				dst = append(dst, record.Target)
			}
		}
		lasetErr = nil
		break
	}
	err = lasetErr

	return
}
