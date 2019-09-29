# Battlecards
Easily search for your a Starwars character, and get their collectible card. Collect as much Starwars character cards as you want!
You can search by name or planets. It also has an analytics tool to compare the stats of all characters in your collection!

## Development
Instructions on how to build and play with your own local copy of Battlecards.
If you already have Go(version 1.11 and above) and postgreSQL, skip ahead to the [building section](##Building-Battlecards).

### Prerequisites
You'll need these dependencies before getting started:
```
javascript
nodejs
npm
golang
postgresql
```

### MacOS
1) Download Golang (**Version 1.11 and up ONLY**) [here](https://golang.org/dl/)
2) Install postgreSQL
```
brew install postgreSQL
brew services start postgreSQL
```
3) Install nodeJS [here](https://nodejs.org/en/)

### Ubuntu
1) Download & Install golang
```
sudo apt-get update
sudo snap install go --classic
```
2) Install postgreSQL
```
sudo apt install postgresql postgresql-contrib
```

### Configuring postgreSQL (Ubuntu)
We will need to set up a user for postgreSQL, and create a database.

Login in to your postgres user account
```
sudo -i -u postgres
```
Create a new role called battlecards. Enter the command below and fill up the prompts
```
createuser --interactive
```
Create a new database, also called battlecards. Make sure you are still logged in to postgres user.
```
createdb battlecards
```
Exit from the postgres account by typing in `exit.

Let's add a useraccount which has the same name as the role we specified earlier, which in this case is `battlecards`

```
sudo add user battlecards
```
Switch over to this new battlecards account
```
sudo -i -u battlecards
```
Log into our `battlecards` database using the `battlecards` user we just created
```
psql
\conninfo
```
Typing in `\conninfo` gives us crucial information to connect to our DB.
```
You are connected to database "battlecards" as user "battlecards" via socket in "/var/run/postgresql" at port "5432".
```
While inside the `battlecards` database, configure the password for battlecards user
```
\password battlecards
```
Remember this password and we will set up our `.env` file.

### Configuring for MacOS
Install and initialize postgreSQL for MacOS
```
brew install postgresql
brew services start postgreSQL
```
Next create a database called `battlecards`
```
createdb battlecards
```
Login in to your new `battlecards` database
```
psql -d battlecards
```

Let's configure the password for the `battlecards` database for your machine username.
```
\password nizamsyahrul
```

Grab all necessary details by typing in :
```
\conninfo
```
We will be using those details to modify our .env file.
### Building Battlecards 
Follow these instructions to get Battlecards running after you've installed the dependencies

First, clone the repository and install all required node modules.
```
git clone https://github.com/syahrul12345/Battlecards.git
cd Battlecards
```
### Caching
The cached file lives in the `cache` folder as `Cache.txt`. If this is the first time that you run `/backend` and you have yet to make a request, `Cache.txt` does not exist. It will be automatically created upon your first request! 

### Compiling the frontend
Next, let's compile the frontend.
```
cd view
npm install
npm run build
```
The frontend will live in the root folder in a newly generated `/dist/` directory. Make sure `/dist/` is in the root folder. Our backend server will target the `index.html` file and serve it on `http://localhost:5555`.

###Compiling the backend

Enter the `backend` folder in the root repository.
```
cd ../
cd backend
```

The server is written in golang and lives in the `backend` folder. It acts as middlelayer between the frontend and the SWAPI, allowing us to re-index the api and to apply caching features. Caching is done by saving it to the `cache` folder in the root directory.

**Ensure that postgreSQL is running on your machine**
Since the backend requires to connect to a local postgreSQL database, we will need to configure our `.env` file. Change the contents of the `sample.env` in the backend. Replace the `db_name`, `db_pass` and `db_user` with the correct information. If you're not sure how to do this go [to the configuration part.](#Configuring-postgreSQL-(all-OS))

