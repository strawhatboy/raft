/*
 * Filename: /mnt/c/Users/PureAdmin/tj/courses/raft/store.go
 * Path: /mnt/c/Users/PureAdmin/tj/courses/raft
 * Created Date: Thursday, January 1st 1970, 8:00:00 am
 * Author: strawhatboy
 * 
 * Copyright (c) 2020 Your Company
 */

package core

type Store struct {

}

func (s *Store) Get(key string) (string, error) {
	return "", nil
}

func (s *Store) Put(key string, value string) error {
	return nil
}

func (s *Store) Join(addr string) error {
	return nil
}
