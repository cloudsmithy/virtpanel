package service

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"virtpanel/internal/model"
)

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
		addrs, _ := iface.Addrs()
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.To4() != nil {
				cidr := ipnet.String()
				gw := fmt.Sprintf("%d.%d.%d.1", ipnet.IP.To4()[0], ipnet.IP.To4()[1], ipnet.IP.To4()[2])
				script.WriteString(fmt.Sprintf(" && ip addr del %s dev %s; ip addr add %s dev %s && ip route add default via %s dev %s 2>/dev/null", cidr, req.SlaveNIC, cidr, name, gw, name))
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
			var script strings.Builder
			script.WriteString(fmt.Sprintf("ip link set %s nomaster", slave))
			for _, a := range addrs {
				if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.To4() != nil {
					cidr := ipnet.String()
					gw := fmt.Sprintf("%d.%d.%d.1", ipnet.IP.To4()[0], ipnet.IP.To4()[1], ipnet.IP.To4()[2])
					script.WriteString(fmt.Sprintf(" && ip addr add %s dev %s && ip route add default via %s dev %s 2>/dev/null", cidr, slave, gw, slave))
				}
			}
			script.WriteString(fmt.Sprintf(" && ip link del %s", full))
			return exec.Command("bash", "-c", script.String()).Run()
		}
	}
	return exec.Command("ip", "link", "del", full).Run()
}
