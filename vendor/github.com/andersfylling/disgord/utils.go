package disgord

import (
	"errors"
	"strings"
)

// GetShardForGuildID converts a GuildID into a ShardID for correct retrieval of guild information
func GetShardForGuildID(guildID Snowflake, shardCount uint) (shardID uint) {
	return uint(guildID>>22) % shardCount
}

//////////////////////////////////////////////////////
//
// Validators
//
//////////////////////////////////////////////////////

func ValidateHandlerInputs(inputs ...interface{}) (err error) {
	var i int
	var ok bool

	// make sure that middlewares are only at beginning
	for j := i; j < len(inputs); j++ {
		if _, ok = inputs[j].(Middleware); ok {
			if j != i {
				return errors.New("middlewares can only be in the beginning. Grouped together")
			}
			i++
		}
	}

	// there should now be N handlers, 0 < N.
	if len(inputs) <= i {
		return errors.New("missing handler(s)")
	}

	// check for extra controllers
	for j := len(inputs) - 2; j >= i; j-- {
		if _, ok = inputs[j].(HandlerCtrl); ok {
			return errors.New("a handlerCtrl's can only be at the end of the definition. Expected a handler")
		}
	}

	// TODO: test for all handler types?

	return nil
}

// https://discordapp.com/developers/docs/resources/user#avatar-data
func validAvatarPrefix(avatar string) (valid bool) {
	if avatar == "" {
		return false
	}

	construct := func(encoding string) string {
		return "data:image/" + encoding + ";base64,"
	}

	if len(avatar) < len(construct("X")) {
		return false // missing base64 declaration
	}

	encodings := []string{
		"jpeg", "png", "gif",
	}
	for _, encoding := range encodings {
		prefix := construct(encoding)
		if strings.HasPrefix(avatar, prefix) {
			valid = len(avatar)-len(prefix) > 0 // it has content
			break
		}
	}

	return true
}

// ValidateUsername uses Discords rule-set to verify user-names and nicknames
// https://discordapp.com/developers/docs/resources/user#usernames-and-nicknames
//
// Note that not all the rules are listed in the docs:
//  There are other rules and restrictions not shared here for the sake of spam and abuse mitigation, but the
//  majority of users won't encounter them. It's important to properly handle all error messages returned by
//  Discord when editing or updating names.
func ValidateUsername(name string) (err error) {
	if name == "" {
		return errors.New("empty")
	}

	// attributes
	length := len(name)

	// Names must be between 2 and 32 characters long.
	if length < 2 {
		err = errors.New("name is too short")
	} else if length > 32 {
		err = errors.New("name is too long")
	}
	if err != nil {
		return err
	}

	// Names are sanitized and trimmed of leading, trailing, and excessive internal whitespace.
	if name[0] == ' ' {
		err = errors.New("contains whitespace prefix")
	} else if name[length-1] == ' ' {
		err = errors.New("contains whitespace suffix")
	} else {
		last := name[1]
		for i := 2; i < length-1; i++ {
			if name[i] == ' ' && last == name[i] {
				err = errors.New("contains excessive internal whitespace")
				break
			}
			last = name[i]
		}
	}
	if err != nil {
		return err
	}

	// Names cannot contain the following substrings: '@', '#', ':', '```'
	illegalChars := []string{
		"@", "#", ":", "```",
	}
	for _, illegalChar := range illegalChars {
		if strings.Contains(name, illegalChar) {
			err = errors.New("can not contain the character " + illegalChar)
			return err
		}
	}

	// Names cannot be: 'discordtag', 'everyone', 'here'
	illegalNames := []string{
		"discordtag", "everyone", "here",
	}
	for _, illegalName := range illegalNames {
		if name == illegalName {
			err = errors.New("the given username is illegal")
			return err
		}
	}

	return nil
}

func validateChannelName(name string) (err error) {
	if name == "" {
		return errors.New("empty")
	}

	// attributes
	length := len(name)

	// Names must be of length of minimum 2 and maximum 100 characters long.
	if length < 2 {
		err = errors.New("name is too short")
	} else if length > 100 {
		err = errors.New("name is too long")
	}
	if err != nil {
		return err
	}

	return nil
}
