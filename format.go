package statuses

import "encoding/json"

type Formatter func(Status) (string, error)

func ShortFormat(status Status) (string, error) {
	dirty := false
	props := make([]rune, 8)
	if status.Ahead >= 10 {
		dirty = true
		props[0], props[1] = '+', '*'
	} else if status.Ahead > 0 {
		dirty = true
		props[0], props[1] = '+', rune('0'+byte(status.Ahead))
	} else {
		props[0], props[1] = ' ', ' '
	}
	props[2] = ' '
	if status.Behind >= 10 {
		dirty = true
		props[3], props[4] = '-', '*'
	} else if status.Behind > 0 {
		dirty = true
		props[3], props[4] = '-', rune('0'+byte(status.Behind))
	} else {
		props[3], props[4] = ' ', ' '
	}
	props[5] = ' '
	if status.Modified {
		dirty = true
		props[6] = 'M'
	} else {
		props[6] = ' '
	}
	if status.Untracked {
		dirty = true
		props[7] = 'U'
	} else {
		props[7] = ' '
	}
	if !dirty {
		return "", nil
	}
	return string(props) + " " + status.Path, nil
}

func JSONFormat(status Status) (string, error) {
	buf, err := json.Marshal(status)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
