package cpu

import (
	"path"
	"strconv"
	"strings"
)

const cgroupRootDir = "/sys/fs/cgroup"

type cgroup struct {
	cgroupSet map[string]string
}

func (c *cgroup) CPUCFSQuotaUs() (int64, error) {
	data, err := readFile(path.Join(c.cgroupSet["cpu"], "cpu.cfs_quota_us"))
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(data, 10, 64)
}

func (c *cgroup) CPUAcctUsage() (uint64, error) {
	data, err := readFile(path.Join(c.cgroupSet["cpuacct"], "cpuacct.usage"))
	if err != nil {
		return 0, err
	}
	return parseUint(data)
}

func (c *cgroup) CPUAcctUsagePreCPU() ([]uint64, error) {
	data, err := readFile(path.Join(c.cgroupSet["cpuacct"], "cpuacct.usage_percpu"))
	if err != nil {
		return nil, err
	}
	var usage []uint64
	for _, v := range strings.Fields(string(data)) {
		var u uint64
		if u, err = parseUint(v); err != nil {
			return nil, err
		}
		if u != 0 {
			usage = append(usage, u)
		}
	}
	return usage, nil
}
