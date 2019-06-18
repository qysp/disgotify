// Package disgord provides Go bindings for the documented Discord API. And allows for a stateful Client using the Session interface, with the option of a configurable caching system or bypass the built-in caching logic all together.
//
// Getting started
//
// Create a DisGord session to get access to the REST API and socket functionality. In the following example, we listen for new messages and write a "hello" message when our handler function gets fired.
//
// Session interface: https://godoc.org/github.com/andersfylling/disgord/#Session
//  discord, err := disgord.NewClient(&disgord.Config{
//    BotToken: "my-secret-bot-token",
//  })
//  if err != nil {
//    panic(err)
//  }
//
//  // listen for incoming messages and reply with a "hello"
//  discord.On(event.MessageCreate, func(session disgord.Session, evt *disgord.MessageCreate) {
//      evt.Message.RespondString("hello")
//  })
//
//  // connect to the socket API to receive events
//  err = discord.Connect()
//  if err != nil {
//      panic(err)
//  }
//  discord.DisconnectOnInterrupt()
//
// If you want some logic to fire when the bot is ready (all shards has received their ready event), please use the Ready method.
//  // ...
//  discord.Ready(func() {
//  	fmt.Println("READY NOW!")
//  })
//  // ...
//
//
//
// Listen for events using channels
//
// Disgord also provides the option to listen for events using a channel, instead of registering a handler. However, before using the event channel, you must notify disgord that you care about the event (this is done automatically in the event handler registration).
//  session.AcceptEvent(event.MessageCreate) // alternative: disgord.EvtMessageCreate
//  session.AcceptEvent(event.MessageUpdate)
//  for {
//      var message *disgord.Message
//      var status string
//      select {
//      case evt, alive := <- session.EventChannels().MessageCreate()
//          if !alive {
//              return
//          }
//          message = evt.Message
//          status = "created"
//      case evt, alive := <- session.EventChannels().MessageUpdate()
//          if !alive {
//              return
//          }
//          message = evt.Message
//          status = "updated"
//      }
//
//      fmt.Printf("A message from %s was %s\n", message.Author.Mention(), status)
//      // output example: "A message from @Anders was created"
//  }
//
// Optimizing your cache logic
//
// > Note: if you create a CacheConfig you don't have to set every field. All the CacheAlgorithms are default to LFU when left blank.
//
// A part of Disgord is the control you have; while this can be a good detail for advanced users, we recommend beginners to utilise the default configurations (by simply not editing the configuration).
// Here we pass the cache config when creating the session to access to the different cache replacement algorithms, lifetime settings, and the option to disable different cache systems.
//  discord, err := disgord.NewClient(&disgord.Config{
//    BotToken: "my-secret-bot-token",
//    Cache: &disgord.CacheConfig{
//              Mutable: false, // everything going in and out of the cache is deep copied
//				// setting Mutable to true, might break your program as this is experimental.
//
//              DisableUserCaching: false, // activates caching for users
//              UserCacheLifetime: time.Duration(4) * time.Hour, // removed from cache after 9 hours, unless updated
//              UserCacheAlgorithm: disgord.CacheAlgLFU,
//
//              DisableVoiceStateCaching: true, // don't cache voice states
//              // VoiceStateCacheLifetime  time.Duration
//              // VoiceStateCacheAlgorithm string
//
//              DisableChannelCaching: false,
//              ChannelCacheLifetime: 0, // lives forever
//              ChannelCacheAlgorithm: disgord.CacheAlgLFU, // lfu (Least Frequently Used)
//
//				GuildCacheAlgorithm: disgord.CacheAlgLFU, // no limit set, so the strategy to replace entries is not used
//           },
//  })
//
// If you just want to change a specific field, you can do so. By either calling the disgord.DefaultCacheConfig which gives you a Cache configuration designed by DisGord. Or you can set specific fields in a new CacheConfig since the different Cache Strategies are automatically set to LFU if missing.
// 	&disgord.Config{}
// Will automatically become
//  &disgord.Config{
//  	UserCacheAlgorithm: disgord.CacheAlgLFU,
//		VoiceStateCacheAlgorithm disgord.CacheAlgLFU,
//		ChannelCacheAlgorithm: disgord.CacheAlgLFU,
//		GuildCacheAlgorithm: disgord.CacheAlgLFU,
//  }
//
// And writing
//  &disgord.Config{
//  	UserCacheAlgorithm: disgord.CacheAlgLRU,
//		VoiceStateCacheAlgorithm disgord.CacheAlgLRU,
//  }
// Becomes
//  &disgord.Config{
//  	UserCacheAlgorithm: disgord.CacheAlgLRU, // unchanged
//		VoiceStateCacheAlgorithm disgord.CacheAlgLRU,  // unchanged
//		ChannelCacheAlgorithm: disgord.CacheAlgLFU,
//		GuildCacheAlgorithm: disgord.CacheAlgLFU,
//  }
//
// > Note: Disabling caching for some types while activating it for others (eg. disabling channels, but activating guild caching), can cause items extracted from the cache to not reflect the true discord state.
//
// Example, activated guild but disabled channel caching: The guild is stored to the cache, but it's channels are discarded. Guild channels are dismantled from the guild object and otherwise stored in the channel cache to improve performance and reduce memory use. So when you extract the cached guild object, all of the channel will only hold their channel ID, and nothing more.
//
//
// Immutable and concurrent accessible cache
//
// The option CacheConfig.Immutable can greatly improve performance or break your system. If you utilize channels or you need concurrent access, the safest bet is to set immutable to `true`. While this is slower (as you create deep copies and don't share the same memory space with variables outside the cache), it increases reliability that the cache always reflects the last known Discord state.
// If you are uncertain, just set it to `true`. The default setting is `true` if `disgord.Cache.CacheConfig` is `nil`.
//
//
// Bypass the built-in REST cache
//
// Whenever you call a REST method from the Session interface; the cache is always checked first. Upon a cache hit, no REST request is executed and you get the data from the cache in return. However, if this is problematic for you or there exist a bug which gives you bad/outdated data, you can bypass it by using the REST functions directly. Remember that this will not update the cache for you, and this needs to be done manually if you depend on the cache.
//  // get a user using the Session implementation (checks cache, and updates the cache on cache miss)
//  user, err := session.GetUser(userID)
//
//  // bypass the cache checking. Same function name, but is found in the disgord package, not the session interface.
//  user, err := disgord.GetUser(userID)
//
// Manually updating the cache
//
// If required, you can access the cache and update it by hand. Note that this should not be required when you use the Session interface.
//  user, err := disgord.GetUser(userID)
//  if err != nil {
//      return err
//  }
//
//  // update the cache
//  cache := discord.Cache()
//  err = cache.Update(disgord.UserCache, user)
//  if err != nil {
//      return err
//  }
//
//
// Build tags
//
// `disgord_diagnosews` will store all the incoming and outgoing json data as files in the directory "diagnose-report/packets". The file format is as follows: unix_clientType_direction_shardID_operationCode_sequenceNumber[_eventName].json
//
// `json-std` switches out jsoniter with the json package from the std libs.
//
// `disgord_removeDiscordMutex` replaces mutexes in discord structures with a empty mutex; removes locking behaviour and any mutex code when compiled.
//
// `disgord_parallelism` activates built-in locking in discord structure methods. Eg. Guild.AddChannel(*Channel) does not do locking by default. But if you find yourself using these discord data structures in parallel environment, you can activate the internal locking to reduce race conditions. Note that activating `disgord_parallelism` and `disgord_removeDiscordMutex` at the same time, will cause you to have no locking as `disgord_removeDiscordMutex` affects the same mutexes.
//
//
// Saving and Deleting Discord data
//
// > Note: when using SaveToDiscord(...) make sure the object reflects the Discord state. Calling Save on default values might overwrite or reset the object at Discord, causing literally.. Hell.
//
// You might have seen the two methods in the session interface: SaveToDiscord(...) and DeleteFromDiscord(...).
// This are as straight forward as they sound. Passing a discord data structure into one of them executes their obvious behavior; to either save it to Discord, or delete it.
//  // create a new role and give it certain permissions
//  role := disgord.Role{}
//  role.Name = "Giraffes"
//  role.GuildID = guild.ID // required, for an obvious reason
//  role.Permissions = disgord.ManageChannelsPermission | disgord.ViewAuditLogsPermission
//  err := session.SaveToDiscord(&role)
//
// You know what.. Let's just remove the role
//  err := session.DeleteFromDiscord(&role)
//
package disgord

