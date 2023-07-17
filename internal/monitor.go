package aikuaimonitor

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gosuri/uilive"
	"github.com/jakeslee/ikuai"
	"github.com/jakeslee/ikuai/action"
)

type IkuaiMonitor struct {
	Users    []*User
	Clients  []*ikuai.IKuai
	Sessions []string
	mu       sync.Mutex
}

type ShowMonitorInterfaceResultWithName struct {
	ClientName      string
	InterfaceResult *action.ShowMonitorInterfaceResult
}
type ShowMonitorResultWithName struct {
	ClientName string
	Result     *action.ShowMonitorResult
}

func (m *IkuaiMonitor) InitAllClinet() {
	if len(m.Users) != 0 {
		for _, user := range m.Users {
			i := ikuai.NewIKuai(user.Url, user.Name, user.Password)
			m.Clients = append(m.Clients, i)
			m.Sessions = append(m.Sessions, "")
		}
	}
}

func (m *IkuaiMonitor) LoginAllClinet() {
	if len(m.Clients) > 0 {
		for index, c := range m.Clients {
			session, err := m.LoginOneClient(c)
			if err != nil {
				log.Println(err)
			} else {
				m.Sessions[index] = session
			}
		}
	}
}

func (m *IkuaiMonitor) LoginOneClient(i *ikuai.IKuai) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return i.Login()
}

func (m *IkuaiMonitor) GetAllMonitorLan() map[string]*action.ShowMonitorResult {
	var wg sync.WaitGroup
	var result = make(map[string]*action.ShowMonitorResult)
	if len(m.Clients) != 0 {
		resultChain := make(chan ShowMonitorResultWithName, len(m.Clients))
		for index, client := range m.Clients {
			wg.Add(1)
			go func(c *ikuai.IKuai, i int) {
				defer wg.Done()
				dobuleCheckLogin := func() error {
					_, err := m.LoginOneClient(c)
					if err != nil {
						log.Println(err)
						resultChain <- ShowMonitorResultWithName{ClientName: m.Users[i].ClientName}
						return err
					}
					return nil
				}
				if len(m.Sessions) <= i || len(m.Sessions[i]) == 0 {
					err := dobuleCheckLogin()
					if err != nil {
						return
					}
				}
				lan, err := c.ShowMonitorLan()
				if err != nil {
					log.Println(err)
				}
				if err != nil || (lan != nil && lan.Result.Result == NOT_LOGIN) {
					err := dobuleCheckLogin()
					if err != nil {
						return
					}
					lan, err = c.ShowMonitorLan()
					if err != nil {
						log.Println(err)
					}
				}
				if err != nil {
					resultChain <- ShowMonitorResultWithName{ClientName: m.Users[i].ClientName}
				} else {
					resultChain <- ShowMonitorResultWithName{ClientName: m.Users[i].ClientName, Result: lan}
				}
			}(client, index)
		}
		wg.Wait()
		close(resultChain)
		for res := range resultChain {
			result[res.ClientName] = res.Result
		}
	}
	return result
}

func (m *IkuaiMonitor) GetMonitorInterface() map[string]*action.ShowMonitorInterfaceResult {
	var wg sync.WaitGroup
	var result = make(map[string]*action.ShowMonitorInterfaceResult)
	if len(m.Clients) != 0 {
		resultChain := make(chan ShowMonitorInterfaceResultWithName, len(m.Clients))
		for index, client := range m.Clients {
			wg.Add(1)
			go func(c *ikuai.IKuai, i int) {
				defer wg.Done()
				dobuleCheckLogin := func() error {
					_, err := m.LoginOneClient(c)
					if err != nil {
						log.Println(err)
						resultChain <- ShowMonitorInterfaceResultWithName{ClientName: m.Users[i].ClientName}
						return err
					}
					return nil
				}
				if len(m.Sessions) <= i || len(m.Sessions[i]) == 0 {
					fmt.Println("session 为空，取调用登录")
					err := dobuleCheckLogin()
					if err != nil {
						return
					}
				}
				res, err := c.ShowMonitorInterface()
				if err != nil {
					log.Println(err)
				}
				if err != nil || (res != nil && res.Result.Result == NOT_LOGIN) {
					fmt.Println("登录校验未通过，或者出错，尝试重新调用接口登录")
					err := dobuleCheckLogin()
					if err != nil {
						return
					}
					res, err = c.ShowMonitorInterface()
					if err != nil {
						log.Println(err)
					}
				}
				if err != nil || res == nil {
					log.Println(err)
					resultChain <- ShowMonitorInterfaceResultWithName{ClientName: m.Users[i].ClientName}
				} else {
					resultChain <- ShowMonitorInterfaceResultWithName{ClientName: m.Users[i].ClientName, InterfaceResult: res}
				}
			}(client, index)
		}
		wg.Wait()
		close(resultChain)
		for res := range resultChain {
			result[res.ClientName] = res.InterfaceResult
		}
	}
	return result
}

func (m *IkuaiMonitor) TransformMonitorInterface() map[string]map[string][][]string {
	interfaceResult := m.GetMonitorInterface()
	var finalResult = map[string]map[string][][]string{}
	for key, data := range interfaceResult {
		if data == nil {
			continue
		}
		d := data.Data
		ifaceChecks := d.IfaceCheck
		IfaceStreams := d.IfaceStream
		strArrIfaceChecks := [][]string{}
		strArrIfaceStreams := [][]string{}
		strArrIfaceChecks = append(strArrIfaceChecks, []string{"Name", "Ip", "Internet", "Result", "Comment"})
		for _, ifaceCheck := range ifaceChecks {
			strArrIfaceChecks = append(strArrIfaceChecks, []string{
				ifaceCheck.Interface,
				ifaceCheck.IPAddr,
				ifaceCheck.Internet,
				ifaceCheck.Result,
				ifaceCheck.Comment})
		}
		strArrIfaceStreams = append(strArrIfaceStreams, []string{"Name", "Ip", "Upload", "Download", "TotoalUpload", "TotalDownload", "Comment"})
		for _, ifaceStream := range IfaceStreams {
			strArrIfaceStreams = append(strArrIfaceStreams, []string{
				ifaceStream.Interface,
				ifaceStream.IPAddr,
				FormatSize(float64(ifaceStream.Upload)),
				FormatSize(float64(ifaceStream.Download)),
				FormatSize(float64(ifaceStream.TotalUp)),
				FormatSize(float64(ifaceStream.TotalDown)),
				ifaceStream.Comment})
		}
		finalResult[key] = map[string][][]string{
			"ifaceChecks":  strArrIfaceChecks,
			"ifaceStreams": strArrIfaceStreams,
		}
	}
	return finalResult
}

var Monitor IkuaiMonitor

func InitMonitor() {
	user, err := InitIni()
	if err != nil {
		panic("获取用户信息错误，请检测conf.ini文件")
	}
	Monitor = IkuaiMonitor{
		Users: user,
	}
	Monitor.InitAllClinet()
	Monitor.LoginAllClinet()
}

func ScheduleMonitorInterface(loop int) {
	writer := uilive.New()
	writer.Start()
	singalChain := make(chan os.Signal, 1)
	signal.Notify(singalChain, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-singalChain
		writer.Stop()
		os.Exit(0)
	}()
	ticker := time.Tick(time.Duration(loop) * time.Second)
	for range ticker {
		PrintData(Monitor.TransformMonitorInterface(), writer)
	}
	// for {
	// 	PrintData(Monitor.TransformMonitorInterface(), writer)
	// 	time.Sleep(time.Second * time.Duration(loop))
	// }

}
