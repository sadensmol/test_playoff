# Задание
### Введение

Наша компания разрабатывает новую онлайн игру. Перед полноценным запуском мы планируем дать доступ ограниченным группам пользователей для получения первой обратной связи. Для этого мы хотим создавать приглашения, по каждому из которых сможет зарегистрироваться определенное ограниченное число новых пользователей. Требуется сделать API сервис для приглашения новых пользователей в игру с закрытой/ограниченной регистрацией.

### Сценарий

1. Сотрудник нашей компании создает 3 приглашения с кодами: “twitter-reg1”, “telegram-test” и “instagram-hello”.
2. В трех соц сетях публикуются ссылки для регистрации с кодами приглашений. Пример: ”https://our-super-game.io/invitations/twitter-reg1”.
3. Пользователи проходят по ссылкам и оставляют адреса своих электронных почт (запрос уходит на **endpoint который требуется сделать**).
4. Далее на адреса электронных почт отправляются письма с дальнейшими инструкциями.

### Functional Requirements

- Приглашение может быть использовано максимум N раз (где N находится в пределах от 1 до 1000 включительно).
- Один пользователь может использовать только одно приглашение и только один раз.
- Для использования приглашения и последующей регистрации пользователь должен предоставить свой Email.

### ❗️Non-functional Requirements

- Приглашение может быть опубликовано в крупном аккаунте одной из популярных социальных сетей, что может привести с существенному количеству пользователей желающих зарегистрироваться по одному и тому же приглашению приблизительно в одно и то же время (сразу после публикации).
- Сервис должен быть масштабируем (горизонтально).
- Возможность zero-downtime деплоймент (сам pipeline создавать не надо, но возможность должна быть).

### Примечания

Для того, чтобы не выходить за разумные рамки тестового задания, API должен содержать только один endpoint — для сохранения желающих зарегистрироваться по приглашению.

Пример:

```jsx
POST /invitations/{code}
Request: 
{
 	“Email”: “[test@example.com](mailto:test@example.com)” 
}
Response: 
 Http 200 / 400 / и тд
 No body.
```

Это лишь пример, вы можете изменить контракт на ваше усмотрение.

Данные должны храниться в БД **MongoDB**. Структуры данных на ваш выбор.

Язык программирования — **Go**.

Как сами приглашения попадают в БД не имеет значения для этого тестового задания. Исходим из того, что у нас изначально уже есть какое-то количество приглашений.

Тесты обязательны. Не требуется большое покрытие, но для примера необходимо что-то показать.

Так как в рамках тестового задания практически невозможно довести проект до состояния production ready, пожалуйста, напишите сопроводительное письмо с описанием всего того, что бы вы добавили (как если бы это было не тестовое задание, а ваш рабочий проект). Какие-то вещи можно комментировать в коде с тегом NOTE (пример: ***// NOTE:** добавить больше тестов*).

Пишите любые советы, рекомендации и мысли (пример: я бы рекомендовал заменить MongoDB на XYZ для того, чтобы…).

Напишите ваши предложения и идеи по выбору инфраструктуры для этого приложения.

Напишите почему выбираете тот или иной фреймвокр или библиотеку.

Добавьте пошагавшую инструкцию как запустить проект на локальной машине.

# Проблемы и решения

### Обработка ошибок HTTP REST API

Для простоты решения обработка ошибок отсутствует - если сервер доступен то всегда возвращается только 200 статус. В случае любой ошибки возвращается 500 статус.

С точки зрения юзера, если он отправил запрос на сервер, то он получит ответ в любом случае но позже в email. Сайт покажет что все ок (так как получит 200 статус) и юзер не увидит ошибок, иначе он может подумать что его запрос не дошел до сервера и будет отправлять раз за разом.

### "Человеческий" сервис

Обработка запроса на сервере сделана асинхронной. Запросы ставятся в очередь в порядке их поступления и обрабатываются в порядке очереди. 

Асинхронная обработка позволит более человечно (справедливо) подходить к обработке запросов - в порядке поступления (очереди ) и не перегружать сервер. Если юзер ввел некоректный email, то мы пропустим такого юзера и "ход"  перейдет к следующему в очереди.


### Лимиты или подход "идите в жопу"

В случае если лимиты на приглашения закончились, то сервер все равно принимает запросы и позже все равно отсылает емейлы с извинениями и "благими пожеланиями". Это сделано для того чтобы не отбивать юзеров от продукта. В дальнейшем эта инфа может быть использоваться для маркетинга.

### Ограничение на количество запросов

У каждого сервиса есть лимит и время необходимое для увеличения пропускной способности (скейлинга). Для укрощения хайпа - на инфраструктуре должен стоять общий rate limiter для всех запросов - который при привышении лимита сервиса будет задерживать запросы (не откидывать).

### Проблемы параллельности и мультисервисности и горизонтального масштабирования

Сервисы паралелятся через LB и кластер настраивается в соотвествии с ожиданиями от промоушена. 
Блокировки для базы не используются (так как обработка асинхронная) - настраивается несколько mongodb в кластере (шардирование) со своим LB. В данном случае явный перекос в сторону записи.


# Выбор версий фреймворков и библиотек

go 1.21 - что локально сейчас стоит.


### HTTP 

стандартный пакет net/http - для простоты и надежности (также довольно таки шустрый). Обычно в проектах используется то с чем есть опыт работы и что уже используется в компании.

### MongoDB

'go.mongodb.org/mongo-driver' - у меня есть опыт с этой либой

### Logs

стандарный пакет slog - для простоты и надежности.
	
