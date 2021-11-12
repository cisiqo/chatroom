#### 1、认证登录

###### 如下字段全部是必须项

`{"Type":"join", "Username":"abc", "Account":"cccc", "Password":"XXXXXXX"}`

*注意：认证并没有实现



###### 成功会返回

`{"Message":"授权认证成功","Status":"ok","Time":"2021-11-12 10:14:25","Type":"login","Username":"abc"}`



###### 失败会返回

`{"Message":"授权认证失败","Status":"error","Time":"2021-11-12 10:14:25","Type":"login","Username":"abc"}`



#### 2、创建房间并加入房间

`{"Type":"join"}`



成功会返回一条广播消息
`{"Message":"abc进入直播房间","RoomID":"38297699","Time":"2021-11-12 10:17:54","Type":"message"}`



#### 3、加入房间

`{"Type":"join", "RoomID":"38297699"}`



###### 如果房间存在，会进入直播间，并广播一条消息

`{"Message":"aaaa进入直播房间","RoomID":"38297699","Time":"2021-11-12 10:20:40","Type":"message"}`



###### 如果房间不存在，会返回

`{"Message":"直播房间ID不存在","RoomID":"3829769911111","Status":"error","Time":"2021-11-12 10:22:48","Type":"join"}`



###### 如果房间为空，会返回

`{"Message":"直播房间ID为空","RoomID":"","Status":"error","Time":"2021-11-12 10:23:53","Type":"join"}`



#### 4、退出房间

`{"Type":"leave", "RoomID":"38297699"}`

###### 如果房间存在，会进入直播间，并广播一条消息

`{"Message":"aaaa退出直播房间","RoomID":"38297699","Time":"2021-11-12 10:20:40","Type":"message"}`



###### 如果房间不存在，会返回

`{"Message":"直播房间ID不存在","RoomID":"3829769911111","Status":"error","Time":"2021-11-12 10:22:48","Type":"leave"}`



###### 如果房间为空，会返回

`{"Message":"直播房间ID为空","RoomID":"","Status":"error","Time":"2021-11-12 10:23:53","Type":"leave"}`



#### 5、发消息

###### 消息的必须项，可以随便自行添加其他字段

`{"Type":"message"}`    

该消息会广播



###### 未登录，会返回

`{"Message":"您未授权认证，请登录授权","Status":"error","Time":"2021-11-12 10:28:02","Type":"message"}`



###### 未加入直播房间，会返回

`{"Message":"您未加入直播房间","Status":"error","Time":"2021-11-12 10:29:13","Type":"message"}`



#### 6、API接口

1. `/api/create_room`

2. `/api/get_rooms_number`

3. `/api/get_rooms`

4. `/api/get_room_members_number`

   

以上接口没有参数，返回都是json字符串





