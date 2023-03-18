# Gofermart

## Use cases

```puml
@startuml
skinparam actorStyle awesome

legend right
| Color           | Mean                         |
|                  | to implement |
|<#lightyellow>    | are not in the task                  |
endlegend

left to right direction

actor "User" as U

package "Система лояльности" as M {
    usecase "регистрация/авторизация" as M_RA
    usecase "отправить номер заказа"as M_SON
    usecase "получить статус начисления по заказу" as M_GOS
    usecase "списать балы" as M_DB
    usecase "получить все списания" as M_DH
}

package "Система расчетов балов(aka accrual)" as ACR #lightyellow {
    usecase "Зарегестрировать номер заказа" as A_RON
    usecase "Расчитать балы лояльности по заказу" as A_CB
}
package "Online shop" as OS #lightyellow {
    usecase "Купить" as OS_B
}

U ..> M_RA : 3
U ..> M_SON : 4
M_SON ..> A_CB : 5
U ..> M_GOS : 6
U ..> M_DB : 7
U ..> OS_B : 1
OS_B ..> A_RON : 2
U ..> M_DH : 8

@enduml
```

## Model

```puml
@startuml

title model
namespace Gofermart {
    entity "User" as U {
        id int
        login string
        pwd string
    }
    entity "Order" as O1 {
        id int
        userId int
        orderNumber string
        status int
        accrualStatus: int
        accrual int
        created_at datetime
    }
'    entity "BalanceAction" as BA {
'        value uint
'        type int (enum: \n\
'        1 = пополнение на основе\n\
'        добавление заказа;\n\
'        2 = списание по инициативе\n\
'        пользователя.\n)
'        orderId int (default = 0 - нет заказа)
'    }
    
    entity "Withdrawal" as W {
        id int
        userId int
        amount float
        created_at datetime
    }
}

namespace Accrual {
    entity "Order" as O2 {
        + orderNumber string
        + status string
        + accrual null|int
    }
}

U "1" *-d- "many" O1 : create >
'U "1" *-d- "many" BA : read >
'O1 "1" -l- "1" BA : add >
U "1" *-d- "many" W : create >

O1 "1"--"1" O2 : read >

@enduml
```

## Classes

```puml
@startuml

package "Business" as BL {
    package "Services" as BL.S{
        class "OrderService" as BL.S.OS  {
            + AddOrder()
            + GetAllOrdersByUserID()
            + UpdateOrdersStatuses()
            + GetSumAccrualByUserID()
        }
        class "WithdrawalService" as BL.S.WS {
            + AddWithdrawal()
            + GetAllWithdrawalsByUserID()
            + GetSumAmountByUserID()
        }
        class "UserService" as BL.S.US {
            + CreateNewUser()
            + GetUserIfPwdValid()
            + GetBalance()
        }
        
        package "Contract" as BL.S.C {
            interface "OrderRepositoryInterface" as BL.S.C.ORI {}
            interface "UserRepositoryInterface" as BL.S.C.URI {}
            interface "WithdrawalRepositoryInterface" as BL.S.C.WRI {}
            BL.S.C.URI -[hidden]d-> BL.S.C.WRI
            BL.S.C.WRI -[hidden]-> BL.S.C.ORI
         }
         BL.S.US o-- BL.S.WS
         BL.S.US o--- BL.S.OS
         
         BL.S.OS -ro  BL.S.C.ORI
         BL.S.WS -ro  BL.S.C.WRI
         BL.S.US -ro  BL.S.C.URI
         
    }
    package "Entities" as BL.E {
'        class "User" as U {}
'        class "Order" as O {}
'        class "Withdrawal" as W {}
'        U *-- O 
'        note left on link : Заказ был\nдобавлен\nпользователем
'        U *-- W 
'        note left on link : Списание было\nдобавлено\nпользователем
    }
}
class "UserController" as UC {
    login()
    register()
    orderList()
    addOrder()
    balance()
    withdraw()
    withdrawalList()
}

package "Repositories" as R {
    class "OrderRepository" as R.OR {}
    class "UserRepository" as R.UR {}
    class "WithdrawalRepository" as R.WR {}
    R.UR -[hidden]-> R.WR
    R.WR -[hidden]-> R.OR
}

BL.S -l-o UC
BL.E o-- BL.S
BL.S.C.ORI <|. R.OR 
BL.S.C.URI <|. R.UR 
BL.S.C.WRI <|. R.WR 

@enduml
```

## Status mapping on accrual

NEW — заказ загружен в систему, но не попал в обработку;
PROCESSING — вознаграждение за заказ рассчитывается; 200(REGISTERED|PROCESSING)
INVALID — система расчёта вознаграждений отказала в расчёте; 200(INVALID), 204
PROCESSED — данные по заказу проверены и информация о расчёте успешно получена. 200(PROCESSED) 

### Accrual statuses
REGISTERED — заказ зарегистрирован, но не начисление не рассчитано;
INVALID — заказ не принят к расчёту, и вознаграждение не будет начислено;
PROCESSING — расчёт начисления в процессе;
PROCESSED — расчёт начисления окончен;