package service

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"virtpanel/internal/model"
)

// getDefaultGateway reads the current default gateway from ip route
func getDefaultGateway() string {
	out, err := exec.Command("ip", "route", "show", "default").Output()
	if err != nil {
		return ""
	}
	// "default via 192.168.8.1 dev eth0 ..."
	fields := strings.Fields(string(out))
	for i, f := range fields {
		if f == "via" && i+1 < len(fields) {
			return fields[i+1]
		}
	}
	return ""
}

const bridgePrefix = "vp-"

func (s *LibvirtService) ListBridges() ([]model.Bridge, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	result := make([]model.Bridge, 0)
	for _, iface := range ifaces {
		if !strings.HasPrefix(iface.Name, bridgePrefix) {
			continue
		}
		// Verify it's actually a bridge
		if _, err := os.Stat("/sys/class/net/" + iface.Name + "/bridge"); err != nil {
			continue
		}
		br := model.Bridge{
			Name: iface.Name,
			Up:   iface.Flags&net.FlagUp != 0,
		}
		// Get slaves
		entries, _ := os.ReadDir("/sys/class/net/" + iface.Name + "/brif")
		for _, e := range entries {
			br.Slaves = append(br.Slaves, e.Name())
		}
		// Get IP
		addrs, _ := iface.Addrs()
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				br.IP = ipnet.IP.String()
				break
			}
		}
		result = append(result, br)
	}
	return result, nil
}

func (s *LibvirtService) CreateBridge(req model.CreateBridgeRequest) error {
	name := bridgePrefix + req.Name
	if !safeNameRe.MatchString(req.Name) {
		return fmt.Errorf("invalid bridge name: %s", req.Name)
	}
	if _, err := net.InterfaceByName(name); err == nil {
		return fmt.Errorf("网桥 %s 已存在", name)
	}

	var script strings.Builder
	script.WriteString(fmt.Sprintf("ip link add %s type bridge && ip link set %s up", name, name))

	if req.SlaveNIC != "" {
		if !safeNameRe.MatchString(req.SlaveNIC) {
			return fmt.Errorf("invalid NIC name: %s", req.SlaveNIC)
		}
		iface, err := net.InterfaceByName(req.SlaveNIC)
		if err != nil {
			return fmt.Errorf("网卡 %s 不存在", req.SlaveNIC)
		}
		script.WriteString(fmt.Sprintf(" && ip link set %s master %s", req.SlaveNIC, name))
		gw := getDefaultGateway()
		addrs, _ := iface.Addrs()
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				cidr := ipnet.String()
				script.WriteString(fmt.Sprintf(" && ip addr del %s dev %s; ip addr add %s dev %s", cidr, req.SlaveNIC, cidr, name))
				if gw != "" {
					script.WriteString(fmt.Sprintf(" && ip route add default via %s dev %s 2>/dev/null", gw, name))
				}
			}
		}
	}
	return exec.Command("bash", "-c", script.String()).Run()
}

func (s *LibvirtService) DeleteBridge(name string) error {
	full := bridgePrefix + name
	if _, err := net.InterfaceByName(full); err != nil {
		return fmt.Errorf("网桥 %s 不存在", full)
	}
	// Move IPs back to slave before deleting
	entries, _ := os.ReadDir("/sys/class/net/" + full + "/brif")
	if len(entries) > 0 {
		slave := entries[0].Name()
		iface, _ := net.InterfaceByName(full)
		if iface != nil {
			addrs, _ := iface.Addrs()
			gw := getDefaultGateway()
			var script strings.Builder
			script.WriteString(fmt.Sprintf("ip link set %s nomaster", slave))
			for _, a := range addrs {
				if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.To4() != nil {
					cidr := ipnet.String()
					script.WriteString(fmt.Sprintf(" && ip addr add %s dev %s", cidr, slave))
					if gw != "" {
						script.WriteString(fmt.Sprintf(" && ip route add default via %s dev %s 2>/dev/null", gw, slave))
					}
				}
			}
			script.WriteString(fmt.Sprintf(" && ip link del %s", full))
			return exec.Command("bash", "-c", script.String()).Run()
		}
	}
	return exec.Command("ip", "link", "del", full).Run()
}
