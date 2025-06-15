package mockRepo

import "github.com/stretchr/testify/mock"

func (m *MovieRepository) ClearAll() {
	m.Mock = mock.Mock{}
}
