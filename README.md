# Overview
Applikasi ini terdiri dari 2 endpoint rest API, yaitu register user dan login, aplikias ini menggunakan data base postgres

## Cara Menjalankan
1. create file .env dan masukan value berikut
    HTTP_PORT=9010
    DB_CONN=sslmode=disable host=localhost port=5431 user=postgres dbname=edufund password=Standar123.
2. Cara Menjalakan manual aplikias
    a. sesuaikan file .env dari segi port dan connection DB
    b. lalu jalankan go run main.go
3. Menjalakan aplikasi di atas docker
    a. seuaikan connection dan port pada file docker-compose.yml
    b. jalakan command docker-compose up -d

## Cara menggunakan endpoint Rest menggunakan postman
1. enpoint register
    a. url (<url>/register)
    b. body, jenis raw JSON
        {
            "full_name": "test",
            "user_name": "test@gmail.com",
            "password": "abcdefghijklmn",
            "confirmation_password": "abcdefghijklmn"
        }
2. enpoint login
    a. url (<url>/login)
    b. body, jenis raw JSON
        {
            "user_name": "test@gmail.com",
            "password": "abcdefghijklmn"
        }