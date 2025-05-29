### API framework:
| Framework                                     | Pros                                                                                            | Cons                                                    | Suitability                                 |
| --------------------------------------------- | ----------------------------------------------------------------------------------------------- | ------------------------------------------------------- | ------------------------------------------- |
| **[Chi](https://github.com/go-chi/chi)**      | - Minimal & idiomatic Go<br>- Built-in middleware<br>- Tiny footprint<br>- Easy to reason about | - No fancy codegen or decorators<br>- Manual validation | **✅ Best for modular, custom servers**      |
| **[Gin](https://github.com/gin-gonic/gin)**   | - Fast<br>- Lots of features<br>- Good docs<br>- JSON binding, validation                       | - Non-idiomatic in places<br>- Can feel too magical     | 👍 Great if you want more built-in features |
| **[Fiber](https://github.com/gofiber/fiber)** | - ExpressJS-like<br>- Fastest raw perf<br>- Easy to learn                                       | - Less idiomatic<br>- Less flexible for advanced use    | 😐 Better for simpler projects              |
| **Echo**                                      | - Feature-rich<br>- Built-in middleware<br>- Structured                                         | - Slightly heavier<br>- Can feel opinionated            | 👍 Solid middle-ground                      |

 - `Chi`

### WebDAV
- `golang.org/x/net/webdav` - Maintained by golang team

### Authorization
Webdav clients support only the basic auth aka. User + Password
JWT based auth might be implemented for webapps, but it is optional.

### Database
| Want to...                     | Use                      |
| ------------------------------ | ------------------------ |
| Full-featured, easy ORM        | **GORM**                 |
| Type-safe, schema-first        | **Ent** or **SQLBoiler** |
| Stick close to SQL, fast       | **pgx** + `sqlc`         |
| Avoid raw SQL but no heavy ORM | **Squirrel**             |

- `squirrel` - due to being used to sqlalchemy
- golang-migrate