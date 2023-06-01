package cronjob

import "sync"

type JobStatus int

const (
	JobStatusStop JobStatus = iota
	JobStatusRunning
	JobStatusRemove
)

func newProfile() *Profile {
	result := &Profile{
		jobStatus:   JobStatusRunning,
		processing:  0,
		overlapping: false,
	}
	return result
}

type Profile struct {
	mx          sync.Mutex
	entryID     EntryID
	jobStatus   JobStatus // 排程運行中
	processing  int       // job執行中
	overlapping bool      // 重疊執行
}

// ----------
// EntryID
// ----------

// GetEntryID get ID
func (p *Profile) GetEntryID() EntryID {
	p.mx.Lock()
	defer p.mx.Unlock()
	return p.entryID
}

func (p *Profile) setEntryID(entryID EntryID) {
	p.mx.Lock()
	p.entryID = entryID
	p.mx.Unlock()
}

// ----------
// JobStatus
// ----------

// Start 啟動排程
func (p *Profile) Start() *Profile {
	p.mx.Lock()
	if p.jobStatus != JobStatusRemove {
		p.jobStatus = JobStatusRunning
	}
	p.mx.Unlock()

	return p
}

// Stop 暫停排程
func (p *Profile) Stop() *Profile {
	p.mx.Lock()
	if p.jobStatus != JobStatusRemove {
		p.jobStatus = JobStatusStop
	}
	p.mx.Unlock()

	return p
}

func (p *Profile) remove() *Profile {
	p.mx.Lock()
	p.jobStatus = JobStatusRemove
	p.mx.Unlock()
	return p
}

func (p *Profile) GetJobStatus() JobStatus {
	p.mx.Lock()
	defer p.mx.Unlock()

	return p.jobStatus
}

// IsRunning 排程是否啟用
func (p *Profile) IsRunning() bool {
	return p.isRunning()
}

func (p *Profile) isRunning() bool {
	p.mx.Lock()
	isRunning := p.jobStatus == JobStatusRunning
	p.mx.Unlock()
	return isRunning
}

// ----------
// processing setting
// ----------

// Processing 排程執行中
func (p *Profile) processingAdd() {
	p.mx.Lock()
	p.processing++
	p.mx.Unlock()
}

func (p *Profile) processingDone() {
	p.mx.Lock()
	p.processing--
	p.mx.Unlock()
}

// ProcessingCount 排程執行中的數量
func (p *Profile) ProcessingCount() int {
	p.mx.Lock()
	count := p.processing
	p.mx.Unlock()
	return count
}

// IsProcessing 排程是否執行中
func (p *Profile) IsProcessing() bool {
	return p.isProcessing()
}

func (p *Profile) isProcessing() bool {
	p.mx.Lock()
	isProcessing := p.processing > 0
	p.mx.Unlock()
	return isProcessing
}

// ----------
// overlapping setting
// ----------

// Overlapping 在同一時間可否重覆執行
func (p *Profile) Overlapping(b bool) *Profile {
	p.mx.Lock()
	p.overlapping = b
	p.mx.Unlock()
	return p
}

func (p *Profile) isOverlapping() bool {
	p.mx.Lock()
	isOverlapping := p.overlapping
	p.mx.Unlock()
	return isOverlapping
}
