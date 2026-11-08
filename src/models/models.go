package models

func All() []any {
	return []any{
		&User{},
		&Session{},
		&Friendship{},
		&List{},
		&ListItem{},
		&Challenge{},
		&ChallengeItem{},
	}
}
