package worker

import (
	"blog/internal/domain/essay/repository/cache"
	"blog/internal/domain/essay/repository/db"
	"go.uber.org/zap"
	"time"
)

type Worker interface {
	Start()
	Stop()
}

type worker struct {
	db       db.DB
	cache    cache.Cache
	ticker   *time.Ticker
	stopChan chan struct{}
}

func NewWorker(db db.DB, cache cache.Cache) Worker {
	return &worker{
		db:       db,
		cache:    cache,
		stopChan: make(chan struct{}),
	}
}

const workerInternal = 1 * time.Minute

func (w *worker) Start() {
	w.ticker = time.NewTicker(workerInternal) // 每5分钟同步一次

	go func() {
		// 启动后立即执行一次同步
		if err := w.syncVisitedTimes(); err != nil {
			zap.L().Error("同步访问次数失败: %v", zap.Error(err))
		}

		for {
			select {
			case <-w.ticker.C:
				if err := w.syncVisitedTimes(); err != nil {
					zap.L().Error("定时同步访问次数失败: %v", zap.Error(err))
				}
			case <-w.stopChan:
				w.ticker.Stop()
				return
			}
		}
	}()
	zap.L().Info("访问次数同步任务已启动")
}

func (w *worker) Stop() {
	close(w.stopChan)
	zap.L().Info("访问次数同步任务已停止")
}

func (w *worker) syncVisitedTimes() error {
	// 从缓存获取所有访问次数
	vtMap, err := w.cache.GetAllVt()
	if err != nil {
		return err
	}

	if len(vtMap) == 0 {
		return nil
	}

	// 保存到数据库
	if err := w.db.SaveVTsByIDs(vtMap); err != nil {
		return err
	}

	// 将普通logger转为Sugar logger
	sugar := zap.L().Sugar()
	sugar.Infof("成功同步 %d 个文章的访问次数", len(vtMap))
	return nil
}
