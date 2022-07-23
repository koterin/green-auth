# Password Manager

Сервис для генерации паролей для различных внутренних сервисов Бери Заряд

## Сборка

> Для запуска скрипта потребуется установленный `ansible`

Чтобы собрать сервис, требуется:
1. Скачать репозиторий

```
git clone git@gitlab.com:berizaryad_project/devops/passwordmanager.git
```

2. Зайти в созданную папку `passwordManager`
3. Добавить IP-адрес нужной машины в свой iventory-файл (или `ansible-hosts`) под именем `password-host`
4. Добавить в папку `/configs/nginx` SSL-сертификаты для NGINX - .crt и .key
5. Сгенерировать ssh-ключи для контейнера бэкенда - и положить внутрь репозитория в /server:
Генерация - с помощью команды
```
ssh-keygen -f ~/[PATH>]/password-manager/server/ssh_password-manager_key -t rsa
```
где `[PATH]` - путь до репозитория менеджера.

6. Добавить публичный ключ (сгенерированный `~/password-mananager/ssh_password-manager_key.pub` в Deploy Keys [репозитория с паролями для сервисной карты.](https://gitlab.com/berizaryad_project/devops/servicepswd/-/tree/main))

7. Запустить скрипт
```
./password-manager.yaml
```

## Общая информация

Сервис состоит из:

- HTTP-сервера на GO
- Клиентской части - статичных HTML-страниц (+ .JS-скрипты и .css-стили)
- Прокси - NGINX
- Веселых текстов

Сервер запускается в Docker-контейнере на внутреннем порте `8080`, порт контейнера прокинут на порт машины `8080`; проксируется с помощью NGINX.
Все страницы возвращаются с бэкенда, а не NGINX (даже доступные без авторизации).

## Роли

В сервисе определены 3 роли:
  - admin
    - Может совершать любые действия на сайте
  - dev
    - Может пользоваться Swagger
    - Может просматривать список учетных записей сервисной карты (в зашифрованном виде)
    - Не может генерировать пароли для сервисной карты
  - service
    - Может просматривать список учетных записей сервисной карты (в зашифрованном виде)
    - Может генерировать пароли для сервисной карты
    - Не может просматривать Swagger

## Если надо что-то поменять

1. Если требуется задать уникальные креды для пользователя, от лица которого бэкенд будет совершать коммиты в репозиторий с паролями - следует править эти строки в `/server/Dockerfile`:
```
RUN git config --global user.email "password@berizaryad.ru"
RUN git config --global user.name "password-manager"
 ```

2. Если надо добавить нового пользователя - следует сделать INSERT в базу данных сервиса:
    1. Зайти в нужный контейнер с базой данных (password-db):
    ```
    docker exec -it password-db psql -d password-db -U password-manager
    ```
    2. Добавить пользователя:
    ```
    INSERT INTO users VALUES (DEFAULT, '[USERNAME]@berizaryad.ru', '[CHAT_ID]', '[ROLE]', DEFAULT, DEFAULT)
    ```
    где `[USERNAME]` - логин пользователя (или часть реальной почты), `[CHAT_ID]` - полученный chat_id пользователя в Телеграм-боте, `[ROLE]` - роль пользователя (admin или user).
    > ВАЖНО! Следует осторожно выдавать роль админа. Такой пользователь сможет свободно управлять паролями учетных записей сервисной карты.

3. Можно добавить авторизованного пользователя сразу при создании базы данных. В таком случае в файл `/db/init.sql` следует добавить строку:
```
INSERT INTO users (id, email, chat_id, role, created_at, updated_at)
VALUES (DEFAULT, '%USERNAME%@berizaryad.ru', '%CHAT_ID%', '%ROLE%', DEFAULT, DEFAULT);
```
где `%CHAT_ID%` - полученный chat_id пользователя в Телеграм-боте, `%ROLE%` - роль пользователя (admin или user).

## Генерация паролей

Those are the extras mentioned before: it may serve as a login-pass generator.

> You may delete this functionality without affecting the authorization part, but you will have to do it on your own. Instead to use just the auth part you may leave the backend methods the way they are but add your own Frontend or necessary resources redirects.

So the passwords here are stored in the GitHub Repo in the `.service` file in the form of `login:MD5-hashed pass`. The job of the gen-pass button on the site is to add new ones and successfuly push them back to the remote Repo.

1. Show passwords

Makes `git pull` from the backend container and outputs the contents of the password-file

2. Generate password

Logins in those accounts may be unique or sequential: `green1`, `green2`, etc. When the `Generate password`-button is pressed, backend checks out the last sequential login and suggests its incremented value to the user.
For example, if the last `green` login is `green34`, it will suggest `green35`.

User may change that login, but:
   - Login must not be blank
   - Login must not consist spaces
   - Login must not be duplicated

Once the login is accepted, backend generates random 8-symbol password, hashes it with MD5, concatenates strings to form an account string, adds it to the file and pushes to GitHub. Once it's done, user sees the generated password in the modal window on the site; button in the middle with the symbol of two pages copies the password to the clipboard. That's it!

3. Generate passwords

This button lets you generate multiple accounts at a time (from 2 to 100).
Instead of asking for a login input it waits for the number of accounts to be created, default being 2.

Logins for the accounts created this way will be all sequential: `green25`, etc.

