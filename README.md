# Для чего создан итогоывый проект

Прежде всего для проверки накопленных знаний в новом для себя языке программирования.
А вообще проект представляет из себя веб-сервер, который реализует функциональность планировщика задач, которые хранятся в базе данных SQLite.

# Список выполненных заданий со звездочкой

Выполнено большинство заданий со звездолчкой, за исключением реализации правил повторения задач по дням недели и месяцам.

# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

Директория `app` содержит директории с пакетами проекта.

Директория `app\database` содержит функции для работы с базой данных, включая ее инициализацию и работы с таблицей `scheduler`.

Директория `app\handler` содержит функцию для запуска файл-сервера, чтобы иметь доступ к файлам фронтенда.

Директория `app\middleware` содержит реализацию функции аутентификации.

Директория `app\server` содержит функции для запуска HTTP-сервера.

Директория `app\service` содержит реализацию служебной функции NextDate для вычисления правила повторения задач.

Директория `app\taskscontrol` содержит обработчики HTTP-запросов, включая служебные.

# Сборка образа проекта в Docker:

`docker build -t go_final_project .`

# Для запуска контейнера Docker необходимо открыть терминал в текущей директории проекта и выполнить две команды: 
1) Первая создает пустой файл БД, чтобы его ретранслировать в Docker.
2) Вторая создает и запускает контейнер на основе образа проекта.

Windows:
`echo $null >> scheduler.db`\
`docker run -p 7666:7666 -v ${PWD}/scheduler.db:/app_bin/scheduler.db go_final_project:latest`

Linux:
`touch scheduler.db`\
`docker run -p 7666:7666 -v $(pwd)/scheduler.db:/app_bin/scheduler.db go_final_project:latest`
