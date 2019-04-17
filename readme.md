# The Gathering 2019: fast AI

## Game Rules

 * You can move up, down, left or right each game tick.
 * You only move one square per game tick.
 * If you don't send a new command until the new game tick the game picks a random direction for you.
 * Each game tick the latest command sent by all players is handled at the same time. 
 * You win by having the most points at the end of the game.
 * The round is over when there are either no pellets left, or only one player left.
 * The round is over when tick 1000 is reached.
 * The round is over when there is less than two snakes left on the map.
 * The snake dies if it collides with an wall.
 * The snake dies if it collides with the side of an other snake.
 * The snake dies if it tries to go back onto itself.
 * Each food is worth 1 point.
 * When the forth player dies the three left are granted 1 point each.
 * When the third player dies the two left are granted 1 point each.
 * The last snake standing is granted 3 points.
 * The map size is decided based on the number of players at the map at once.
 * The map will have other walls. (Without walls is not implemented yet!)
 * If any bugs that do not change the game outside the rules is found they will be patched during the week (Just poke us!
 ).
 
## Binaries

You can find the precompiled binaries for the  application under releases
https://github.com/N0K0/aicompo-tg19/releases

## Communication Protocol
This game is based on communication via JSON objects over websockets
An example bot is found in the *simple_player.js* can be used as an example. 
We will try to push more examples during the week.

