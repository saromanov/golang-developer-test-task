<p align="center">
  <img src="https://hsto.org/webt/ih/ds/fu/ihdsfuqni5apj0my18tnukzztw0.png" alt="Logo" width="128" />
</p>

# Тестовое задание для GoLang-разработчика

![Project language][badge_language]
[![Build Status][badge_build]][link_build]
[![Do something awesome][badge_use_template]][use_this_repo_template]

# Описание структуры проекта
[cmd](https://github.com/saromanov/golang-developer-test-task/tree/master/cmd)
Вход в приложение

[assets](https://github.com/saromanov/golang-developer-test-task/tree/master/assets)
Файлы нужные для запуска проекта

[pkg/cmd](https://github.com/saromanov/golang-developer-test-task/tree/master/pkg/cmd)
Стартовая точка для запуска приложения

[pkg/config](https://github.com/saromanov/golang-developer-test-task/tree/master/pkg/config)
Модуль для конфигурации

[pkg/data](https://github.com/saromanov/golang-developer-test-task/tree/master/pkg/data)
Модуль для загрузки данных о парковках

[pkg/logger](https://github.com/saromanov/golang-developer-test-task/tree/master/pkg/logger)
Конфигурирование логгера

[pkg/models](https://github.com/saromanov/golang-developer-test-task/tree/master/pkg/models)
Модельки для проекта

[pkg/server](https://github.com/saromanov/golang-developer-test-task/tree/master/pkg/server)
Сервер приложения

[pkg/storage](https://github.com/saromanov/golang-developer-test-task/tree/master/pkg/storage)
База данных

# Запуск

```
docker-compose up --build
```

# Использование

Для поиска по парковкам нужно использовать конечную точку `/v1/search`.
Например запрос
```
http://localhost:3000/v1/search?id=122
```

Получает информацию о парковке с id=122 

```
[ 
   { 
      "global_id":1704751,
      "id":122,
      "system_object_id":"122",
      "name":"Парковка такси по адресу Беговая улица, дом 22, строение 13",
      "mode":"круглосуточно"
   }
]

```

Поддреживаемые параметры для запроса:
`?global_id=1` - поиск по Global ID
`?mode=круглосуточно` - поиск по Mode


Для получения метрик из Prometheus, нужно использовать конечную точку
```
http://localhost:3000/v1/metrics
```

```
...
http_requests_total{code="200",method="GET"} 6
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 0.52
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1.048576e+06
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 8

# HELP total_reads 
# TYPE total_reads counter
total_reads 0
# HELP total_requests 
# TYPE total_requests counter
total_requests 6
...


