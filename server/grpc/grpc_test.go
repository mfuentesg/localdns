package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/mfuentesg/localdns/pb"
	"github.com/mfuentesg/localdns/storage"
	"github.com/stretchr/testify/assert"
)

type fakeStorage struct {
	storage.Storage
}

func (fs *fakeStorage) Get(identifier string) (*storage.Record, error) {
	if identifier == "unknown" {
		return nil, errors.New("unhandled error")
	}

	return &storage.Record{
		Type:   "A",
		Domain: identifier,
		IP:     "192.168.1.10",
	}, nil
}

func (fs *fakeStorage) Put(r storage.Record) error {
	if r.Type == "AAAA" {
		return errors.New("unhandled error")
	}
	return nil
}

func (fs *fakeStorage) Delete(identifier string) error {
	if identifier == "www.delete.com." {
		return errors.New("unhandled error")
	}
	return nil
}

func (fs *fakeStorage) List() ([]*storage.Record, error) {
	return []*storage.Record{
		{Type: "A", Domain: "www.fake.com", IP: "10.168.1.1"},
	}, nil
}

func TestNew(t *testing.T) {
	s := New()
	assert.IsType(t, new(Server), s)
	assert.Equal(t, ":8080", s.Addr)
	assert.Nil(t, s.st)
}

func TestWithAddr(t *testing.T) {
	s := New(WithAddr(":9090"))
	assert.Equal(t, ":9090", s.Addr)
}

func TestWithStorage(t *testing.T) {
	s := New(WithStorage(new(fakeStorage)))
	assert.NotNil(t, s.st)
}

func TestServer_DeleteRecord(t *testing.T) {
	s := New(WithStorage(new(fakeStorage)))

	t.Run("storage failure", func(tt *testing.T) {
		_, err := s.DeleteRecord(context.Background(), &pb.Record{
			Domain: "www.delete.com.",
		})

		assert.Error(tt, err)
	})

	t.Run("deleted correctly", func(tt *testing.T) {
		_, err := s.DeleteRecord(context.Background(), &pb.Record{
			Domain: "www.fake.com.",
		})

		assert.Nil(tt, err)
	})
}

func TestServer_GetRecord(t *testing.T) {
	s := New(WithStorage(new(fakeStorage)))

	t.Run("error on storage", func(tt *testing.T) {
		record, err := s.GetRecord(context.Background(), &pb.Record{
			Domain: "unknown",
		})

		assert.NotNil(tt, err)
		assert.Error(tt, err)
		assert.Nil(tt, record)
	})

	t.Run("valid storage response", func(tt *testing.T) {
		record, err := s.GetRecord(context.Background(), &pb.Record{
			Domain: "www.my-domain.com.",
		})

		assert.Nil(tt, err)
		assert.NoError(tt, err)
		assert.NotNil(tt, record)
		assert.Equal(tt, "192.168.1.10", record.Ip)
		assert.Equal(tt, "www.my-domain.com.", record.Domain)
	})
}

func TestServer_PutRecord(t *testing.T) {
	s := New(WithStorage(new(fakeStorage)))

	t.Run("storage failure", func(tt *testing.T) {
		record, err := s.PutRecord(context.Background(), &pb.Record{
			Type:   "AAAA",
			Domain: "www.my-domain.com",
		})

		assert.Error(tt, err)
		assert.Nil(tt, record)
	})

	t.Run("append `.` as part of the domain", func(tt *testing.T) {
		record, err := s.PutRecord(context.Background(), &pb.Record{
			Type:   "A",
			Domain: "www.my-domain.com",
		})

		assert.Nil(tt, err)
		assert.NotNil(tt, record)
		assert.Equal(tt, "www.my-domain.com.", record.Domain)
	})
}

func TestServer_ListRecords(t *testing.T) {
	s := New(WithStorage(new(fakeStorage)))

	list, err := s.ListRecords(context.Background(), nil)

	assert.Nil(t, err)
	assert.Len(t, list.Records, 1)
}

func TestServer_ListenAndServe(t *testing.T) {
}
