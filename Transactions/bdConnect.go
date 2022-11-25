package main

import (
	"context"
	"errors"
	"fmt"
	"time"
	"github.com/jackc/pgx/v4/pgxpool"
)

type dbConn struct {
	conn *pgxpool.Pool
}

func (d *dbConn) createConnection(name, pass, url, port, dbName string) (err error) {
	after := time.After(10 * time.Second)
	for {
		select {
		case <-after:
			return err
		default:
			d.conn, err = pgxpool.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s", name, pass, url, port, dbName))
			if err == nil {
			return err
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func (d *dbConn) Close() {
	d.conn.Close()
}

func (d *dbConn) makeTransaction(tr Transaction) (ResultTransaction, error) {
	var (
		rtr           ResultTransaction
		balanceSender int64
	)
	err := d.conn.QueryRow(context.Background(), "SELECT balance FROM users WHERE user_id = $1", tr.SenderId).Scan(&balanceSender)
	balanceSender -= tr.Count
	if balanceSender < 0 {
		return rtr, errors.New("sender have not enough money")
	}
	if err != nil {
		return rtr, err
	}
	err = newTransaction(d, tr)
	if err != nil {
		return rtr, err
	}
	rtr.Remains = balanceSender
	rtr.UserId = tr.SenderId
	return rtr, nil
}

func newTransaction(d *dbConn, tr Transaction) (err error) {
	tx, err := d.conn.Begin(context.Background())
	if err != nil {
		return errors.New("Cant create prepeared transaction")
	}
	defer tx.Rollback(context.Background())
	_, err = tx.Exec(context.Background(), "UPDATE users SET balance=balance-$1 WHERE user_id=$2", tr.Count, tr.SenderId)
	if err != nil {
		return errors.New("Cant update client")
	}
	_, err = tx.Exec(context.Background(), "UPDATE users SET balance=balance+$1 WHERE user_id=$2", tr.Count, tr.ReceiverId)
	if err != nil {
		return errors.New("Cant update client")
	}
	_, err = tx.Exec(context.Background(), "INSERT INTO transactions(sender, receiver, count) VALUES($1, $2, $3)", tr.SenderId, tr.ReceiverId, tr.Count)
	if err != nil {
		return errors.New("Cant execute transaction")
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return
}
