package lib

import (
    "errors"
    "fmt"
    "log"
    "net"
    "net/http"
    // "net/url" // Use only in parsing request Header.
)

// https://github.com/gin-gonic/gin/issues/604

// GetClientIPHelper gets the client IP using a mixture of techniques.
// This is how it is with golang at the moment.
func GetClientIPHelper(req *http.Request) (ipResult string, errResult error) {

    // Try lots of ways :) Order is important.

    //  IA: Turn off this method in favor of trying remote address.
    //  
    //  Try Request Header ("Origin")
    // url, err := url.Parse(req.Header.Get("Origin"))
    // if err == nil {
    //     host := url.Host
    //     ip, _, err := net.SplitHostPort(host)
    //     if err == nil {
    //         log.Printf("debug: Found IP using Header (Origin) sniffing. ip: %v", ip)
    //         return ip, nil
    //     }
    // }

    // Try by Request
    ip, err := getClientIPByRequestRemoteAddr(req)
    if err == nil {
        log.Printf("debug: Found IP using Request sniffing. ip: %v", ip)
        return ip, nil
    }

    // Try Request Headers (X-Forwarder). Client could be behind a Proxy
    ip, err = getClientIPByHeaders(req)
    if err == nil {
        log.Printf("debug: Found IP using Request Headers sniffing. ip: %v", ip)
        return ip, nil
    }

    err = errors.New("error: Could not find clients IP address")
    return "", err
}

// getClientIPByRequest tries to get directly from the Request.
// https://blog.golang.org/context/userip/userip.go
func getClientIPByRequestRemoteAddr(req *http.Request) (ip string, err error) {

    // Try via request
    ip, port, err := net.SplitHostPort(req.RemoteAddr)
    if err != nil {
        log.Printf("debug: Getting req.RemoteAddr %v", err)
        return "", err
    } else {
        log.Printf("debug: With req.RemoteAddr found IP:%v; Port: %v", ip, port)
    }

    userIP := net.ParseIP(ip)
    if userIP == nil {
        message := fmt.Sprintf("debug: Parsing IP from Request.RemoteAddr got nothing.")
        log.Printf(message)
        return "", fmt.Errorf(message)

    }
    log.Printf("debug: Found IP: %v", userIP)
    return userIP.String(), nil

}

// getClientIPByHeaders tries to get directly from the Request Headers.
// This is only way when the client is behind a Proxy.
func getClientIPByHeaders(req *http.Request) (ip string, err error) {

    // Client could be behid a Proxy, so Try Request Headers (X-Forwarder)
    ipSlice := []string{}

    ipSlice = append(ipSlice, req.Header.Get("X-Forwarded-For"))
    ipSlice = append(ipSlice, req.Header.Get("x-forwarded-for"))
    ipSlice = append(ipSlice, req.Header.Get("X-FORWARDED-FOR"))

    for _, v := range ipSlice {
        log.Printf("debug: client request header check gives ip: %v", v)
        if v != "" {
            return v, nil
        }
    }
    err = errors.New("error: Could not find clients IP address from the Request Headers")
    return "", err

}

// getMyInterfaceAddr gets this private network IP. Basically the Servers IP.
func getMyInterfaceAddr() (net.IP, error) {

    ifaces, err := net.Interfaces()
    if err != nil {
        return nil, err
    }
    addresses := []net.IP{}
    for _, iface := range ifaces {

        if iface.Flags&net.FlagUp == 0 {
            continue // interface down
        }
        if iface.Flags&net.FlagLoopback != 0 {
            continue // loopback interface
        }
        addrs, err := iface.Addrs()
        if err != nil {
            continue
        }

        for _, addr := range addrs {
            var ip net.IP
            switch v := addr.(type) {
            case *net.IPNet:
                ip = v.IP
            case *net.IPAddr:
                ip = v.IP
            }
            if ip == nil || ip.IsLoopback() {
                continue
            }
            ip = ip.To4()
            if ip == nil {
                continue // not an ipv4 address
            }
            addresses = append(addresses, ip)
        }
    }
    if len(addresses) == 0 {
        return nil, fmt.Errorf("no address Found, net.InterfaceAddrs: %v", addresses)
    }
    //only need first
    return addresses[0], nil
}