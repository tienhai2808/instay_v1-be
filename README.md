## GET STARTED

**Require Go version 1.25+**

```bash
make run
```

### Project Structure

```
├── 📁 cmd
│   └── 📁 api
│       └── 🐹 main.go
├── 📁 configs
│   └── ⚙️ main.yaml
├── 📁 internal
│   ├── 📁 common
│   │   ├── 🐹 constants.go
│   │   ├── 🐹 errors.go
│   │   ├── 🐹 mapper.go
│   │   └── 🐹 utils.go
│   ├── 📁 config
│   │   └── 🐹 main_config.go
│   ├── 📁 container
│   │   ├── 🐹 auth_container.go
│   │   ├── 🐹 file_container.go
│   │   ├── 🐹 main_container.go
│   │   └── 🐹 user_container.go
│   ├── 📁 handler
│   │   ├── 🐹 auth_handler.go
│   │   ├── 🐹 file_handler.go
│   │   └── 🐹 user_handler.go
│   ├── 📁 initialization
│   │   ├── 🐹 logger.go
│   │   ├── 🐹 postgresql.go
│   │   ├── 🐹 rabbitmq.go
│   │   ├── 🐹 redis.go
│   │   ├── 🐹 s3.go
│   │   └── 🐹 snowflake.go
│   ├── 📁 middleware
│   │   └── 🐹 authentication.go
│   ├── 📁 model
│   │   └── 🐹 user_model.go
│   ├── 📁 provider
│   │   ├── 📁 jwt
│   │   │   └── 🐹 jwt.go
│   │   ├── 📁 mq
│   │   │   └── 🐹 message_queue.go
│   │   └── 📁 smtp
│   │       ├── 📁 templates
│   │       │   └── 🌐 auth.html
│   │       └── 🐹 smtp.go
│   ├── 📁 repository
│   │   ├── 📁 implement
│   │   │   └── 🐹 user_repo_impl.go
│   │   └── 🐹 user_repository.go
│   ├── 📁 router
│   │   ├── 🐹 auth_router.go
│   │   ├── 🐹 file_router.go
│   │   └── 🐹 user_router.go
│   ├── 📁 server
│   │   └── 🐹 server.go
│   ├── 📁 service
│   │   ├── 📁 implement
│   │   │   ├── 🐹 auth_svc_impl.go
│   │   │   ├── 🐹 file_svc_impl.go
│   │   │   └── 🐹 user_svc_impl.go
│   │   ├── 🐹 auth_service.go
│   │   ├── 🐹 file_service.go
│   │   └── 🐹 user_service.go
│   ├── 📁 types
│   │   ├── 🐹 data.go
│   │   ├── 🐹 request.go
│   │   └── 🐹 response.go
│   └── 📁 worker
│       └── 🐹 email_worker.go
├── 📁 logs
│   └── 📄 app.log
├── 📁 pkg
│   ├── 📁 bcrypt
│   │   └── 🐹 bcrypt.go
│   └── 📁 snowflake
│       └── 🐹 snowflake.go
├── ⚙️ .gitignore
├── 📄 Makefile
├── 📝 README.md
├── 📄 go.mod
└── 📄 go.sum
```