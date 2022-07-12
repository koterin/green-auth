# green-auth Service

Simple authorization service serving as a firewall for all your online content

> Extra: there are some features described in another section (TBD)

## Build

> For building you'll need `ansible` and `ansible-playbook` installed on your machine

To build `green-auth` you need to:

1. Download the Repo

```
git clone git@github.com:koterin/green-auth.git
```

2. Add the IP-address of the machine which will host `green-auth` to the `/green-auth.yaml` ansible script in the host section. `green-host` is written as an example
4. Add to the `/configs/nginx` folder on the current machine SSL-sertificates for your site - `.crt` and `.key` files
5. Generate SSH-keys for the backend server:
Generate them with the command
```
ssh-keygen -f ~/[PATH]/green-auth/server/ssh_green-auth_key -t rsa
```
where `[PATH]` - path to the `green-auth` folder on your machine.

6. Add the public SSH key to the Deploy keys of the Repo with the passwords (GitLab or GitHub - doesn't matter)

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

Server starts via Docker container on the `8080` port. All static pages are returned from the backend, even the authorization-free ones (for example, log in pages).

## Roles

Service has 3 roles described:

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

1. If you need to set unique credentials for the user who will be the author of the commits made by the `green-auth`, you'll jave to adjust those lines in the `/server/Dockerfile`:

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
    INSERT INTO users VALUES (DEFAULT, '[USERNAME]@[SAMPLE].com', '[CHAT_ID]', '[ROLE]', DEFAULT, DEFAULT);
    ```
    where `[USERNAME]` - the user login, `[CHAT_ID]` - received Telegram chat_id of user and the bot, `[ROLE]` - user's assigned role (`admin`, `dev` or `service`).

    > ATTENTION! Please, assign the `admin` role carefully. A user with admin rights may take some unrecoverable actions.

3. You may wanna set your users list with the database initialization, without the manual input. In that case add to your `/db/init.sql` file a line like that:
```
INSERT INTO users (id, email, chat_id, role, created_at, updated_at)
VALUES (DEFAULT, '[USERNAME]@[SAMPLE].com', '%CHAT_ID%', '%ROLE%', DEFAULT, DEFAULT);
```
where the meaning of the fields is the same as explained in the 2.2 step.

## How the authorization works

#### Log In

When user is registered (see step 2 above) he may authorize via `green-auth`.
First he must provide his login; you may add exclusive pattern for matching those logins (for example, only emails for some particular domain) in the `/server/static/index.html` file. Current pattern set to 2-20 symbols.
Once the field is validated, the Proceed Button activates.

#### Validate login

Proceed button sends request to the backend; it checks if that login actually exists in the database. In the negative case the message `incorrect login` appears and Proceed Button is deactivated again (until some letters are changed in the input field). If the login is correct, user is being redirected to the OTP-page.

#### One Time Password

One Time Password (OTP) is being generated as a 7-digit sequence. Backend sends it via Telegram API to the Telegram Chat between that user and the Bot.

As it's sent with the code style formatting, user may just click on the code (on the PC) to copy it to the clipboard.

Next the code must be inserted in the code input field on the `green-auth` site. If the code is incorrect, user will see the message `wrong code`. If he inserts wrong codes 5 times in a row, he will be asked to get the new code (attempts are being checked for security reasons). 

So the rules for the OTP are:
   - No more than 5 attempts to insert wrong code (without getting the new one)
   - No more than 5 attempts to get the code (5 messages) in 5 minutes
   - No more than 1 attempt to get the code (1 message) in 30 seconds

If the code is right, user will be redirected to the Homepage.

#### Generate passwords

Those are the extras mentioned before: it may serve as a login-pass generator.

> You may delete this functionality without affecting the authorization part, but you will have to do it on your own. Instead to use just the auth part you may leave the backend methods the way they are but add your own Frontend or necessary resources redirects.

So the passwords here are stored in the GitHub Repo in the `.service` file in the form of `login:MD5-hashed pass`. The task of gen-pass button on the site is to add new ones and successfuly push them back to the remote Repo.

1. Show passwords
Makes `git pull` from the backend container and outputs the contents of the password-file

2. Generate password
Logins in those accounts may be unique or sequential: `green1`, `green2`, etc. When the `Generate password`-button is pressed, backend checks out the last sequential login and suggests its incremented value to the user.
For example, if the last `green` login is `green34`, it will suggest `green35`.

User may change that login, but:
   - Login must not be blank
   - Login must not consist spaces
   - Login must not be duplicated

Once the login is accepted, backend generates random 8-symbol password, hashes it with MD5, contatenates strings to form an account string, adds it to the file and pushes to GitHub. Once it's done, user sees the generated password in the modal window on the site; button in the middle with the symbol of two pages copies the password to the clipboard. That's it!

3. Generate passwords
This button lets you generate multiple accounts at a time (from 2 to 100).
Instead of asking for a login input it waits for the number of accounts to be created, default being 2.

Logins for the accounts created this way will be all sequential: `green25`, etc.

#### Swagger

Another feature of the site - swagger firewalling.

Swagger's default use mode - static pages - don't provide any authorization feature (to see and use your methods). So `green-auth` solves it.

Put your swagger-definition file in the `/server/static/swagger/swagger.yaml` and build the server.

> Keep in mind that all frontend pages and Swagger included lie in the server Docker container, so you'll have to rebuild the container with every change or edit those files directly in the container.
> Little tip: to edit files in the running server Docker container. enter it like that:
`docker exec -it green-server sh`
(`sh` - because server uses the bullseye Linux image which does not have bash)
Next download some text editor, like this:
`apt-get install vim`
and happily edit your files inside. Remember, if you rebuild the container you'll have to install vim again.
