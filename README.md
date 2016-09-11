# vindinium-go
A Vindinium bot using GoLang and the [A* Search Algorithm](https://en.wikipedia.org/wiki/A*_search_algorithm)

It takes two possible command line arguments - 

| Mode        | API Key           |
| ------------- |:-------------:|
| Training     | yourKey |
| Arena      | yourKey      |

```
./vindinium-go training yourKey e.g. awffitw4
```
If all is working well you should see output that looks like the following:

```
Starting GoBot with settings mode=training, key=1234

Game starting at url=http://vindinium.org/api/training, viewurl=http://vindinium.org/1234

Located a player mine to go after x: 3, y: 5
Moving in direction: North

Move #: 4 of 800
The closest mine is owned by me x: 3, y: 5, looking for an alternative
Located a player mine to go after x: 3, y: 10
Moving in direction: East
```
