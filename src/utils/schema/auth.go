package schema

import (
	"bucket-list-be/src/utils"
	"regexp"
	"strings"

	z "github.com/Oudwins/zog"
)

var SignupSchema = z.Struct(z.Shape{
	"Name": z.String().
		Trim().
		Min(2).
		Max(100).
		Required(),
	"Username": z.String().
		Trim().
		Transform(func(value *string, ctx z.Ctx) error {
			*value = strings.ToLower(*value)
			return nil
		}).
		Min(3).
		Max(30).
		Match(
			regexp.MustCompile(`^[a-z0-9_]+$`),
			z.Message("Username can contain only lowercase letters, numbers, and underscores"),
		).
		Required(),
	"Email": z.String().
		Trim().
		Transform(func(value *string, ctx z.Ctx) error {
			*value = strings.ToLower(*value)
			return nil
		}).
		Email().
		Max(254).
		Required(),
	"Password": z.String().
		Required().
		Test(z.Test[*string]{
			Func: func(value *string, ctx z.Ctx) {
				if err := utils.ValidatePassword(*value); err != nil {
					ctx.AddIssue(ctx.Issue().SetMessage(err.Error()))
				}
			},
		}),
})