```
db_name = CHANGE_THIS_TO_YOUR_DB_NAME
db_pass = CHANGE_THIS_TO_YOUR_DB_PASSWORD
db_user = CHANGE_THIS_TO_YOUR_USER_ID_OF_MACHINE
db_type = postgres
db_host = localhost
db_port = 5432
```

Change your credentials accordingly. Rename the `sample.env` file to `.env`

```
mv sample.env .env
```

We can then compile the backend and execute the go binary generated.
```
go build
./backend
```

That's it! If you've followed this steps in order, the backend will automatically serve the `index.html` file in the `/dist/` folder. You can visit your own Battlecards website 

### Backend architecture & design

The purpose of the backend is to
1) Reindex the SWAPI so that is searchable by name
2) Stores the reindexed SWAPI in a local postgreSQL DB
3) Exposes an new set of API calls so the frontend can send requests to the backend with character name as as the payload
4) Provides a caching service.

## Reindexing

The backend implements an `indexer` module which first starts by making an RPC call to `"https://swapi.co/api/people/"`. 

The returning payload returns the `next` page link and the set of `characters`. The indexer spawns a goroutine(async thread) to handle the current set of `characters`. Since the original character object returned by SWAPI does not include specific information about vehicles,starship and the homeworld. Instead,only an array of links is provided. The indexer then spawns 1 goroutine each for every link and get the correct information.

Also reindexing is only done on the first time the binary is executed.

## Storage

Once the reindexing process is done, the CharacterResponse struct is modelled to a Character struct for saving in postgreSQL. Careful consideration have to be done when storing an array of Vehicle structs and Starship structs : `Vehicles []Vehicle` and `StarShips []StarShip`. Instead, Character struct uses a bytearray for `Vehicles` and `StarShips`.

## API

I've used a lightweight package called gorilla/mux to expose API endpoints. Mux is also used to serve the HTML built by Vue.JS. It exposes two endpoints `/api/character` and `/api/cache`. `api/character` is a `POST` request which requires the user to provide the name of character in the payload. This name is searched against a normalized string in the database. The normalization of strings in the database is to ensure that requests such as `ObiWANKEnobi` and `Obi-WanKEnobi` always returns `Obi-Wan Kenobi`. This can be seen in the `NameSearch` parameter in the `Character Struct`

```
//Character represents the JSON object of a StarWars character after reindexing
type Character struct {
	gorm.Model
	NameSearch string
	Name       string
	Gender     string
	StarShips  pq.ByteaArray `gorm:"type:varchar(1000)[]"`
	Vehicles   pq.ByteaArray `gorm:"type:varchar(1000)[]"`
	Home       pq.ByteaArray `gorm:"type:varchar(8000)[]"`
}
```
Once the correct Character struct is returned, the `StartShips`,`Vehicles` and `Home` is marshalled to the correct format of:

```
//CharacterResult is the JSON object send to the front end
type CharacterResult struct {
	Name      string
	Gender    string
	StarShips []StarShip
	Vehicles  []Vehicle
	HomeWorld Home
}
```

All models are stored in [character.go](https://github.com/syahrul12345/Battlecards/blob/master/backend/models/character.go)

## Caching Service
The cache file lives as a `Cache.txt`. Everytime a response is made, it automatically caches it in the text file. The caching logic is implemented on the `CharacterResult` struct, which is the object that will be returned to the frontend. It can be seen in [character.go](https://github.com/syahrul12345/Battlecards/blob/master/backend/models/character.go) as `Cache()`

Retrieving the cache is simply reading `Cache.txt` and sending it back to the frontend. Everytime the frontend is refreshed, it makes a GET request `/api/getCache/`. Everytime `/api/getCache/` is called, it checks if the currentTime and the time stored in `Cache.txt` is more than 7 days. If so, it deletes `Cache.txt` and returns an empty array. Otherwise, it simply reads `Cache.txt` and returns the json object written inside.

Caching logic can be seen in [characterController.go](https://github.com/syahrul12345/Battlecards/blob/master/backend/controllers/characterController.go)