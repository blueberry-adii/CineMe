package booking

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Book() (Booking, error) {
	return Booking{}, nil
}
