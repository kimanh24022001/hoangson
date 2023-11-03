**NOTE (duong)**: As we journey through the swirling mists of creation, this masterpiece is still being woven with threads of enchantment and whispers of forgotten legends; tread lightly, for its magic is still taking shape.

# Welcome to the Smatyx Quest Realms! ✨

Embark on an epic quest with us! Within the hallowed halls of this repository, you will uncover a myriad of directories. Each directory is a realm unto itself, harboring treasures and secrets for those deemed worthy:

- `cmd`: A trove of mystical `main.go` files, each imbued with its own unique power and purpose.
- `config`: The sacred heart of our universe, holding the ancient config file that binds our world together. (Prophecy foretells a time when each directory will wield its own config file.)
- `continuous`: A concealed sanctum where the revered script for our CI-CD pipeline resides, ever-ready to be awakened.
- `shared`: A vast library housing ancient scrolls of source code from our web framework, beckoning those wise enough to harness their wisdom.
- `hoangson`: Hallowed grounds dedicated to the source code that powers Hoang Son's digital realm.
- `misc`: An enigmatic chamber filled with a vast array of mysterious artifacts.

# Embark on the Journey! 🧭

Prepare for your quest by arming yourself with the mighty Go version 1.21.3. Seize your weapon from the armory at [Go Downloads](https://go.dev/dl/).

With your newfound power, conjure the server to life with this sacred incantation:

```sh
cd hoangson/ && go run .
```

# Entities & Migrations: The Ancient Rituals 📜

In the hallowed ground of `shared/entities`, you will find Entities - ancient spirits representing the state of our application.

To grant these spirits a physical form in the database realm, inscribe your structures with the ancient runes:

```go
//entity:table
type Migration struct {
	Id          *types.String
	QueryText   *types.String
	AppliedTime *types.Time
}
```

The incantations to summon new tables are as follows:

```go
//entity:table
//entity:table(new_table_name)
//entity:table(old_table_name -> new_table_name)
```

To unveil the hidden structures of existing entities, recite the spell:

```
go run cmd/make_migrate/main.go
```

Finally, bring your newly crafted structures into the physical realm by chanting:

```
go run cmd/do_migrate/main.go
```

# The Final Act 🌟

Should you find yourself lost in the mists of confusion or simply seeking an ally, I am but a Slack message or email away.

Here's to a grand adventure filled with code and wonder!

Yours in magic and mystery,

Dương
