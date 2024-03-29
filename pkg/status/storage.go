package status

import (
	"github.com/cossim/hipush/pkg/consts"
	"github.com/cossim/hipush/pkg/store"
)

type StateStorage struct {
	store store.Store
}

func NewStateStorage(store store.Store) *StateStorage {
	return &StateStorage{store: store}
}

func (s *StateStorage) Init() error {
	return s.store.Init()
}

func (s *StateStorage) Close() error {
	return s.store.Close()
}

// Reset Client storage.
func (s *StateStorage) Reset() {
	s.store.Set(consts.HiPushTotal, 0)
	s.store.Set(consts.HiPushSuccess, 0)
	s.store.Set(consts.HiPushFailed, 0)

	s.store.Set(consts.HTTPTotal, 0)
	s.store.Set(consts.HTTPSuccess, 0)
	s.store.Set(consts.HTTPFailed, 0)

	s.store.Set(consts.GRPCTotal, 0)
	s.store.Set(consts.GRPCSuccess, 0)
	s.store.Set(consts.GRPCFailed, 0)

	s.store.Set(consts.IosTotal, 0)
	s.store.Set(consts.IosSuccess, 0)
	s.store.Set(consts.IosFailed, 0)

	s.store.Set(consts.HuaweiTotal, 0)
	s.store.Set(consts.HuaweiSuccess, 0)
	s.store.Set(consts.HuaweiFailed, 0)

	s.store.Set(consts.AndroidTotal, 0)
	s.store.Set(consts.AndroidSuccess, 0)
	s.store.Set(consts.AndroidFailed, 0)

	s.store.Set(consts.VivoTotal, 0)
	s.store.Set(consts.VivoSuccess, 0)
	s.store.Set(consts.VivoFailed, 0)

	s.store.Set(consts.OppoTotal, 0)
	s.store.Set(consts.OppoSuccess, 0)
	s.store.Set(consts.OppoFailed, 0)

	s.store.Set(consts.XiaomiTotal, 0)
	s.store.Set(consts.XiaomiSuccess, 0)
	s.store.Set(consts.XiaomiFailed, 0)

	s.store.Set(consts.MeizuTotal, 0)
	s.store.Set(consts.MeizuSuccess, 0)
	s.store.Set(consts.MeizuFailed, 0)
}

func (s *StateStorage) AddTotalCount(count int64) {
	s.store.Add(consts.HiPushTotal, count)
}

func (s *StateStorage) AddHttpTotal(count int64) {
	s.store.Add(consts.HTTPTotal, count)
}

func (s *StateStorage) AddHttpSuccess(count int64) {
	s.store.Add(consts.HTTPSuccess, count)
}

func (s *StateStorage) AddHttpFailed(count int64) {
	s.store.Add(consts.HTTPFailed, count)
}

func (s *StateStorage) AddGrpcTotal(count int64) {
	s.store.Add(consts.GRPCTotal, count)
}

func (s *StateStorage) AddGrpcSuccess(count int64) {
	s.store.Add(consts.GRPCSuccess, count)
}

func (s *StateStorage) AddGrpcFailed(count int64) {
	s.store.Add(consts.GRPCFailed, count)
}

func (s *StateStorage) AddIosTotal(count int64) {
	s.store.Add(consts.IosTotal, count)
}

func (s *StateStorage) AddIosSuccess(count int64) {
	s.store.Add(consts.IosSuccess, count)
}

func (s *StateStorage) AddIosFailed(count int64) {
	s.store.Add(consts.IosFailed, count)
}

func (s *StateStorage) AddAndroidTotal(count int64) {
	s.store.Add(consts.IosTotal, count)
}

func (s *StateStorage) AddAndroidSuccess(count int64) {
	s.store.Add(consts.AndroidSuccess, count)
}

func (s *StateStorage) AddAndroidFailed(count int64) {
	s.store.Add(consts.AndroidFailed, count)
}

func (s *StateStorage) AddHuaweiTotal(count int64) {
	s.store.Add(consts.HuaweiTotal, count)
}

func (s *StateStorage) AddHuaweiSuccess(count int64) {
	s.store.Add(consts.HuaweiSuccess, count)
}

func (s *StateStorage) AddHuaweiFailed(count int64) {
	s.store.Add(consts.HuaweiFailed, count)
}

