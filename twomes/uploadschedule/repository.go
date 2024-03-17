package uploadschedule

type UploadScheduleRepository interface {
	Create(UploadSchedule) (UploadSchedule, error)
	Delete(UploadSchedule) error
}
