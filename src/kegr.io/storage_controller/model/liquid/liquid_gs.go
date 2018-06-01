package liquid

// GetID getter
func (l *Liquid) GetID() string {
	return l.id
}

// SetID setter
func (l *Liquid) SetID(id string) {
	l.id = id
}

// GetFileHash getter
func (l *Liquid) GetFileHash() []byte {
	return l.fileHash
}

// SetFileHash setter
func (l *Liquid) SetFileHash(fileHash []byte) {
	l.fileHash = fileHash
}

// GetContent getter
func (l *Liquid) GetContent() []byte {
	return l.content
}

// SetContent setter
func (l *Liquid) SetContent(content []byte) {
	l.content = content
}

// GetSize getter
func (l *Liquid) GetSize() int64 {
	return l.size
}

// SetSize setter
func (l *Liquid) SetSize(size int64) {
	l.size = size
}

// GetLastUpdated getter
func (l *Liquid) GetLastUpdated() int64 {
	return l.lastUpdated
}

// SetLastUpdated setter
func (l *Liquid) SetLastUpdated(lastUpdated int64) {
	l.lastUpdated = lastUpdated
}

// IsDeleted getter
func (l *Liquid) IsDeleted() bool {
	return l.deleted
}

// SetDeleted setter
func (l *Liquid) SetDeleted(deleted bool) {
	if deleted {
		l.content = nil
	}
	l.deleted = deleted
}

// GetOptions getter
func (l *Liquid) GetOptions() IOptions {
	return l.options
}

// SetOptions setter
func (l *Liquid) SetOptions(options IOptions) {
	l.options = options
}