The user should first connect to the communication endpoint for users. 
Which is located at [ws://localhost:8080/ws](). 
Any attempt to connect to any of the other endpoint will end in an disqualification

When a player is connected successfully it will receive the following message:
```JSON
{
  "type":"info",
  "message":"Hi!"
}
```

### For setting a username:
```JSON
{
  "type":"username",
  "value":"yLIDs"
}
```

If good:
```JSON
{
  "type":"info",
  "message":"Username OK!"
}
```

### For setting a color:
The ```value``` field can be any color that is accepted by the [fillStyle property]
```JSON
{
  "type":"color",
  "value":"rgb(149, 160, 10)"
}
```

If good:
```JSON
{
  "type":"info",
  "message":"Color OK!"
}
```

When the game is started the player will start to receive objects that looks like the following:
### Understanding the map during the game

```JSON
{
  "NumPlayers": 2,
  "Players": {
    "0M3xh": {
      "username": "0M3xh",
      "Color": "rgb(146, 166, 9)",
      "PosX": [
        7,
        6,
        6
      ],
      "PosY": [
        10,
        10,
        10
      ],
      "Head": {
        "X": 7,
        "Y": 10
      },
      "Tail": {
        "X": 6,
        "Y": 10
      },
      "Size": 3,
      "TotalScore": 0,
      "RoundScore": 0
    },
    "2Qdau": {
      "username": "2Qdau",
      "Color": "rgb(248, 200, 16)",
      "PosX": [
        4,
        4,
        4
      ],
      "PosY": [
        2,
        2,
        2
      ],
      "Head": {
        "X": 4,
        "Y": 2
      },
      "Tail": {
        "X": 4,
        "Y": 2
      },
      "Size": 3,
      "TotalScore": 0,
      "RoundScore": 0
    }
  },
  "GameStatus": {
    "Status": "running",
    "RoundNumber": 1,
    "TotalRounds": 10,
    "GameMap": {
      "SizeX": 14,
      "SizeY": 14,
      "Content": [
        ["'X'","'X'","'X'","'X'","'X'","'X'","'X'","'X'","'X'","'X'","'X'","'X'","'X'","'X'"],
        ["'X'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'X'"],
        ["'X'","'_'","'_'","'^'","'*'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'X'"],
        ["'X'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'X'"],
        ["'X'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'X'"],
        ["'X'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'X'"],
        ["'X'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'X'"],
        ["'X'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'X'"],
        ["'X'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'X'"],
        ["'X'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'X'"],
        ["'X'","'_'","'_'","'_'","'_'","'_'","'*'","'*'","'_'","'_'","'_'","'_'","'^'","'X'"],
        ["'X'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'X'"],
        ["'X'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'_'","'X'"],
        ["'X'","'X'","'X'","'X'","'X'","'X'","'X'","'X'","'X'","'X'","'X'","'X'","'X'","'X'"]
      ],
      "Heads": null,
      "Walls": [
      {"X":0,"Y":0},{"X":1,"Y":0},{"X":2,"Y":0},{"X":3,"Y":0},{"X":4,"Y":0},{"X":5,"Y":0},{"X":6,"Y":0},{"X":7,"Y":0},{"X":8,"Y":0},{"X":9,"Y":0},{"X":10,"Y":0},{"X":11,"Y":0},{"X":12,"Y":0},
      {"X":13,"Y":0},{"X":0,"Y":13},{"X":1,"Y":13},{"X":2,"Y":13},{"X":3,"Y":13},{"X":4,"Y":13},{"X":5,"Y":13},{"X":6,"Y":13},{"X":7,"Y":13},{"X":8,"Y":13},{"X":9,"Y":13},{"X":10,"Y":13},
      {"X":11,"Y":13},{"X":12,"Y":13},{"X":13,"Y":13},{"X":0,"Y":0},{"X":0,"Y":1},{"X":0,"Y":2},{"X":0,"Y":3},{"X":0,"Y":4},{"X":0,"Y":5},{"X":0,"Y":6},{"X":0,"Y":7},{"X":0,"Y":8},{"X":0,"Y":9},
      {"X":0,"Y":10},{"X":0,"Y":11},{"X":0,"Y":12},{"X":0,"Y":13},{"X":13,"Y":0},{"X":13,"Y":1},{"X":13,"Y":2},{"X":13,"Y":3},{"X":13,"Y":4},{"X":13,"Y":5},{"X":13,"Y":6},{"X":13,"Y":7},
      {"X":13,"Y":8},{"X":13,"Y":9},{"X":13,"Y":10},{"X":13,"Y":11},{"X":13,"Y":12},{"X":13,"Y":13}
      ],
      "Foods": [
        {
          "X": 12,
          "Y": 10
        },
        {
          "X": 3,
          "Y": 2
        }
      ]
    },
    "CurrentTick": 1
  },
  "Self": {
    "username": "0M3xh",
    "Color": "rgb(146, 166, 9)",
    "PosX": [
      7,
      6,
      6
    ],
    "PosY": [
      10,
      10,
      10
    ],
    "Head": {
      "X": 7,
      "Y": 10
    },
    "Tail": {
      "X": 6,
      "Y": 10
    },
    "Size": 3,
    "TotalScore": 0,
    "RoundScore": 0
  }
}
```

The structure contains the following:
* Numplayers: The number of players in this game <3
* Players: A map of the players where the key is their name
    * username: Their username
    * Color: The fillstyle property for that player. Generally not something you need to care about
    * PosX: An list of all the X positions of the snakes blocks
    * PosY: An list of all the Y positions of the snakes blocks
    * Head: The front block of the snake
    * Tail: The last block of the snake
    * Size: Number of blocks in a snake
    * TotalScore: The accumulated score of an snake this game
    * Roundscore: The score a snake has gotten this round
* GameStatus: All the map data and more
    * Status: Will generally only be ```Running``` when messages is received. 
        * The other states are as follows, but might not be exposed to the user:
            * pregame
            * initRound
            * running
            * roundDone
            * gameDone
        * RoundNumber: The current round number
        * TotalRounds: The total number of rounds
        * GameMap
            * SizeX: The size of the map in the X dimension
            * SizeY: THe size of the map in the Y dimension
            * Content: Contains a list of lists. The outer list is a list of rows, and the inner list is a list of the tiles in said row.
                * The possible tiles is as follows:
                    * Clear     = '_'
                    * Wall      = 'X'
                    * Snake     = '*'
                    * SnakeHead = 'H'
                    * Food      = '^' 
            * Heads: Will contain a list of head locations too (not implemented, use from the player object)
            * Walls: A list of dicts which shows where the walls are.
                * Each wall object contains a X and Y field
            * Foods: The same as walls, just with food.
        * CurrentTick: Contains data about which tick we are currently at
* Self: Contains your own object that is also found under ```Players```

### Moving your player
To move a player send the following object:

```JSON
{
  "type":"move",
  "value":"down"
}
```

- Value can be one of the following:
    * up
    * down
    * left
    * right
    
If no value is send each tick the server assigns a random one. 
This has a 25% chance of killing the snake due to crashing into itself.


## Running some tests
When you have the binary running you will find a server listening on [localhost:8080]()  
On top of this there is a simple debugging panel at [localhost:8080/debug.html]()

From this panel you can add the premade AI made in Javascript

You may also 



## How to compile

We need two things to run this server:
* [A installation of Go](https://golang.org/doc/install)
* A working browser (Tested in Firefox and Chrome)


We also need two external packages for Go. Just run the following commands in your terminal of choice:
```
go get github.com/google/logger
go get github.com/gorilla/websocket
```

Now you can either just run ```go build``` or ```go run``` 
