package grpc

import (
	"context"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/mfuentesg/localdns/pb"
	"github.com/mfuentesg/localdns/storage"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type fakeStorage struct {
	storage.Storage
}

const (
	invalidRecordID = "7af2bc99-600d-468d-b1a8-da2a72ba8538"
	validRecordID   = "9964f6cb-29f4-41bd-9770-4c175098206c"
)

func (fs *fakeStorage) Get(identifier string) (*storage.Record, error) {
	if identifier == invalidRecordID {
		return nil, errors.New("unhandled error")
	}

	return &storage.Record{
		ID:     validRecordID,
		Type:   "A",
		Domain: "www.my-domain.com.",
		IPv4:   "192.168.1.10",
	}, nil
}

func (fs *fakeStorage) Put(r storage.Record) (string, error) {
	if r.Type == "AAAA" {
		return "", errors.New("unhandled error")
	}
	return validRecordID, nil
}

func (fs *fakeStorage) Delete(identifier string) error {
	if identifier == invalidRecordID {
		return errors.New("unhandled error")
	}
	return nil
}

func (fs *fakeStorage) List() ([]*storage.Record, error) {
	return []*storage.Record{
		{ID: validRecordID, Type: "A", Domain: "www.fake.com", IPv4: "10.168.1.1"},
	}, nil
}

func TestNew(t *testing.T) {
	s := New(nil)
	assert.IsType(t, new(Server), s)
	assert.Equal(t, ":8080", s.Addr)
	assert.Nil(t, s.st)
}

func TestWithAddr(t *testing.T) {
	s := New(nil, WithAddr(":9090"))
	assert.Equal(t, ":9090", s.Addr)
}

func TestServer_DeleteRecord(t *testing.T) {
	s := New(new(fakeStorage))
	log.SetOutput(ioutil.Discard)

	t.Run("storage failure", func(tt *testing.T) {
		_, err := s.DeleteRecord(context.Background(), &pb.Record{
			Id: invalidRecordID,
		})

		assert.Error(tt, err)
	})

	t.Run("deleted correctly", func(tt *testing.T) {
		_, err := s.DeleteRecord(context.Background(), &pb.Record{
			Id: validRecordID,
		})

		assert.Nil(tt, err)
	})
}

func TestServer_GetRecord(t *testing.T) {
	s := New(new(fakeStorage))

	t.Run("error on storage", func(tt *testing.T) {
		record, err := s.GetRecord(context.Background(), &pb.Record{
			Id: invalidRecordID,
		})

		assert.NotNil(tt, err)
		assert.Error(tt, err)
		assert.Nil(tt, record)
	})

	t.Run("valid storage response", func(tt *testing.T) {
		record, err := s.GetRecord(context.Background(), &pb.Record{
			Id: validRecordID,
		})

		assert.Nil(tt, err)
		assert.NoError(tt, err)
		assert.NotNil(tt, record)
		assert.Equal(tt, "192.168.1.10", record.Ipv4)
		assert.Equal(tt, "www.my-domain.com.", record.Domain)
	})
}

func TestServer_PutRecord(t *testing.T) {
	s := New(new(fakeStorage))

	t.Run("storage failure", func(tt *testing.T) {
		record, err := s.PutRecord(context.Background(), &pb.Record{
			Type:   "AAAA",
			Domain: "www.my-domain.com",
		})

		assert.Error(tt, err)
		assert.Nil(tt, record)
	})

	t.Run("append `.` as part of the domain", func(tt *testing.T) {
		log.SetOutput(ioutil.Discard)

		record, err := s.PutRecord(context.Background(), &pb.Record{
			Type:   "A",
			Domain: "www.my-domain.com",
		})

		assert.Nil(tt, err)
		assert.NotNil(tt, record)
		assert.Equal(tt, "www.my-domain.com.", record.Domain)
		assert.Equal(tt, validRecordID, record.Id)
	})
}

func TestServer_ListRecords(t *testing.T) {
	s := New(new(fakeStorage))

	list, err := s.ListRecords(context.Background(), nil)

	assert.Nil(t, err)
	assert.Len(t, list.Records, 1)
}

func TestServer_ListenAndServe(t *testing.T) {
}
