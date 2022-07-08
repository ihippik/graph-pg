# GraphPG
### Simple query builder for go-pg ORM.

* parameters are separated by commas `,`
* name and value of the variable are separated by a colon `:`

Simple case

#### VARCHAR

`name:Ivan, surname:Petrov` we'll get `WHERE name = 'Ivan' AND surnam = 'Petrov'`

`name:~van` we'll get `WHERE name ILIKE '%van%'`

`name:Ivan|Ann` we'll get `WHERE (name = 'Ivan' OR name = 'Petrov')`


#### INT
`age_1|age_2:10` we'll get `WHERE (age_1 = 10 OR age_2 = 10)`

`age:10` we'll get `WHERE age = 10`

`age:||10` we'll get `WHERE age <= 10`

`age:10||` we'll get `WHERE age >= 10`

`age:10~12~15` we'll get `WHERE age IN (10,12,15)`

`age:10;;20` we'll get `WHERE age > 10 AND age <=20`