func (s *StateStorage) AddXiaomiTotal(count int64) {
	s.store.Add(consts.XiaomiTotal, count)
}

func (s *StateStorage) AddXiaomiSuccess(count int64) {
	s.store.Add(consts.XiaomiSuccess, count)
}

func (s *StateStorage) AddXiaomiFailed(count int64) {
	s.store.Add(consts.XiaomiFailed, count)
}

func (s *StateStorage) AddOppoTotal(count int64) {
	s.store.Add(consts.OppoTotal, count)
}

func (s *StateStorage) AddOppoSuccess(count int64) {
	s.store.Add(consts.OppoSuccess, count)
}

func (s *StateStorage) AddOppoFailed(count int64) {
	s.store.Add(consts.OppoFailed, count)
}

func (s *StateStorage) AddVivoTotal(count int64) {
	s.store.Add(consts.VivoTotal, count)
}

func (s *StateStorage) AddVivoSuccess(count int64) {
	s.store.Add(consts.VivoSuccess, count)
}

func (s *StateStorage) SetVivoSend(count int64) {
	s.store.Set(consts.VivoSend, count)
}

func (s *StateStorage) SetVivoReceive(count int64) {
	s.store.Set(consts.VivoReceive, count)
}

func (s *StateStorage) SetVivoDisplay(count int64) {
	s.store.Set(consts.VivoDisplay, count)
}

func (s *StateStorage) SetVivoClick(count int64) {
	s.store.Set(consts.VivoClick, count)
}

func (s *StateStorage) AddVivoFailed(count int64) {
	s.store.Add(consts.VivoFailed, count)
}

func (s *StateStorage) AddMeizuTotal(count int64) {
	s.store.Add(consts.MeizuTotal, count)
}

func (s *StateStorage) AddMeizuSuccess(count int64) {
	s.store.Add(consts.MeizuSuccess, count)
}

func (s *StateStorage) AddMeizuFailed(count int64) {
	s.store.Add(consts.MeizuFailed, count)
}

func (s *StateStorage) AddHonorTotal(count int64) {
	s.store.Add(consts.HonorTotal, count)
}

func (s *StateStorage) AddHonorSuccess(count int64) {
	s.store.Add(consts.HonorSuccess, count)
}

func (s *StateStorage) AddHonorFailed(count int64) {
	s.store.Add(consts.HonorFailed, count)
}

// GetTotalCount show counts of all notification.
func (s *StateStorage) GetTotalCount() int64 {
	//return s.store.Get(consts.HiPushTotal)
	var count int64
	count += s.GetAndroidTotal()
	count += s.GetIosTotal()
	count += s.GetHuaweiTotal()
	count += s.GetXiaomiTotal()
	count += s.GetVivoTotal()
	count += s.GetOppoTotal()
	count += s.GetMeizuTotal()
	count += s.GetHonorTotal()
	//count = s.GetHttpTotal() + s.GetGrpcTotal()
	return count
}

func (s *StateStorage) GetSuccessCount() int64 {
	var count int64
	count += s.GetAndroidSuccess()
	count += s.GetIosSuccess()
	count += s.GetHuaweiSuccess()
	count += s.GetXiaomiSuccess()
	count += s.GetVivoSuccess()
	count += s.GetOppoSuccess()
	count += s.GetMeizuSuccess()
	count += s.GetHonorSuccess()
	//count = s.GetHttpSuccess() + s.GetGrpcSuccess()
	return count
}

func (s *StateStorage) GetFailedCount() int64 {
	var count int64
	count += s.GetAndroidFailed()
	count += s.GetIosFailed()
	count += s.GetHuaweiFailed()
	count += s.GetXiaomiFailed()
	count += s.GetVivoFailed()
	count += s.GetOppoFailed()
	count += s.GetMeizuFailed()
	count += s.GetHonorFailed()
	//count = s.GetHttpFailed() + s.GetGrpcFailed()
	return count
}

func (s *StateStorage) GetSendCount() int64 {
	var count int64
	count += s.GetAndroidSend()
	count += s.GetIosSend()
	count += s.GetHuaweiSend()
	count += s.GetXiaomiSend()
	count += s.GetVivoSend()
	count += s.GetOppoSend()
	count += s.GetMeizuSend()
	count += s.GetHonorSend()
	return count
}

