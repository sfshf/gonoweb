# gonoweb

This is a monolithic Go web application scaffold that uses Gin, GORM, and Casbin for the backend, and React.js and Next.js for the frontend.

## 高层设计

```mermaid
flowchart LR

UI[Web UI]
Gateway[API Gateway]
Auth[Auth Service]
User[User Service]
Role[Role Service]
Permission[Permission Service]
DB[(MySQL)]

UI --> Gateway

Gateway --> Auth
Gateway --> User
Gateway --> Role
Gateway --> Permission

Auth --> DB
User --> DB
Role --> DB
Permission --> DB
```

## 低层设计

## Backend

- [Gin](https://github.com/gin-gonic/gin)
- [Gorm](https://github.com/go-gorm/gorm)
- [Casbin](https://github.com/casbin/casbin)

## Frontend

- [React](https://github.com/facebook/react)
- [Nextjs](https://nextjs.org/)
- [HeroUI](https://www.heroui.com/docs/guide/introduction)

