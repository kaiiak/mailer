package main

type mockSender struct {
	config *mailConfig
}

func newMockSender(config *mailConfig) *mockSender {
	return &mockSender{config}
}

func (sender *mockSender) Send(m email) error {
	if _, err := parseMailList(m.recipient()); err != nil {
		return err
	}
	return nil
}
