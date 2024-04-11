package mempool

import (
	"encoding/binary"
	"sync/atomic"
	"testing"

	"github.com/cometbft/cometbft/abci/example/kvstore"
	"github.com/cometbft/cometbft/proxy"
)

func BenchmarkReap(b *testing.B) {
	app := kvstore.NewInMemoryApplication()
	cc := proxy.NewLocalClientCreator(app)
	mp, cleanup := newMempoolWithApp(cc)
	defer cleanup()

	mp.config.Size = 100000

	size := 10000
	for i := 0; i < size; i++ {
		tx := make([]byte, 8)
		binary.BigEndian.PutUint64(tx, uint64(i))
		if _, err := mp.CheckTxSync(tx, TxInfo{}); err != nil {
			b.Fatal(err)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mp.ReapMaxBytesMaxGas(100000000, 10000000)
	}
}

func BenchmarkCheckTxSync(b *testing.B) {
	app := kvstore.NewInMemoryApplication()
	cc := proxy.NewLocalClientCreator(app)
	mp, cleanup := newMempoolWithApp(cc)
	defer cleanup()

	mp.config.Size = 1000000

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		tx := make([]byte, 8)
		binary.BigEndian.PutUint64(tx, uint64(i))
		b.StartTimer()

		if _, err := mp.CheckTxSync(tx, TxInfo{}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParallelCheckTxSync(b *testing.B) {
	app := kvstore.NewInMemoryApplication()
	cc := proxy.NewLocalClientCreator(app)
	mp, cleanup := newMempoolWithApp(cc)
	defer cleanup()

	mp.config.Size = 100000000

	var txcnt uint64
	next := func() uint64 {
		return atomic.AddUint64(&txcnt, 1) - 1
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tx := make([]byte, 8)
			binary.BigEndian.PutUint64(tx, next())
			if _, err := mp.CheckTxSync(tx, TxInfo{}); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkCheckDuplicateTxSync(b *testing.B) {
	app := kvstore.NewInMemoryApplication()
	cc := proxy.NewLocalClientCreator(app)
	mp, cleanup := newMempoolWithApp(cc)
	defer cleanup()

	mp.config.Size = 1000000

	for i := 0; i < b.N; i++ {
		tx := make([]byte, 8)
		binary.BigEndian.PutUint64(tx, uint64(i))
		if _, err := mp.CheckTxSync(tx, TxInfo{}); err != nil {
			b.Fatal(err)
		}

		if _, err := mp.CheckTxSync(tx, TxInfo{}); err == nil {
			b.Fatal("tx should be duplicate")
		}
	}
}

func BenchmarkCheckTxAsync(b *testing.B) {
	app := kvstore.NewInMemoryApplication()
	cc := proxy.NewLocalClientCreator(app)
	mp, cleanup := newMempoolWithApp(cc)
	defer cleanup()

	mp.config.Size = 1000000

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		tx := make([]byte, 8)
		binary.BigEndian.PutUint64(tx, uint64(i))
		b.StartTimer()
		mp.CheckTxAsync(tx, TxInfo{}, nil, nil)
	}
}

func BenchmarkParallelCheckTxAsync(b *testing.B) {
	app := kvstore.NewInMemoryApplication()
	cc := proxy.NewLocalClientCreator(app)
	mp, cleanup := newMempoolWithApp(cc)
	defer cleanup()

	mp.config.Size = 100000000

	var txcnt uint64
	next := func() uint64 {
		return atomic.AddUint64(&txcnt, 1) - 1
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tx := make([]byte, 8)
			binary.BigEndian.PutUint64(tx, next())
			mp.CheckTxAsync(tx, TxInfo{}, nil, nil)
		}
	})
}

func BenchmarkCheckDuplicateTxAsync(b *testing.B) {
	app := kvstore.NewInMemoryApplication()
	cc := proxy.NewLocalClientCreator(app)
	mp, cleanup := newMempoolWithApp(cc)
	defer cleanup()

	mp.config.Size = 1000000

	for i := 0; i < b.N; i++ {
		tx := make([]byte, 8)
		binary.BigEndian.PutUint64(tx, uint64(i))
		mp.CheckTxAsync(tx, TxInfo{}, nil, nil)
		mp.CheckTxAsync(tx, TxInfo{}, nil, nil)
	}
}
