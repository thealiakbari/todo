package enum

type Topic string

func (t Topic) String() string {
	return string(t)
}
