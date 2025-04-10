package worker

import (
	"blog/internal/domain/friendLink/infrastructure/db"
	"blog/internal/domain/friendLink/infrastructure/notifier"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

type Worker interface {
	Start()
	Stop()
}

type worker struct {
	db       db.DB
	notifier notifier.Notifier
	ticker   *time.Ticker
	stopChan chan struct{}
}

func NewWorker(db db.DB, notifier notifier.Notifier) Worker {
	return &worker{
		db:       db,
		notifier: notifier,
		stopChan: make(chan struct{}),
	}
}

const workerInternal = 24 * time.Hour

func (w *worker) Start() {
	w.ticker = time.NewTicker(workerInternal) // 每5分钟同步一次

	go func() {
		// 启动后立即执行一次
		if err := w.notifyPendingLinks(); err != nil {
			sugar := zap.L().Sugar()
			sugar.Infof("友链审核提醒通知发送失败%+v", err)
		}

		for {
			select {
			case <-w.ticker.C:
				if err := w.notifyPendingLinks(); err != nil {
					zap.L().Error("发送友链审核提醒失败", zap.Error(err))
				}
			case <-w.stopChan:
				w.ticker.Stop()
				return
			}
		}
	}()

	zap.L().Info("友链审核提醒任务已启动")
}

func (w *worker) Stop() {
	close(w.stopChan)
	zap.L().Info("友链审核提醒任务已停止")
}

func (w *worker) notifyPendingLinks() error {
	// 获取所以pending状态的友链
	links, err := w.db.GetPendingLinks()
	if err != nil {
		return err
	}

	if len(links) == 0 {
		return nil
	}

	if err = w.notifier.SendPendingNotification(links); err != nil {
		return errors.WithStack(err)
	}

	// 将普通logger转为Sugar logger
	sugar := zap.L().Sugar()
	sugar.Infof("成功读取 %d 个等待审核的友链，并发送提醒通知", len(links))
	return nil
}
