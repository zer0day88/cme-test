### How to run this project on kubernetes (local using kind)
1. Deploy kind, cassandra, redis, and nginx-ingress
```shell
cd deploy-kind
terraform init
terraform apply
```

2. default user for redis and cassandra
```text
redis
password: 12345

cassandra:
user: cassandra
password: 12345
```

3. Connect to cassandra service using port forwarding to init db and table , use default user and password to login<br>
    create new keyscpaces named cme and required tables.
```cassandraql
CREATE KEYSPACE cme WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
describe keyspaces;

create table users
(
   id         uuid,
   username   varchar,
   password   varchar,
   created_at timestamp,
   primary key ( username,id )
);

create table messages
(
   id uuid,
   sender varchar,
   recipient varchar,
   content varchar,
   message_time timestamp,
   primary key (recipient, sender,id)
);         
```


4. Deploy user service, message service, and ingress
```shell
cd deploy-services
terraform init
terraform apply
```

5. service could be accessed from public ingress endpoint
```text
user service : http://<public ip or localhost>/user/v1/register
message service : http://<public ip or localhost>/message/v1/send
```

### User Service APIs:
1. POST `/v1/register` (to register user) <br>
   payload <br>
   ```json
   {
      "username": "johndoe",
      "password": "abcdef"
   }
   ```
   success response
   ```json
   {
      "code": 200,
      "message": "success"
   }
   ```
2. POST `/v1/login` (to login) <br>
   payload
   ```json
   {
      "username": "johndoe",
      "password": "abcdef"
   }
   ```
   success response
   ```json
   {
    "code": 200,
    "message": "Success",
    "data": {
        "AccessToken": {
            "Token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTY5NTA2MTMsImlhdCI6MTcxNjk0ODgxMywibmJmIjoxNzE2OTQ4ODEzLCJzdWIiOiI4ODk4MDlkZi1mMDNlLTQzOTMtODdmZS03NzE2YWJhYjcxZDUiLCJ0b2tlbl91dWlkIjoiMjIzNjRjNjUtMmJhNS00NDFlLTg4YzMtZmU3MmM1OTYxYWE5In0.LMAdfYLmjjfmyVFAneTScgMOBSyxqF4eRYWv3jzhuzG6Ge0A-nM2fTzz4mdZb0m11_fc-D4kX-9bBvyx_5HOx6naCu-mEjX7rOOeDyAc4Oriwr-Il5Jegkdrp-uY0-RXSvlLexb-3Mdkgr1aDCkcITZETN9M3tMpDAPv4n75zh7qNKiuDRp8Bw8Fuwf9MgimGI3JNKZLwqOtIUTmNA2ZRSmXr7hDDkT-SYhmW5mB4GJSkYTSno_T0wdQOekCF-It-zGczsKE5BlM91iIzB_fSgkbpH5HVlmNGW4ku-C2SZiDOvVo0jp-IcxYBs_MOEzk3OJ-bsZZjg-uUFj0gyTtnndAMfemlVKtzY4WzK1wg0GHSD6Td5DDmfjpZVhN0Ojs9HLPFoCmFxug_-FLwy49KjMCZAil-W_UQhQzhchlThKvg_S4WtU0DWbKiOgrvB8pMu8kf77PWBcH-1eq_IQ5B1PfBeLe3NnLqllvu49wPl6yvnRIncfNhpTiKwYHpklNCbPwbA4EriCSw4WxAgsDk4CpGUhfBHT960OVV0WlXxDi4ZsSJbzVO0hBpVZtYCzbD_vIz4ym_3TMXS2D6blP9dmKMYQ6sO6TV0f6FARkbJSMqPoW83lNCy-mAAAqcylxJhRoFaPh2I7L2u1XcNwQDr6bw-FKzagbbArRXhRIOaM",
            "TokenUuid": "22364c65-2ba5-441e-88c3-fe72c5961aa9",
            "UserID": "889809df-f03e-4393-87fe-7716abab71d5",
            "ExpiresIn": 1716950613
        },
        "RefreshToken": {
            "Token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTY5NTI0MTMsImlhdCI6MTcxNjk0ODgxMywibmJmIjoxNzE2OTQ4ODEzLCJzdWIiOiI4ODk4MDlkZi1mMDNlLTQzOTMtODdmZS03NzE2YWJhYjcxZDUiLCJ0b2tlbl91dWlkIjoiZTA3MWZiM2UtNTI2YS00ZjJlLTg2MTItMDQ0YjkxOTdiZjE3In0.JMiGk8ma1JmCPl1uVNxnmoESv6Hsxs8j8_i2vbwSWcIYTdwxJqdCIQ7UDiWj9NJXAY7T1yuXcbwdjrOZalyHckp_dbyEK9u9uw6gwrYdp6UBATNI5ylxXBOuUWiRIVbIrtAHTsjkZ-NvEtLF7wADX0_Ah0Ippw8iVM9rEeaJGNC7RyUSyl0KVmfW-l6B4Apl4yatk-JouGqDW7B_sxU9JIyNnHmQgN-5GrAmMLJQCrqp49x6oc7_2Z3-xXQoawpmcRE6DjyVRJjIgTzVyt2U03Le62X1EfrVhIzPLhNv2eniJWMc64NWY6rl0-zEaCcQY6lShn-XN7D3q4kdKPH4loyCuISxN9DINAc8U1sOnVc0-RQzS1QX20f-F--ZWf9nhUzwhOnYeN2eRNgeawipu3lT2xWrAgdbuSguLwh2zPo6oZYazQ33N3duPuin21L6PkQs2_h_zxtWn-F3gD5k0JXBmd5PepjiGB-NIxjFplmB1Zkqt7UfItj_lgZMtBWVN3e6GUMhfwe3JRYpip1BBFOyBQ8LKUVht9sXxBGtOxubfp_NB-mQbRodDHVRxkGj0uL8E6f2o-B7IPf5VkKFsbOLacWReOWWGTahbuHiefExc67uexWjN4vCvToAcPe6hhHuyC-Ndkp_OFY_CzFJYtSDntlKBL3t43eAtjn52fQ",
            "TokenUuid": "e071fb3e-526a-4f2e-8612-044b9197bf17",
            "UserID": "889809df-f03e-4393-87fe-7716abab71d5",
            "ExpiresIn": 1716952413
         }
      }
   }
   ```

### Message Service APIs:
1. POST `/v1/send`  <br>
   payload <br>
   ```json
   {
      "to": "johndoe",
      "content": "well done"
   }
   ```
   success response
   ```json
   {
      "code": 200,
      "message": "success"
   }
   ```
2. GET `/v1/messages`  <br>
  
   success response
   ```json
   {
      "code": 200,
      "message": "Success",
      "data": [
      {
         "id": "14cbe435-e1cc-4311-b566-09c7d289a052",
         "sender": "joko",
         "recipient": "paijo",
         "content": "well done",
         "message_time": "2024-07-29T19:22:44Z"
      }
      ]
   }
   ```
   
> no unit test, if I have more time to do that I will use this one https://testcontainers.com/ for end-to-end testing 

### Explanation
This project has 2 services user and messaging service and infra on the local kubernetes
I use redis only for login and refresh token, not in messaging services
in real case / production grade the tech would be very different
and the message could be cached using TTL on redis based on frequent request,
more frequent the message requested and then more longer TTL