func (s *StateStorage) GetReceiveCount() int64 {
	var count int64
	count += s.GetAndroidReceive()
	count += s.GetIosReceive()
	count += s.GetHuaweiReceive()
	count += s.GetXiaomiReceive()
	count += s.GetVivoReceive()
	count += s.GetOppoReceive()
	count += s.GetMeizuReceive()
	count += s.GetHonorReceive()
	return count
}

func (s *StateStorage) GetDisplayCount() int64 {
	var count int64
	count += s.GetAndroidDisplay()
	count += s.GetIosDisplay()
	count += s.GetHuaweiDisplay()
	count += s.GetXiaomiDisplay()
	count += s.GetVivoDisplay()
	count += s.GetOppoDisplay()
	count += s.GetMeizuDisplay()
	count += s.GetHonorDisplay()
	return count
}

func (s *StateStorage) GetClickCount() int64 {
	var count int64
	count += s.GetAndroidClick()
	count += s.GetIosClick()
	count += s.GetHuaweiClick()
	count += s.GetXiaomiClick()
	count += s.GetVivoClick()
	count += s.GetOppoClick()
	count += s.GetMeizuClick()
	count += s.GetHonorClick()
	return count
}

func (s *StateStorage) GetHttpTotal() int64 {
	return s.store.Get(consts.HTTPTotal)
}

func (s *StateStorage) GetHttpSuccess() int64 {
	return s.store.Get(consts.HTTPSuccess)
}

func (s *StateStorage) GetHttpFailed() int64 {
	return s.store.Get(consts.HTTPFailed)
}

func (s *StateStorage) GetGrpcTotal() int64 {
	return s.store.Get(consts.GRPCTotal)
}

func (s *StateStorage) GetGrpcSuccess() int64 {
	return s.store.Get(consts.GRPCSuccess)
}

func (s *StateStorage) GetGrpcFailed() int64 {
	return s.store.Get(consts.GRPCFailed)
}

func (s *StateStorage) GetIosTotal() int64 {
	return s.store.Get(consts.IosTotal)
}

func (s *StateStorage) GetIosSuccess() int64 {
	return s.store.Get(consts.IosSuccess)
}

func (s *StateStorage) GetIosFailed() int64 {
	return s.store.Get(consts.IosFailed)
}

func (s *StateStorage) GetIosSend() int64 {
	return s.store.Get(consts.IosSend)
}

func (s *StateStorage) GetIosReceive() int64 {
	return s.store.Get(consts.IosReceive)
}

func (s *StateStorage) GetIosDisplay() int64 {
	return s.store.Get(consts.IosDisplay)
}

func (s *StateStorage) GetIosClick() int64 {
	return s.store.Get(consts.IosClick)
}

func (s *StateStorage) GetAndroidTotal() int64 {
	return s.store.Get(consts.AndroidTotal)
}

func (s *StateStorage) GetAndroidSuccess() int64 {
	return s.store.Get(consts.AndroidSuccess)
}

func (s *StateStorage) GetAndroidFailed() int64 {
	return s.store.Get(consts.AndroidFailed)
}

func (s *StateStorage) GetAndroidSend() int64 {
	return s.store.Get(consts.AndroidSend)
}

func (s *StateStorage) GetAndroidReceive() int64 {
	return s.store.Get(consts.AndroidReceive)
}

func (s *StateStorage) GetAndroidDisplay() int64 {
	return s.store.Get(consts.AndroidDisplay)
}

func (s *StateStorage) GetAndroidClick() int64 {
	return s.store.Get(consts.AndroidClick)
}

func (s *StateStorage) GetHuaweiTotal() int64 {
	return s.store.Get(consts.AndroidTotal)
}

func (s *StateStorage) GetHuaweiSuccess() int64 {
	return s.store.Get(consts.HuaweiSuccess)
}

func (s *StateStorage) GetHuaweiFailed() int64 {
	return s.store.Get(consts.HuaweiFailed)
}

func (s *StateStorage) GetHuaweiSend() int64 {
	return s.store.Get(consts.AndroidSend)
}

func (s *StateStorage) GetHuaweiReceive() int64 {
	return s.store.Get(consts.HuaweiReceive)
}

func (s *StateStorage) GetHuaweiDisplay() int64 {
	return s.store.Get(consts.HuaweiDisplay)
}

