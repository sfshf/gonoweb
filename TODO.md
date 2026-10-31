# TODO

## IAM 

- [nist iam](https://csrc.nist.gov/glossary/term/Identity_and_access_management?utm_source=chatgpt.com)
- [microsoft iam](https://learn.microsoft.com/en-us/entra/fundamentals/identity-fundamental-concepts?utm_source=chatgpt.com)

### [wiki -- IAM](https://en.wikipedia.org/wiki/Identity_and_access_management)

这个 Wiki 是 IAM（Identity and Access Management，身份与访问管理） 最经典的入门资料之一，也是很多企业（微软、Google、AWS、Okta、Auth0 等）的基础理念来源。

下面我（ChatGPT）按照 Wiki 的结构，结合企业实际项目（尤其是你最近在设计的 RBAC/IAM 微服务）进行翻译和讲解。

#### 一、IAM 是什么（Wiki 第一段）

Wiki 的定义：

```txt
Identity and access management (IAM) is a framework of policies and technologies to ensure that the right users have the appropriate access to technology resources.
```

翻译：

```txt
IAM（身份与访问管理）是一套**策略（Policy）+流程（Process）+技术（Technology）**组成的体系，它确保：

正确的人（Right User）

在

正确的时间（Right Time）

对

正确的资源（Right Resource）

拥有

正确的权限（Right Permission）
```

这是 IAM 最经典的一句话：

```txt
Right Person → Right Resource → Right Permission → Right Time
```

例如：

员工 Alice：

- 可以登录公司系统
- 可以访问 GitLab
- 可以访问 Jenkins
- 可以访问 Kubernetes
- 不能访问财务系统

IAM 就负责这一切。

#### 二、IAM 为什么存在

Wiki 提到：

IAM 属于：

- IT Security（信息安全）
- Data Management（数据管理）

因为企业越来越大：

假设一家企业：

```txt
10000 员工

↓

500 套系统

↓

100 万个账号

↓

几千万条权限
```

如果没有 IAM：

```txt
每个系统

自己维护账号

自己维护密码

自己维护权限
```

后果就是：

```txt
账号重复

密码不同

权限不同步

员工离职还保留账号

风险巨大
```

IAM 就是统一管理这些。

#### 三、Identity（身份）是什么

Wiki 专门解释了 Identity。

很多新人理解错。

Identity 不是用户名。

Identity 指的是：

```txt
一个实体（Entity）在数字世界中的身份。
```

实体可以是：

```txt
Person（人）

Machine（服务器）

Application（应用）

Service（服务）

API

Robot

IoT Device
```

所以：

IAM 管理的不只是用户。

还包括：

```txt
服务器身份

Docker

Kubernetes

Service Account

OAuth Client

机器人账号
```

现在越来越重要的是：

Machine Identity（机器身份）

甚至很多公司：

机器账号已经比员工账号还多。

#### 四、Digital Identity（数字身份）

Wiki 提到：

Digital Identity 包括：

```txt
Identity

+

Attributes
```

举例：

```txt
User

id = 1001

username = alice

email = xxx@gmail.com

department = IT

role = Admin

title = Engineer

phone

avatar

status
```

这些属性（Attribute）

共同组成：

Digital Identity。

所以：

Identity 不等于用户名。

Identity 是：

```txt
Identity

↓

Attributes

↓

Credential

↓

Role

↓

Permission
```

#### 五、Identity Management（身份管理）

Wiki 说：

Identity Management（IdM）

负责：

整个身份生命周期（Lifecycle）。

例如：

```txt
创建员工

↓

创建账号

↓

修改信息

↓

部门调整

↓

升职

↓

增加权限

↓

离职

↓

删除账号
```

这就是：

Identity Lifecycle。

很多企业叫：

Identity Governance。

#### 六、Access Management（访问管理）

Identity 有了以后：

下一步就是：

Access。

Wiki：

Access Management

负责：

控制：

谁可以访问什么。

例如：

```txt
Alice

↓

GitLab

↓

Read
```

或者：

```txt
Bob

↓

Jenkins

↓

Admin
```

访问管理回答两个问题：

```txt
Who？

能访问谁？

What？

能访问什么？
```

#### 七、IAM 的两个核心

Wiki 实际上强调：

IAM 包括：

##### 第一部分

Authentication（认证）

证明：

```txt
你是谁
```

例如：

```txt
用户名密码

MFA

短信

Google Authenticator

Face ID

Fingerprint

OIDC

OAuth Login
```

认证成功：

```txt
Alice

Verified
```

##### 第二部分

Authorization（授权）

认证以后：

继续判断：

```txt
Alice

↓

Can Access ?

↓

YES / NO
```

例如：

```txt
GET /api/users

↓

RBAC

↓

Admin

↓

Allow
```

所以：

认证：

```txt
Who are you？
```

授权：

```txt
What can you do？
```

很多新人会混淆。

这是 IAM 最重要的一点。

##### 八、IAM 管理哪些对象

Wiki 给出了几个对象。

不仅仅是：

User。

还包括：

```txt
User

Role

Permission

Device

Application

Service

Hardware

Certificate
```

现代 IAM：

越来越强调：

Machine Identity。

例如：

```txt
GitHub Action

↓

AWS IAM Role

↓

S3

↓

KMS
```

整个过程没有人参与。

也是 IAM。

##### 九、IAM 生命周期（最重要）

Wiki 把 IAM 分成：

Configuration Phase

Operation Phase

可以理解成：

第一阶段

配置

```txt
创建账号

创建角色

创建权限

绑定关系
```

例如：

```txt
Role

Admin

↓

Permission

user:create
```

这是：

Provisioning。

第二阶段：

运行。

例如：

用户登录：

```txt
Alice

↓

Login

↓

Authentication

↓

Authorization

↓

Access Resource
```

这是：

Operation。

所以：

IAM：

其实包括：

```txt
Provision

↓

Authentication

↓

Authorization

↓

Audit
```

##### 十、Wiki 中的几个重要概念

###### ① Identity

身份。

例如：

```txt
Alice
```


###### ② Credential

证明身份。

例如：

```txt
Password

Certificate

JWT

FIDO

Passkey
```

###### ③ Authentication

验证身份。

```txt
是不是 Alice？
```

###### ④ Authorization

授权。

```txt
Alice 能不能访问？
```

###### ⑤ Access Control

真正执行权限控制。

例如：

Go：

```txt
Gin Middleware

↓

Casbin

↓

Allow
```

###### ⑥ Audit

审计。

例如：

```txt
Alice

删除用户

时间

IP

UserAgent
```

全部记录下来。

##### 十一、Wiki 没有深入讲，但现代 IAM 已经发展出的能力

Wiki 偏基础定义，而现代企业 IAM（微软 Entra、Okta、Auth0、Keycloak 等）已经扩展出很多能力，例如：

|模块	|作用|
|--|--|
|SSO|	单点登录|
|MFA|	多因素认证|
|OAuth2|	授权|
|OIDC|	身份认证|
|SCIM|	用户同步|
|RBAC|	基于角色授权|
|ABAC|	基于属性授权|
|PBAC|	基于策略授权|
|Federation	|身份联邦|
|Identity Governance|	身份治理|
|Zero Trust	|零信任|
|PAM	|特权账号管理|

所以现在很多公司说 IAM，其实指的是一整套身份安全平台，而不是单纯的用户登录。

##### 十二、结合你的 RBAC 项目，该如何理解 IAM

根据你之前介绍的项目（包含 User、Role、Resource、Domain、Casbin 等模型），实际上你已经完成了 IAM 中授权（Authorization）层 的一部分。

可以把整体架构理解为：

```txt
                IAM Platform
                     │
     ┌───────────────┼───────────────┐
     │               │               │
 Identity       Authentication   Authorization
   (IdM)             (AuthN)         (AuthZ)
     │                 │               │
 User             OAuth2/OIDC      RBAC/Casbin
 Role             Password         ABAC
 Group            MFA              Permission
 Department       JWT              Resource
 Organization     SSO              Policy
```

也就是说：

- Identity（身份）：管理用户、组织、部门、角色等主体信息。
- Authentication（认证）：负责登录、OAuth2、OIDC、JWT、MFA 等。
- Authorization（授权）：负责 Casbin、RBAC、ABAC、资源权限等。
- Audit（审计）：负责记录所有登录、授权和资源访问行为。

这也是现代企业级 IAM（如 Keycloak、Okta、Microsoft Entra、Auth0）普遍采用的总体架构。

## microservices

- [microservices go](https://encore.dev/articles/microservices-go?utm_source=chatgpt.com)
- [Microservices in go](https://www.heroku.com/blog/microservices_in_go_using_go_kit/?utm_source=chatgpt.com)
- [go mocroservices](https://skoredin.pro/blog/tutorials/go-microservices?utm_source=chatgpt.com)
- [go food delivery microservices](https://github.com/mehdihadeli/go-food-delivery-microservices?utm_source=chatgpt.com)

### 当前项目如何拆分为微服务架构

你的这个项目实际上已经具备了**IAM（Identity & Access Management）**的雏形，而不仅仅是一个 RBAC 项目。

目前你的核心模型大概是：

```txt
User
Role
Domain（租户/组织）
Resource（菜单、按钮、API）
Casbin Policy
```

很多人在这里会犯一个错误：按照数据库表拆微服务。

例如：

```txt
❌ User Service
❌ Role Service
❌ Resource Service
❌ Casbin Service
❌ Domain Service
```

这是最典型的按 CRUD 拆服务，最后会导致：

- 每个请求跨 5 个 RPC
- 一个页面要调十几个接口
- 事务越来越复杂
- Casbin 成为所有服务依赖
- 性能越来越差

#### 我(ChatGPT)的建议：按业务能力（Bounded Context）拆，而不是按表拆

RBAC 本质上可以拆成下面几个领域。

```txt
                +----------------+
                | API Gateway    |
                +-------+--------+
                        |
        +---------------+---------------+
        |                               |
+-------v------+               +--------v-------+
| Identity     |               | Authorization  |
| Service      |               | Service        |
+-------+------+               +--------+-------+
        |                               |
        |                               |
+-------v-------+               +--------v--------+
| Organization  |               | Resource        |
| Service       |               | Service         |
+---------------+               +-----------------+
```

而不是

```txt
Role Service
User Service
Domain Service
Policy Service
```

##### 第一部分：Identity Service（身份中心）

负责

```txt
用户

登录

OAuth2

OIDC

JWT

Token

密码

MFA

Session
```

数据库

```txt
user

user_profile

user_password

user_login_history

refresh_token
```

这里不要放 Role。

Role 不属于身份。

##### 第二部分：Organization Service（组织中心）

负责

```txt
Tenant

Domain

Organization

Department

Group

Member
```

数据库

```txt
domain

department

group

user_domain

user_department
```

这里管理的是：

```txt
张三

属于：

XX公司

研发部

Go组
```

不是权限。

##### 第三部分：Authorization Service（授权中心）

这是整个 RBAC 最核心的。

负责

```txt
Role

Permission

Casbin

Policy

Role Binding

Permission Binding
```

数据库例如：

```txt
role

policy

role_policy

user_role

group_role

domain_role
```

这里只有：

```txt
谁

拥有什么角色

角色拥有哪些权限
```

不关心菜单长什么样。

##### 第四部分：Resource Service（资源中心）

你的 Resource 表建议单独独立出来。

Resource 不只是菜单。

Resource 应该是：

```txt
Menu

Button

API

Data Permission

Operation

Route

Widget

Feature Flag
```

例如：

```txt
resource
---------------------
id
type

MENU

BUTTON

API

PAGE

COLUMN

ACTION
```

例如：

```txt
用户管理
    |
    +--新增按钮
    |
    +--删除按钮
    |
    +--修改按钮

API

POST /user

DELETE /user

GET /user
```

全部都是 Resource。

Authorization 只引用 Resource ID。

不要复制资源。

##### 那 Casbin 放哪里？

Casbin 不应该成为独立服务。

应该属于

```txt
Authorization Service
```

里面。

例如：

```txt
Authorization

    Casbin Adapter

    Policy

    RBAC

    Permission Check

    Cache
```

外部服务：

```txt
CanAccess()

BatchCanAccess()

GetPermission()
```

而不是暴露 Casbin。

##### Resource 如何引用？

例如：

```txt
Resource

id=1001

MENU

/user

-------------

id=1002

BUTTON

user:add

-------------

id=1003

API

POST /user
```

Role

```txt
Admin

↓

Policy

↓

1001

1002

1003
```

永远引用 ResourceID。

##### Resource 不应该知道 Role

很多项目：

```txt
resource_role
```

这是错误设计。

正确应该：

```txt
Role

↓

Policy

↓

Resource
```

Resource 不知道谁拥有它。

##### API 服务怎么办？

真正拆微服务以后：


```txt
User Service

Order Service

Chat Service

Payment Service
```

每个服务都有自己的 API。

这些 API 应该自动注册 Resource。

例如：

```txt
POST /user

↓

启动时

↓

注册

↓

Resource Center
```

Resource Center

```txt
resource

id=10086

POST /user

service=user

module=user
```

这样后台菜单能自动发现 API。

##### 前端菜单怎么办？

很多公司：

```txt
Menu Service
```

其实没必要。

菜单也是 Resource。

例如：

```txt
Resource

MENU

Dashboard

MENU

User

MENU

Setting

BUTTON

Delete

BUTTON

Export

API

POST /user

API

DELETE /user
```

统一管理。

#### 我(ChatGPT)建议你的最终拆分

```txt
IAM
│
├── identity-service
│      登录
│      用户
│      Token
│      OAuth2
│      OIDC
│
├── organization-service
│      Domain
│      Tenant
│      Department
│      Group
│
├── authorization-service
│      Role
│      Casbin
│      Policy
│      Permission
│      Cache
│
├── resource-service
│      Menu
│      Button
│      API
│      Route
│      Widget
│      Metadata
│
└── gateway
       JWT
       鉴权
       路由
```

#### 如果是一个成熟的企业级 IAM，我还会进一步演进

考虑到你之前提到希望把它设计成可复用的 RBAC 微服务平台，我会采用更接近大型企业（如阿里云 RAM、腾讯 CAM、Keycloak、Auth0）的架构：

```txt
IAM Platform
│
├── identity-service        # 身份认证（Authentication）
├── authorization-service   # 权限决策（Authorization）
├── resource-service        # 资源注册中心
├── organization-service    # 租户/组织架构
├── audit-service           # 审计日志
├── notification-service    # 消息通知
├── policy-engine           # Casbin/OpenFGA 策略引擎
├── admin-console           # 后台管理
└── gateway                 # API 网关
```

其中 authorization-service 只负责权限关系和决策，resource-service 维护系统中所有菜单、按钮、API 等资源元数据，而业务微服务通过资源注册机制把自己的资源同步到资源中心。这样新增一个业务服务几乎不需要修改 IAM 的核心代码，扩展性会比直接围绕 user、role、resource 等表拆分高得多，也更符合领域驱动设计（DDD）的边界划分。


## tdd

TDD（Test Driven Development，测试驱动开发）是一种开发方法论，核心思想是：

先写测试，再写业务代码，用测试驱动设计，而不是写完代码再补测试。

很多优秀的 Go 开源项目（如 Kubernetes、Docker、Etcd、Prometheus 等）都大量采用了 TDD 或 TDD 的思想。

### TDD 的核心理念

传统开发流程：

```txt
需求
   │
   ▼
写业务代码
   │
   ▼
运行
   │
   ▼
发现 Bug
   │
   ▼
补测试
```

TDD 的开发流程：

```txt
需求
   │
   ▼
先写测试
   │
   ▼
测试失败（Red）
   │
   ▼
写最少的业务代码
   │
   ▼
测试通过（Green）
   │
   ▼
重构代码（Refactor）
   │
   ▼
测试仍然通过
``

这就是经典的 Red → Green → Refactor。

### TDD 三大循环

#### 第一阶段：Red（红灯）

先写测试。

例如：

```go
func TestAdd(t *testing.T) {
    if Add(1, 2) != 3 {
        t.Fatal("expect 3")
    }
}
```

运行：

```bash
go test
```

结果：

```txt
undefined: Add
```

失败。

这是故意的。

因为：

你还没有写业务代码。

#### 第二阶段：Green（绿灯）

写最简单的实现。

```go
func Add(a, b int) int {
    return 3
}
```

再次运行：

```bash
go test
```

结果：

```txt
PASS
```

虽然代码很烂：

```go
return 3
```

但是：

测试通过了。

继续增加测试：

```go
func TestAdd(t *testing.T) {

    if Add(1, 2) != 3 {
        t.Fatal()
    }

    if Add(2, 5) != 7 {
        t.Fatal()
    }
}
```

再次运行：

```txt
FAIL
```

于是修改：

```go
func Add(a, b int) int {
    return a + b
}
```

再次：

```txt
PASS
```

#### 第三阶段：Refactor（重构）

此时：

代码可以优化。

例如：

```go
func Add(a, b int) int {
    sum := a + b
    return sum
}
```

或者：

增加注释：

```go
// Add returns a+b.
func Add(a, b int) int {
    return a + b
}
```

重点：

重构以后测试必须仍然全部通过。

### 一个完整 Go 示例

假设实现一个字符串反转。

#### 第一步：先写测试

```go
package stringutil

import "testing"

func TestReverse(t *testing.T) {
    got := Reverse("hello")

    if got != "olleh" {
        t.Fatalf("expect olleh, got %s", got)
    }
}
```

运行：

```txt
undefined Reverse
```

#### 第二步：写业务代码

```go
package stringutil

func Reverse(s string) string {

    r := []rune(s)

    for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
        r[i], r[j] = r[j], r[i]
    }

    return string(r)
}
```

运行：

```txt
PASS
```

#### 第三步：增加更多测试

```go
func TestReverse(t *testing.T) {

    tests := []struct {
        input string
        want  string
    }{
        {"hello", "olleh"},
        {"abc", "cba"},
        {"", ""},
        {"中A文", "文A中"},
    }

    for _, tt := range tests {

        got := Reverse(tt.input)

        if got != tt.want {
            t.Fatalf("%v -> %v", tt.input, got)
        }
    }
}
```

这样就驱动你的实现不断完善。

### 为什么 TDD 能提升代码质量？

例如你设计：

```go
func CreateUser(name string)
```

写测试时：

```go
CreateUser("")
```

你会立即思考：

空名字怎么办？

于是增加：

```go
func TestEmptyName(t *testing.T)
```

然后想到：

```go
CreateUser("   ")
```

再想到：

```go
CreateUser("非常非常非常长...")
```

最后想到：

```go
CreateUser("重复用户名")
```

你会发现：

测试实际上帮助你设计接口，而不是验证代码。

因此很多人说：

TDD 驱动的是设计（Design），不是测试（Testing）。

### Go 中常用的测试方式

#### 1. Table Driven Test（最推荐）

Go 官方最推荐。

```go
func TestAdd(t *testing.T) {

    tests := []struct {
        a int
        b int
        want int
    }{
        {1, 2, 3},
        {2, 5, 7},
        {-1, 3, 2},
    }

    for _, tt := range tests {

        got := Add(tt.a, tt.b)

        if got != tt.want {
            t.Fatal()
        }
    }
}
```

几乎所有 Go 项目都这样写。

#### 2. 子测试（Subtest）

```go
for _, tt := range tests {

    t.Run(tt.name, func(t *testing.T) {

    })
}
```

好处：

失败时：

```txt
=== RUN TestAdd/negative

=== RUN TestAdd/zero
```

定位更快。

#### 3. Mock

例如：

```go
type UserRepo interface {

    Save(User) error

}
```

测试：

```go
repo := &MockRepo{}
service := NewUserService(repo)
```

无需真正连接数据库。

#### 4. Benchmark

```go
func BenchmarkAdd(b *testing.B) {

    for i := 0; i < b.N; i++ {
        Add(1, 2)
    }

}
```

运行：

```bash
go test -bench=.
```

#### 5. Fuzz Test（Go 1.18+）

```go
func FuzzReverse(f *testing.F) {

    f.Add("hello")

    f.Fuzz(func(t *testing.T, s string) {

        Reverse(Reverse(s))

    })

}
```

自动生成大量随机输入，帮助发现边界问题。

### Go 常用测试第三方库

|库|	作用	|推荐|
|-|-|-|
|testing（标准库）|	Go 官方测试框架|	⭐⭐⭐⭐⭐|
|github.com/stretchr/testify|	提供断言（assert/require）、Mock、Suite，Go 社区最常用	⭐⭐⭐⭐⭐
|go.uber.org/mock	|官方推荐的 Mock 框架（原 golang/mock 的后续）|	⭐⭐⭐⭐⭐|
|github.com/google/go-cmp/cmp	|深度比较复杂对象，比 reflect.DeepEqual 更灵活|	⭐⭐⭐⭐|
|github.com/onsi/ginkgo/v2 + github.com/onsi/gomega|	BDD 风格测试框架，适合大型项目|	⭐⭐⭐⭐|

例如使用 testify：

```go
func TestAdd(t *testing.T) {
    assert.Equal(t, 3, Add(1, 2))
    assert.Equal(t, 7, Add(3, 4))
}

```

相比手写 if 判断，更简洁、错误信息也更清晰。

### TDD 在真实项目中的实践（以 Go Web 项目为例）

一个典型的分层结构如下：

```txt
internal/
├── handler/
├── service/
├── repository/
├── model/
└── test/
```

推荐按下面的顺序进行 TDD：

1. 先写 Service 层测试：定义业务行为，例如“注册用户成功”“余额不足返回错误”。
2. Repository 使用 Mock：不依赖 MySQL、SQLite、Redis 等外部资源，让测试快速、稳定。
3. 实现 Service：编写最少代码让测试通过。
4. 补充 Handler 测试：验证 HTTP 请求、状态码和 JSON 返回。
5. 最后做集成测试：连接真实数据库、消息队列等验证整体流程。

### TDD 的核心要点总结

真正的 TDD 不只是“先写测试”，而是坚持以下原则：

- 测试先于实现：每增加一个功能，先写一个失败的测试。
- 一次只做一件事：每轮只增加一个测试，只写刚好让它通过的代码。
- 持续重构：测试通过后再优化代码结构，而不是一开始就追求完美设计。
- 测试描述业务行为：测试应该体现需求，例如“余额不足不能转账”，而不是“调用了某个私有函数”。
- 保持测试快速、可重复：大量使用接口和 Mock，避免单元测试依赖网络、数据库等外部环境。

对于你最近一直在做的 Go 后端、RBAC、CGO、SQLite、微服务 等项目，我建议把 Service 层作为 TDD 的起点。这一层业务逻辑最集中，也最容易通过接口和 Mock 隔离依赖，是 Go 项目中收益最高的 TDD 实践方式。
