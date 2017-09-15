package nut

import (
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/user"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/astaxie/beego"
)

// GetAdminSiteStatus site status
// @router /admin/status [get]
func (p *Plugin) GetAdminSiteStatus() {
	p.LayoutDashboard()
	p.MustAdmin()

	var err error
	p.Data["os"], err = p._osStatus()
	if err == nil {
		p.Data["network"], err = p._networkStatus()
	}
	if err == nil {
		p.Data["cache"], err = p._cacheStatus()
	}
	if err == nil {
		p.Data["database"], err = p._dbStatus()
	}
	if err == nil {
		p.Data["jobs"], err = p._jobStatus()
	}

	p.Flash(nil, err)

	p.Data[TITLE] = Tr(p.Locale(), "nut.admin.site.status.title")
	p.TplName = "nut/admin/site/status.html"
}

func (p *Plugin) _jobStatus() (H, error) {
	val := H{}
	for k, v := range JOBBER().handlers {
		val[k] = runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name()
	}
	return val, nil
}
func (p *Plugin) _osStatus() (H, error) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	hn, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	hu, err := user.Current()
	if err != nil {
		return nil, err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	var ifo syscall.Sysinfo_t
	if err := syscall.Sysinfo(&ifo); err != nil {
		return nil, err
	}
	return H{
		"app author":           fmt.Sprintf("%s <%s>", AuthorName, AuthorEmail),
		"app licence":          Copyright,
		"app version":          fmt.Sprintf("%s(%s) - %s", Version, BuildTime, beego.BConfig.RunMode),
		"app root":             pwd,
		"who-am-i":             fmt.Sprintf("%s@%s", hu.Username, hn),
		"go version":           runtime.Version(),
		"go root":              runtime.GOROOT(),
		"go runtime":           runtime.NumGoroutine(),
		"go last gc":           time.Unix(0, int64(mem.LastGC)).Format(time.ANSIC),
		"os cpu":               runtime.NumCPU(),
		"os ram(free/total)":   fmt.Sprintf("%dM/%dM", ifo.Freeram/1024/1024, ifo.Totalram/1024/1024),
		"os swap(free/total)":  fmt.Sprintf("%dM/%dM", ifo.Freeswap/1024/1024, ifo.Totalswap/1024/1024),
		"go memory(alloc/sys)": fmt.Sprintf("%dM/%dM", mem.Alloc/1024/1024, mem.Sys/1024/1024),
		"os time":              time.Now().Format(time.ANSIC),
		"os arch":              fmt.Sprintf("%s(%s)", runtime.GOOS, runtime.GOARCH),
		"os uptime":            (time.Duration(ifo.Uptime) * time.Second).String(),
		"os loads":             ifo.Loads,
		"os procs":             ifo.Procs,
	}, nil
}
func (p *Plugin) _networkStatus() (H, error) {
	sts := H{}
	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, v := range ifs {
		ips := []string{v.HardwareAddr.String()}
		adrs, err := v.Addrs()
		if err != nil {
			return nil, err
		}
		for _, adr := range adrs {
			ips = append(ips, adr.String())
		}
		sts[v.Name] = ips
	}
	return sts, nil
}

func (p *Plugin) _dbStatus() (H, error) {
	val := H{
		"drivers": strings.Join(sql.Drivers(), ", "),
	}
	switch beego.AppConfig.String("databasedriver") {
	case "postgres":
		// err := orm.NewOrm().Raw("select version()").V
		// var version string
		// row.Scan(&version)
		// val["version"] = version
		// // http://blog.javachen.com/2014/04/07/some-metrics-in-postgresql.html
		// row = p.Db.Raw("select pg_size_pretty(pg_database_size('postgres'))").Row()
		// var size string
		// row.Scan(&size)
		// val["size"] = size
		// if rows, err := p.Db.
		// 	Raw("select pid,current_timestamp - least(query_start,xact_start) AS runtime,substr(query,1,25) AS current_query from pg_stat_activity where not pid=pg_backend_pid()").
		// 	Rows(); err == nil {
		// 	defer rows.Close()
		// 	for rows.Next() {
		// 		var pid int
		// 		var ts time.Time
		// 		var qry string
		// 		row.Scan(&pid, &ts, &qry)
		// 		val[fmt.Sprintf("pid-%d", pid)] = fmt.Sprintf("%s (%v)", ts.Format("15:04:05.999999"), qry)
		// 	}
		// } else {
		// 	return nil, err
		// }
		val["url"] = beego.AppConfig.String("databasesource")
	}
	return val, nil
}

func (p *Plugin) _cacheStatus() ([]string, error) {
	// str, err := Cache().
	// if err != nil {
	// 	return nil, err
	// }
	// return strings.Split(str, "\n"), nil
	return []string{}, nil
}