import (
	"fmt"

	"github.com/andersfylling/disgord/constant"
	"github.com/andersfylling/snowflake/v3"
)

// LibraryInfo returns name + version
func LibraryInfo() string {
	return fmt.Sprint(constant.Name, constant.Version)
}

// Wrapper for github.com/andersfylling/snowflake
// ------------------

// Snowflake twitter snowflake identification for Discord
type Snowflake = snowflake.Snowflake

// GetSnowflake see snowflake.GetSnowflake
func GetSnowflake(v interface{}) (Snowflake, error) {
	s, err := snowflake.GetSnowflake(v)
	return Snowflake(s), err
}

// NewSnowflake see snowflake.NewSnowflake
func NewSnowflake(id uint64) Snowflake {
	return Snowflake(snowflake.NewSnowflake(id))
}

// ParseSnowflakeString see snowflake.ParseSnowflakeString
func ParseSnowflakeString(v string) Snowflake {
	return Snowflake(snowflake.ParseSnowflakeString(v))
}

func newErrorMissingSnowflake(message string) *ErrorMissingSnowflake {
	return &ErrorMissingSnowflake{
		info: message,
	}
}

// ErrorMissingSnowflake used by methods about to communicate with the Discord API. If a snowflake value is required
// this is used to identify that you must set the value before being able to interact with the Discord API
type ErrorMissingSnowflake struct {
	info string
}

func (e *ErrorMissingSnowflake) Error() string {
	return e.info
}

func newErrorEmptyValue(message string) *ErrorEmptyValue {
	return &ErrorEmptyValue{
		info: message,
	}
}

// ErrorEmptyValue when a required value was set as empty
type ErrorEmptyValue struct {
	info string
}

func (e *ErrorEmptyValue) Error() string {
	return e.info
}
