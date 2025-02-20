# Тестовое задание на позицию стажёра-бэкендера

## Микросервис для работы с балансом пользователей

**Проблема:**

В нашей компании есть много различных микросервисов. Многие из них так или иначе хотят взаимодействовать с балансом пользователя. На архитектурном комитете приняли решение централизовать работу с балансом пользователя в отдельный сервис. 

**Задача:**

Необходимо реализовать микросервис для работы с балансом пользователей (зачисление средств, списание средств, перевод средств от пользователя к пользователю, а также метод получения баланса пользователя). Сервис должен предоставлять HTTP API и принимать/отдавать запросы/ответы в формате JSON. 

**Сценарии использования:**

Далее описаны несколько упрощенных кейсов приближенных к реальности.
1. Сервис биллинга с помощью внешних мерчантов (аля через visa/mastercard) обработал зачисление денег на наш счет. Теперь биллингу нужно добавить эти деньги на баланс пользователя. 
2. Пользователь хочет купить у нас какую-то услугу. Для этого у нас есть специальный сервис управления услугами, который перед применением услуги резервирует деньги на отдельном счете и потом списывает в доход компании. 
3. Бухгалтерия раз в месяц хочет получить сводный отчет по всем пользователям в разрезе каждой услуги.


**Требования к сервису:**

1. Сервис должен предоставлять HTTP API с форматом JSON как при отправке запроса, так и при получении результата.
2. Язык разработки: Golang.
2. Фреймворки и библиотеки можно использовать любые.
3. Реляционная СУБД: MySQL или PostgreSQL.
4. Использование docker и docker-compose для поднятия и развертывания dev-среды.
4. Весь код должен быть выложен на Github с Readme файлом с инструкцией по запуску и примерами запросов/ответов (можно просто описать в Readme методы, можно через Postman, можно в Readme curl запросы скопировать, и так далее).
5. Если есть потребность в асинхронных сценариях, то использование любых систем очередей - допускается.
6. При возникновении вопросов по ТЗ оставляем принятие решения за кандидатом (в таком случае Readme файле к проекту должен быть указан список вопросов с которыми кандидат столкнулся и каким образом он их решил).
7. Разработка интерфейса в браузере НЕ ТРЕБУЕТСЯ. Взаимодействие с API предполагается посредством запросов из кода другого сервиса. Для тестирования можно использовать любой удобный инструмент. Например: в терминале через curl или Postman.

**Будет плюсом:**

1. Покрытие кода тестами.
2. [Swagger](https://swagger.io/solutions/api-design/) файл для вашего API.
3. Реализовать сценарий разрезервирования денег, если услугу применить не удалось.

**Основное задание (минимум):**

- Метод начисления средств на баланс. Принимает id пользователя и сколько средств зачислить.
- Метод резервирования средств с основного баланса на отдельном счете. Принимает id пользователя, ИД услуги, ИД заказа, стоимость.
- Метод признания выручки – списывает из резерва деньги, добавляет данные в отчет для бухгалтерии. Принимает id пользователя, ИД услуги, ИД заказа, сумму.
- Метод получения баланса пользователя. Принимает id пользователя.

**Детали по заданию:**

1. По умолчанию сервис не содержит в себе никаких данных о балансах (пустая табличка в БД). Данные о балансе появляются при первом зачислении денег. 
2. Валидацию данных и обработку ошибок оставляем на усмотрение кандидата. 
3. Список полей к методам не фиксированный. Перечислен лишь необходимый минимум. В рамках выполнения доп. заданий возможны дополнительные поля.
4. Механизм миграции не нужен. Достаточно предоставить конечный SQL файл с созданием всех необходимых таблиц в БД. 
5. Баланс пользователя - очень важные данные в которых недопустимы ошибки (фактически мы работаем тут с реальными деньгами). Необходимо всегда держать баланс в актуальном состоянии и не допускать ситуаций когда баланс может уйти в минус. 
6. Мультивалютность реализовывать не требуется.

**Дополнительные задания**

Далее перечислены дополнительные задания. Они не являются обязательными, но их выполнение даст существенный плюс перед другими кандидатами. 

*Доп. задание 1:*

Бухгалтерия раз в месяц просит предоставить сводный отчет по всем пользователем, с указанием сумм выручки по каждой из предоставленной услуги для расчета и уплаты налогов.

Задача: реализовать метод для получения месячного отчета. На вход: год-месяц. На выходе ссылка на CSV файл.

Пример отчета:

название услуги 1;общая сумма выручки за отчетный период

название услуги 2;общая сумма выручки за отчетный период

*Доп. задание 2:*

Пользователи жалуются, что не понимают за что были списаны (или зачислены) средства. 

Задача: необходимо предоставить метод получения списка транзакций с комментариями откуда и зачем были начислены/списаны средства с баланса. Необходимо предусмотреть пагинацию и сортировку по сумме и дате. 
