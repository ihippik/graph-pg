# GraphPG
### Simple query builder for go-pg ORM.

* parameters are separated by commas `,`
* name and value of the variable are separated by a colon `:`

Simple case

#### VARCHAR

`name:Ivan, surname:Petrov` we get `WHERE name = 'Ivan' AND surnam = 'Petrov'`

`name:~van` we get `WHERE name ILIKE '%van%'`

`name:Ivan|Ann` we get `WHERE (name = 'Ivan' OR name = 'Petrov')`


#### INT
`age:10` we get `WHERE age = 10`

`age:||10` we get `WHERE age <= 10`

`age:10||` we get `WHERE age >= 10`

`age:10~12~15` we get `WHERE age IN (10,12,15)`

`age:10;;20` we get `WHERE age > 10 AND age <=20`