func (s *StateStorage) GetHuaweiClick() int64 {
	return s.store.Get(consts.HuaweiClick)
}

func (s *StateStorage) GetXiaomiTotal() int64 {
	return s.store.Get(consts.XiaomiTotal)
}

func (s *StateStorage) GetXiaomiSuccess() int64 {
	return s.store.Get(consts.XiaomiSuccess)
}

func (s *StateStorage) GetXiaomiFailed() int64 {
	return s.store.Get(consts.XiaomiFailed)
}

func (s *StateStorage) GetXiaomiSend() int64 {
	return s.store.Get(consts.XiaomiSend)
}

func (s *StateStorage) GetXiaomiReceive() int64 {
	return s.store.Get(consts.XiaomiReceive)
}

func (s *StateStorage) GetXiaomiDisplay() int64 {
	return s.store.Get(consts.XiaomiDisplay)
}

func (s *StateStorage) GetXiaomiClick() int64 {
	return s.store.Get(consts.XiaomiClick)
}

func (s *StateStorage) GetVivoTotal() int64 {
	return s.store.Get(consts.VivoTotal)
}

func (s *StateStorage) GetVivoSuccess() int64 {
	return s.store.Get(consts.VivoSuccess)
}

func (s *StateStorage) GetVivoFailed() int64 {
	return s.store.Get(consts.VivoFailed)
}

func (s *StateStorage) GetVivoSend() int64 {
	return s.store.Get(consts.VivoSend)
}

func (s *StateStorage) GetVivoReceive() int64 {
	return s.store.Get(consts.VivoReceive)
}

func (s *StateStorage) GetVivoDisplay() int64 {
	return s.store.Get(consts.VivoDisplay)
}

func (s *StateStorage) GetVivoClick() int64 {
	return s.store.Get(consts.VivoClick)
}

func (s *StateStorage) GetOppoTotal() int64 {
	return s.store.Get(consts.OppoTotal)
}

func (s *StateStorage) GetOppoSuccess() int64 {
	return s.store.Get(consts.OppoSuccess)
}

func (s *StateStorage) GetOppoFailed() int64 {
	return s.store.Get(consts.OppoFailed)
}

func (s *StateStorage) GetOppoSend() int64 {
	return s.store.Get(consts.OppoSend)
}

func (s *StateStorage) GetOppoReceive() int64 {
	return s.store.Get(consts.OppoReceive)
}

func (s *StateStorage) GetOppoDisplay() int64 {
	return s.store.Get(consts.OppoDisplay)
}

func (s *StateStorage) GetOppoClick() int64 {
	return s.store.Get(consts.OppoClick)
}

func (s *StateStorage) GetMeizuTotal() int64 {
	return s.store.Get(consts.MeizuTotal)
}

func (s *StateStorage) GetMeizuSuccess() int64 {
	return s.store.Get(consts.MeizuSuccess)
}

func (s *StateStorage) GetMeizuFailed() int64 {
	return s.store.Get(consts.MeizuFailed)
}

func (s *StateStorage) GetMeizuSend() int64 {
	return s.store.Get(consts.MeizuSend)
}

func (s *StateStorage) GetMeizuReceive() int64 {
	return s.store.Get(consts.MeizuReceive)
}

func (s *StateStorage) GetMeizuDisplay() int64 {
	return s.store.Get(consts.MeizuDisplay)
}

func (s *StateStorage) GetMeizuClick() int64 {
	return s.store.Get(consts.MeizuClick)
}

func (s *StateStorage) GetHonorTotal() int64 {
	return s.store.Get(consts.HonorTotal)
}

func (s *StateStorage) GetHonorSuccess() int64 {
	return s.store.Get(consts.HonorSuccess)
}

func (s *StateStorage) GetHonorFailed() int64 {
	return s.store.Get(consts.HonorFailed)
}

func (s *StateStorage) GetHonorSend() int64 {
	return s.store.Get(consts.HonorSend)
}

func (s *StateStorage) GetHonorReceive() int64 {
	return s.store.Get(consts.HonorReceive)
}

func (s *StateStorage) GetHonorDisplay() int64 {
	return s.store.Get(consts.HonorDisplay)
}

func (s *StateStorage) GetHonorClick() int64 {
	return s.store.Get(consts.HonorClick)
}
