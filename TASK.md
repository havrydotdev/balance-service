# Test task for the position of an intern-backender 

## Microservice for working with user balance. 

**Problem:** 

There are many different microservices in our company. Many of them somehow want to interact with the user's balance. The architectural committee decided to centralize the work with the user's balance into a separate service. 

**Task:** 

It is necessary to implement a microservice for working with user balances (crediting funds, debiting funds, transferring funds from user to user, as well as a method for obtaining a user's balance). The service must provide an HTTP API and accept/return requests/responses in JSON format. 

**Usage scenarios:** 

Here are some simplified cases close to reality.
1. The billing service with the help of external merchants (ala via visa/mastercard) has processed the transfer of money to our account. Now billing needs to add this money to the user's balance. 
2. The user wants to buy some service from us. To do this, we have a special service management service, which checks the balance before using the service and then writes off the required amount. 
3. In the near future, it is planned to give users the opportunity to transfer money to each other within our platform. We decided to foresee such an opportunity in advance and put it into the architecture of our service. 

**Code requirements:** 

1. Development language: Go. We are ready to consider solutions in PHP/Python/other languages, but golang is our priority.
2. You can use any frameworks and libraries 
3. Relational DBMS: MySQL or PostgreSQL 
4. All code must be posted on Github with a Readme file with instructions for launching and sample requests / responses (you can simply describe methods in Readme, you can use Postman, you can copy requests in Readme curl, you get the idea ...) 5. If there is 
a 
need, you can connect caches (Redis) and / or queues (RabbitMQ, Kafka) Readme file for the project should contain a list of issues that the candidate faced and how he solved them)
7. Development of the interface in the browser is NOT REQUIRED. Interaction with the API is assumed through requests from the code of another service. For testing, you can use any convenient tool. For example: in the terminal via curl or Postman. 

**Will be a plus:** 

1. Using docker and docker-compose to raise and deploy a dev environment. 
2. API methods return human-readable error descriptions and corresponding status codes when they occur. 
3. Everything is implemented on GO, yet we are interviewing a developer on GO. HINT: At the interview one way or another there will be questions about Go. Whoever read it, well done :) 
4. Unit / integration tests are written. 

**Main task (minimum):**

The method of accruing funds to the balance. Accepts user id and how much money to deposit. 

The method of debiting funds from the balance. Accepts user id and how much to write off. 

The method of transferring funds from user to user. It accepts the id of the user from which the funds should be debited, the id of the user to whom the funds should be credited, as well as the amount. 

Method for obtaining the user's current balance. Accepts a user id. The balance is always in rubles. 

**Task details:** 

1. The accrual and write-off methods can be combined into one, if the common architecture allows it. 
2. By default, the service does not contain any data on balances (an empty table in the database). Balance data appears when the money is first credited.
3. Data validation and error handling are left up to the candidate. 
4. The list of fields for methods is not fixed. Only the bare minimum is listed. As part of the implementation of additional tasks, additional fields are possible. 
5. No migration mechanism needed. It is enough to provide the final SQL file with the creation of all the necessary tables in the database. 
6. User balance - very important data in which errors are unacceptable (in fact, we work here with real money). It is necessary to always keep the balance up to date and avoid situations when the balance can go negative. 
7. The default balance currency is always rubles. 

**Additional tasks**

The following are the extras. tasks. They are not mandatory, but their implementation will give a significant advantage over other candidates. 

*Additional task 1:* 

Effective managers wanted to add goods and services to our applications in currencies other than the ruble. It is necessary to be able to display the user's balance in a currency other than the ruble. 

Task: add to the method of obtaining the balance additional. parameter. Example: ?currency=USD. 
If this parameter is present, then we must convert the user's balance from the ruble to the specified currency. Data on the current exchange rate can be taken from here https://exchangeratesapi.io/ or from any other open source.

Note: we remind you that the base currency that is stored on our balance sheet is always the ruble. As part of this task, the conversion always takes place from the base currency. 

*Additional task 2:* 

Users complain that they do not understand why the funds were debited (or credited). 

Task: it is necessary to provide a method for obtaining a list of transactions with comments from where and why the funds were credited/debited from the balance. It is necessary to provide for pagination and sorting by amount and date.
