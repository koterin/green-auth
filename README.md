# green-auth Service

Simple authorization service serving as a firewall for all your online content

> Extra: there are some features described in another section (TBD)

## Build

> For building you'll need `ansible` and `ansible-playbook` installed on your machine

To build `green-auth` you need to:

1. Download the Repo

```
git@github.com:koterin/green-auth.git
```

2. Enter the new folder `green-auth`
3. Add the IP-address of the machine which will host `green-auth` to the `/green-auth.yaml` ansible script in the hist section. `green-host` is written as an example
4. Add to the `/configs/nginx` folder on the current machine SSL-sertificates for your site - `.crt` and `.key` files
5. Generate SSH-keys for the backend server:
Generate them with the command
```
ssh-keygen -f ~/[PATH]/green-auth/server/ssh_green-auth_key -t rsa
```
where `[PATH]` - path to the `green-auth` folder on your machine.

6. Add the public SSH key to the Deploy keys of the Repo with the passwords

7. Change the names of your added `.crt` and `.key` files in the `/green-auth.yaml` script - find the comments pointing to those places.

8. Change the name of the site and stuff in your NGINX config:

TBD

9. Create your TG-bot and add its key to the `/docker-compose.yaml` file (follow the comments)

10. Start the script
```
./green-auth.yaml
```
It will install all neccessary dependencies on your host machine and start the `green-auth`-service via Docker.

## Main info

Green-auth Service consists of:

- HTTP-server (GO)
- Client part (static HTML-pages + .js-scripts and .css-styles)
- Proxy (NGINX)
- PostgreSQL Database
- Funny captions

Server starts via Docker container on the `8080` port. All static pages are returned from the backend, even the authroization-free ones (for example, log in pages).

## Roles

Service has 3 roled described:

  - admin
    - Able to perform any action on the site
  - dev
    - Able to use Swagger
    - Able to see the account list
    - Unable to generate new accounts
  - service
    - Able to see the account list
    - Able to generate new accounts
    - Unable to see Swagger

## If you need to change something

1. If you need to set unique credentials for the user who will be the author of the commits made by the `green-auth`, you'll jave to adjust those line in the `/server/Dockerfile`:

```
RUN git config --global user.email "green.auth@sample.com"
RUN git config --global user.name "green-auth"
 ```

2. If you need to add new user you'll have to make an `INSERT` into the database:
   1. Enter the target Docker container with your running database (`green-db`):
    ```
    docker exec -it green-db psql -d green-db -U green-manager
    ```
    2. Add the user:
    ```
    INSERT INTO users VALUES (DEFAULT, '[USERNAME]@[SAMPLE].com', '[CHAT_ID]', '[ROLE]', DEFAULT, DEFAULT)
    ```
    where `[USERNAME]` - the user login, `[CHAT_ID]` - received Telegram chat_id of user and the bot, `[ROLE]` - user's assigned role (`admin`, `dev` or `service`).

    > ATTENTION! Please, assign the `admin` role carefully. A user with admin right may take some unresetable actions.

3. You may wanna set you users list with the database initialization, without the manual input. In that case add to your `/db/init.sql` file a line like that:
```
INSERT INTO users (id, email, chat_id, role, created_at, updated_at)
VALUES (DEFAULT, '[USERNAME]@[SAMPLE].com', '%CHAT_ID%', '%ROLE%', DEFAULT, DEFAULT);
```
where the meaning of the fields is the same as explained in the 2.2 step.
